package repositories

import (
	"errors"
	"time"

	"github.com/byte4bite/byte4bite/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DonationRepository handles database operations for donations
type DonationRepository struct {
	db *gorm.DB
}

// NewDonationRepository creates a new donation repository
func NewDonationRepository(db *gorm.DB) *DonationRepository {
	return &DonationRepository{db: db}
}

// Create creates a new donation
func (r *DonationRepository) Create(donation *models.Donation) error {
	return r.db.Create(donation).Error
}

// FindByID finds a donation by ID
func (r *DonationRepository) FindByID(id uuid.UUID) (*models.Donation, error) {
	var donation models.Donation
	err := r.db.Preload("Pantry").First(&donation, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("donation not found")
		}
		return nil, err
	}
	return &donation, nil
}

// Update updates a donation
func (r *DonationRepository) Update(donation *models.Donation) error {
	return r.db.Save(donation).Error
}

// Delete deletes a donation
func (r *DonationRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Donation{}, "id = ?", id).Error
}

// FindAll finds all donations with optional filters
func (r *DonationRepository) FindAll(pantryID *uuid.UUID, receiptSent *bool, startDate, endDate *time.Time, limit, offset int) ([]models.Donation, error) {
	var donations []models.Donation
	query := r.db.Preload("Pantry")

	if pantryID != nil {
		query = query.Where("pantry_id = ?", *pantryID)
	}

	if receiptSent != nil {
		query = query.Where("receipt_sent = ?", *receiptSent)
	}

	if startDate != nil {
		query = query.Where("donation_date >= ?", *startDate)
	}

	if endDate != nil {
		query = query.Where("donation_date <= ?", *endDate)
	}

	err := query.Order("donation_date DESC").
		Limit(limit).Offset(offset).
		Find(&donations).Error
	return donations, err
}

// Count counts all donations with optional filters
func (r *DonationRepository) Count(pantryID *uuid.UUID, receiptSent *bool, startDate, endDate *time.Time) (int64, error) {
	var count int64
	query := r.db.Model(&models.Donation{})

	if pantryID != nil {
		query = query.Where("pantry_id = ?", *pantryID)
	}

	if receiptSent != nil {
		query = query.Where("receipt_sent = ?", *receiptSent)
	}

	if startDate != nil {
		query = query.Where("donation_date >= ?", *startDate)
	}

	if endDate != nil {
		query = query.Where("donation_date <= ?", *endDate)
	}

	err := query.Count(&count).Error
	return count, err
}

// FindByPantryID finds all donations for a pantry
func (r *DonationRepository) FindByPantryID(pantryID uuid.UUID, limit, offset int) ([]models.Donation, error) {
	var donations []models.Donation
	err := r.db.Where("pantry_id = ?", pantryID).
		Order("donation_date DESC").
		Limit(limit).Offset(offset).
		Find(&donations).Error
	return donations, err
}

// FindByDonorEmail finds all donations by donor email
func (r *DonationRepository) FindByDonorEmail(email string) ([]models.Donation, error) {
	var donations []models.Donation
	err := r.db.Preload("Pantry").
		Where("LOWER(donor_email) = LOWER(?)", email).
		Order("donation_date DESC").
		Find(&donations).Error
	return donations, err
}

// MarkReceiptSent marks a donation's receipt as sent
func (r *DonationRepository) MarkReceiptSent(id uuid.UUID) error {
	return r.db.Model(&models.Donation{}).
		Where("id = ?", id).
		Update("receipt_sent", true).Error
}

// GetTotalDonations calculates total monetary donations
func (r *DonationRepository) GetTotalDonations(pantryID *uuid.UUID, startDate, endDate *time.Time) (float64, error) {
	var total float64
	query := r.db.Model(&models.Donation{}).
		Select("COALESCE(SUM(amount), 0)")

	if pantryID != nil {
		query = query.Where("pantry_id = ?", *pantryID)
	}

	if startDate != nil {
		query = query.Where("donation_date >= ?", *startDate)
	}

	if endDate != nil {
		query = query.Where("donation_date <= ?", *endDate)
	}

	err := query.Scan(&total).Error
	return total, err
}

// GetDonorCount counts unique donors
func (r *DonationRepository) GetDonorCount(pantryID *uuid.UUID) (int64, error) {
	var count int64
	query := r.db.Model(&models.Donation{}).
		Select("COUNT(DISTINCT donor_email)")

	if pantryID != nil {
		query = query.Where("pantry_id = ?", *pantryID)
	}

	err := query.Scan(&count).Error
	return count, err
}

// Search searches donations by donor name or description
func (r *DonationRepository) Search(query string, limit, offset int) ([]models.Donation, error) {
	var donations []models.Donation
	searchPattern := "%" + query + "%"
	err := r.db.Preload("Pantry").
		Where("LOWER(donor_name) LIKE LOWER(?) OR LOWER(description) LIKE LOWER(?)",
			searchPattern, searchPattern).
		Order("donation_date DESC").
		Limit(limit).Offset(offset).
		Find(&donations).Error
	return donations, err
}
