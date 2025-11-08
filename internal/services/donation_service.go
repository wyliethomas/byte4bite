package services

import (
	"errors"
	"time"

	"github.com/byte4bite/byte4bite/internal/models"
	"github.com/byte4bite/byte4bite/internal/repositories"
	"github.com/google/uuid"
)

// DonationService handles business logic for donations
type DonationService struct {
	donationRepo *repositories.DonationRepository
	pantryRepo   *repositories.PantryRepository
}

// NewDonationService creates a new donation service
func NewDonationService(donationRepo *repositories.DonationRepository, pantryRepo *repositories.PantryRepository) *DonationService {
	return &DonationService{
		donationRepo: donationRepo,
		pantryRepo:   pantryRepo,
	}
}

// CreateDonationRequest represents a request to create a donation
type CreateDonationRequest struct {
	PantryID     uuid.UUID `json:"pantry_id" binding:"required"`
	DonorName    string    `json:"donor_name" binding:"required"`
	DonorEmail   string    `json:"donor_email" binding:"email"`
	DonorPhone   string    `json:"donor_phone"`
	Amount       *float64  `json:"amount"`
	Description  string    `json:"description" binding:"required"`
	DonationDate time.Time `json:"donation_date"`
}

// UpdateDonationRequest represents a request to update a donation
type UpdateDonationRequest struct {
	DonorName    *string    `json:"donor_name"`
	DonorEmail   *string    `json:"donor_email"`
	DonorPhone   *string    `json:"donor_phone"`
	Amount       *float64   `json:"amount"`
	Description  *string    `json:"description"`
	DonationDate *time.Time `json:"donation_date"`
	ReceiptSent  *bool      `json:"receipt_sent"`
}

// GetDonationsRequest represents a request to get donations
type GetDonationsRequest struct {
	PantryID    *uuid.UUID
	ReceiptSent *bool
	StartDate   *time.Time
	EndDate     *time.Time
	Page        int
	PageSize    int
}

// GetDonationsResponse represents the response containing donations
type GetDonationsResponse struct {
	Donations []models.Donation `json:"donations"`
	Total     int64             `json:"total"`
	Page      int               `json:"page"`
	Pages     int               `json:"pages"`
}

// DonationStatsResponse represents donation statistics
type DonationStatsResponse struct {
	TotalDonations   int64   `json:"total_donations"`
	TotalAmount      float64 `json:"total_amount"`
	DonorCount       int64   `json:"donor_count"`
	ReceiptsPending  int64   `json:"receipts_pending"`
	MonetaryCount    int64   `json:"monetary_count"`
	InKindCount      int64   `json:"in_kind_count"`
}

// CreateDonation creates a new donation
func (s *DonationService) CreateDonation(req *CreateDonationRequest) (*models.Donation, error) {
	// Verify pantry exists
	_, err := s.pantryRepo.FindByID(req.PantryID)
	if err != nil {
		return nil, errors.New("pantry not found")
	}

	// Validate amount if provided
	if req.Amount != nil && *req.Amount < 0 {
		return nil, errors.New("donation amount cannot be negative")
	}

	// Set donation date to now if not provided
	donationDate := req.DonationDate
	if donationDate.IsZero() {
		donationDate = time.Now()
	}

	donation := &models.Donation{
		PantryID:     req.PantryID,
		DonorName:    req.DonorName,
		DonorEmail:   req.DonorEmail,
		DonorPhone:   req.DonorPhone,
		Amount:       req.Amount,
		Description:  req.Description,
		DonationDate: donationDate,
		ReceiptSent:  false,
	}

	if err := s.donationRepo.Create(donation); err != nil {
		return nil, err
	}

	// Reload to get pantry association
	return s.donationRepo.FindByID(donation.ID)
}

// GetDonation retrieves a donation by ID
func (s *DonationService) GetDonation(id uuid.UUID) (*models.Donation, error) {
	return s.donationRepo.FindByID(id)
}

// GetDonations retrieves a list of donations
func (s *DonationService) GetDonations(req GetDonationsRequest) (*GetDonationsResponse, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 100 {
		req.PageSize = 20
	}

	offset := (req.Page - 1) * req.PageSize

	donations, err := s.donationRepo.FindAll(
		req.PantryID,
		req.ReceiptSent,
		req.StartDate,
		req.EndDate,
		req.PageSize,
		offset,
	)
	if err != nil {
		return nil, err
	}

	total, err := s.donationRepo.Count(
		req.PantryID,
		req.ReceiptSent,
		req.StartDate,
		req.EndDate,
	)
	if err != nil {
		return nil, err
	}

	pages := int(total) / req.PageSize
	if int(total)%req.PageSize != 0 {
		pages++
	}

	return &GetDonationsResponse{
		Donations: donations,
		Total:     total,
		Page:      req.Page,
		Pages:     pages,
	}, nil
}

