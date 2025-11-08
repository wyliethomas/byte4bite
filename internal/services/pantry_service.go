package services

import (
	"errors"

	"github.com/byte4bite/byte4bite/internal/models"
	"github.com/byte4bite/byte4bite/internal/repositories"
	"github.com/google/uuid"
)

// PantryService handles business logic for pantries
type PantryService struct {
	pantryRepo *repositories.PantryRepository
}

// NewPantryService creates a new pantry service
func NewPantryService(pantryRepo *repositories.PantryRepository) *PantryService {
	return &PantryService{
		pantryRepo: pantryRepo,
	}
}

// CreatePantryRequest represents a request to create a pantry
type CreatePantryRequest struct {
	Name         string `json:"name" binding:"required"`
	Address      string `json:"address" binding:"required"`
	City         string `json:"city" binding:"required"`
	State        string `json:"state" binding:"required"`
	ZipCode      string `json:"zip_code" binding:"required"`
	ContactEmail string `json:"contact_email" binding:"required,email"`
	ContactPhone string `json:"contact_phone"`
	IsActive     bool   `json:"is_active"`
}

// UpdatePantryRequest represents a request to update a pantry
type UpdatePantryRequest struct {
	Name         *string `json:"name"`
	Address      *string `json:"address"`
	City         *string `json:"city"`
	State        *string `json:"state"`
	ZipCode      *string `json:"zip_code"`
	ContactEmail *string `json:"contact_email"`
	ContactPhone *string `json:"contact_phone"`
	IsActive     *bool   `json:"is_active"`
}

// GetPantriesRequest represents a request to get pantries
type GetPantriesRequest struct {
	IsActive *bool
	Page     int
	PageSize int
}

// GetPantriesResponse represents the response containing pantries
type GetPantriesResponse struct {
	Pantries []models.Pantry `json:"pantries"`
	Total    int64           `json:"total"`
	Page     int             `json:"page"`
	Pages    int             `json:"pages"`
}

// CreatePantry creates a new pantry
func (s *PantryService) CreatePantry(req *CreatePantryRequest) (*models.Pantry, error) {
	pantry := &models.Pantry{
		Name:         req.Name,
		Address:      req.Address,
		City:         req.City,
		State:        req.State,
		ZipCode:      req.ZipCode,
		ContactEmail: req.ContactEmail,
		ContactPhone: req.ContactPhone,
		IsActive:     req.IsActive,
	}

	if err := s.pantryRepo.Create(pantry); err != nil {
		return nil, err
	}

	return pantry, nil
}

// GetPantry retrieves a pantry by ID
func (s *PantryService) GetPantry(id uuid.UUID) (*models.Pantry, error) {
	return s.pantryRepo.FindByID(id)
}

// GetPantries retrieves a list of pantries
func (s *PantryService) GetPantries(req GetPantriesRequest) (*GetPantriesResponse, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 100 {
		req.PageSize = 20
	}

	offset := (req.Page - 1) * req.PageSize

	pantries, err := s.pantryRepo.FindAll(req.IsActive, req.PageSize, offset)
	if err != nil {
		return nil, err
	}

	total, err := s.pantryRepo.Count(req.IsActive)
	if err != nil {
		return nil, err
	}

	pages := int(total) / req.PageSize
	if int(total)%req.PageSize != 0 {
		pages++
	}

	return &GetPantriesResponse{
		Pantries: pantries,
		Total:    total,
		Page:     req.Page,
		Pages:    pages,
	}, nil
}

// UpdatePantry updates a pantry
func (s *PantryService) UpdatePantry(id uuid.UUID, req *UpdatePantryRequest) (*models.Pantry, error) {
	pantry, err := s.pantryRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Name != nil {
		pantry.Name = *req.Name
	}
	if req.Address != nil {
		pantry.Address = *req.Address
	}
	if req.City != nil {
		pantry.City = *req.City
	}
	if req.State != nil {
		pantry.State = *req.State
	}
	if req.ZipCode != nil {
		pantry.ZipCode = *req.ZipCode
	}
	if req.ContactEmail != nil {
		pantry.ContactEmail = *req.ContactEmail
	}
	if req.ContactPhone != nil {
		pantry.ContactPhone = *req.ContactPhone
	}
	if req.IsActive != nil {
		pantry.IsActive = *req.IsActive
	}

	if err := s.pantryRepo.Update(pantry); err != nil {
		return nil, err
	}

	return pantry, nil
}

// DeletePantry deletes a pantry
func (s *PantryService) DeletePantry(id uuid.UUID) error {
	// Check if pantry exists
	_, err := s.pantryRepo.FindByID(id)
	if err != nil {
		return err
	}

	return s.pantryRepo.Delete(id)
}

// SearchPantries searches for pantries by query
func (s *PantryService) SearchPantries(query string) ([]models.Pantry, error) {
	if query == "" {
		return nil, errors.New("search query cannot be empty")
	}
	return s.pantryRepo.Search(query)
}

// GetPantriesByCity retrieves pantries in a specific city
func (s *PantryService) GetPantriesByCity(city string) ([]models.Pantry, error) {
	if city == "" {
		return nil, errors.New("city cannot be empty")
	}
	return s.pantryRepo.FindByCity(city)
}

// GetPantriesByZipCode retrieves pantries in a specific zip code
func (s *PantryService) GetPantriesByZipCode(zipCode string) ([]models.Pantry, error) {
	if zipCode == "" {
		return nil, errors.New("zip code cannot be empty")
	}
	return s.pantryRepo.FindByZipCode(zipCode)
}

// TogglePantryStatus toggles the active status of a pantry
func (s *PantryService) TogglePantryStatus(id uuid.UUID) (*models.Pantry, error) {
	pantry, err := s.pantryRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	pantry.IsActive = !pantry.IsActive
	if err := s.pantryRepo.Update(pantry); err != nil {
		return nil, err
	}

	return pantry, nil
}