// UpdateDonation updates a donation
func (s *DonationService) UpdateDonation(id uuid.UUID, req *UpdateDonationRequest) (*models.Donation, error) {
	donation, err := s.donationRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.DonorName != nil {
		donation.DonorName = *req.DonorName
	}
	if req.DonorEmail != nil {
		donation.DonorEmail = *req.DonorEmail
	}
	if req.DonorPhone != nil {
		donation.DonorPhone = *req.DonorPhone
	}
	if req.Amount != nil {
		if *req.Amount < 0 {
			return nil, errors.New("donation amount cannot be negative")
		}
		donation.Amount = req.Amount
	}
	if req.Description != nil {
		donation.Description = *req.Description
	}
	if req.DonationDate != nil {
		donation.DonationDate = *req.DonationDate
	}
	if req.ReceiptSent != nil {
		donation.ReceiptSent = *req.ReceiptSent
	}

	if err := s.donationRepo.Update(donation); err != nil {
		return nil, err
	}

	return donation, nil
}

// DeleteDonation deletes a donation
func (s *DonationService) DeleteDonation(id uuid.UUID) error {
	// Check if donation exists
	_, err := s.donationRepo.FindByID(id)
	if err != nil {
		return err
	}

	return s.donationRepo.Delete(id)
}

// MarkReceiptSent marks a donation receipt as sent
func (s *DonationService) MarkReceiptSent(id uuid.UUID) (*models.Donation, error) {
	if err := s.donationRepo.MarkReceiptSent(id); err != nil {
		return nil, err
	}

	return s.donationRepo.FindByID(id)
}

// GetDonationsByDonor gets all donations from a specific donor
func (s *DonationService) GetDonationsByDonor(email string) ([]models.Donation, error) {
	if email == "" {
		return nil, errors.New("email is required")
	}
	return s.donationRepo.FindByDonorEmail(email)
}

// SearchDonations searches donations by donor name or description
func (s *DonationService) SearchDonations(query string, page, pageSize int) (*GetDonationsResponse, error) {
	if query == "" {
		return nil, errors.New("search query cannot be empty")
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	donations, err := s.donationRepo.Search(query, pageSize, offset)
	if err != nil {
		return nil, err
	}

	// For search, we'll use the count of results as an approximation
	// (In production, you'd want a proper count query)
	total := int64(len(donations))

	pages := int(total) / pageSize
	if int(total)%pageSize != 0 {
		pages++
	}

	return &GetDonationsResponse{
		Donations: donations,
		Total:     total,
		Page:      page,
		Pages:     pages,
	}, nil
}

// GetDonationStats retrieves donation statistics
func (s *DonationService) GetDonationStats(pantryID *uuid.UUID, startDate, endDate *time.Time) (*DonationStatsResponse, error) {
	// Get total donations count
	totalDonations, err := s.donationRepo.Count(pantryID, nil, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Get total amount
	totalAmount, err := s.donationRepo.GetTotalDonations(pantryID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Get donor count
	donorCount, err := s.donationRepo.GetDonorCount(pantryID)
	if err != nil {
		return nil, err
	}

	// Get receipts pending count
	receiptSent := false
	receiptsPending, err := s.donationRepo.Count(pantryID, &receiptSent, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Count monetary vs in-kind donations
	// This is a simplified approach - you might want to add specific queries
	donations, err := s.donationRepo.FindAll(pantryID, nil, startDate, endDate, 10000, 0)
	if err != nil {
		return nil, err
	}

	monetaryCount := int64(0)
	inKindCount := int64(0)
	for _, d := range donations {
		if d.Amount != nil && *d.Amount > 0 {
			monetaryCount++
		} else {
			inKindCount++
		}
	}

	return &DonationStatsResponse{
		TotalDonations:  totalDonations,
		TotalAmount:     totalAmount,
		DonorCount:      donorCount,
		ReceiptsPending: receiptsPending,
		MonetaryCount:   monetaryCount,
		InKindCount:     inKindCount,
	}, nil
}
