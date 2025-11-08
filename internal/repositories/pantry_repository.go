package repositories

import (
	"errors"

	"github.com/byte4bite/byte4bite/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PantryRepository handles database operations for pantries
type PantryRepository struct {
	db *gorm.DB
}

// NewPantryRepository creates a new pantry repository
func NewPantryRepository(db *gorm.DB) *PantryRepository {
	return &PantryRepository{db: db}
}

// Create creates a new pantry
func (r *PantryRepository) Create(pantry *models.Pantry) error {
	return r.db.Create(pantry).Error
}

// FindByID finds a pantry by ID
func (r *PantryRepository) FindByID(id uuid.UUID) (*models.Pantry, error) {
	var pantry models.Pantry
	err := r.db.First(&pantry, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("pantry not found")
		}
		return nil, err
	}
	return &pantry, nil
}

// Update updates a pantry
func (r *PantryRepository) Update(pantry *models.Pantry) error {
	return r.db.Save(pantry).Error
}

// Delete deletes a pantry
func (r *PantryRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Pantry{}, "id = ?", id).Error
}

// FindAll finds all pantries with optional filters
func (r *PantryRepository) FindAll(isActive *bool, limit, offset int) ([]models.Pantry, error) {
	var pantries []models.Pantry
	query := r.db.Model(&models.Pantry{})

	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}

	err := query.Order("name ASC").
		Limit(limit).Offset(offset).
		Find(&pantries).Error
	return pantries, err
}

// Count counts all pantries with optional filters
func (r *PantryRepository) Count(isActive *bool) (int64, error) {
	var count int64
	query := r.db.Model(&models.Pantry{})

	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}

	err := query.Count(&count).Error
	return count, err
}

// FindByCity finds pantries by city
func (r *PantryRepository) FindByCity(city string) ([]models.Pantry, error) {
	var pantries []models.Pantry
	err := r.db.Where("LOWER(city) = LOWER(?) AND is_active = ?", city, true).
		Order("name ASC").
		Find(&pantries).Error
	return pantries, err
}

// FindByZipCode finds pantries by zip code
func (r *PantryRepository) FindByZipCode(zipCode string) ([]models.Pantry, error) {
	var pantries []models.Pantry
	err := r.db.Where("zip_code = ? AND is_active = ?", zipCode, true).
		Order("name ASC").
		Find(&pantries).Error
	return pantries, err
}

// Search searches pantries by name or city
func (r *PantryRepository) Search(query string) ([]models.Pantry, error) {
	var pantries []models.Pantry
	searchPattern := "%" + query + "%"
	err := r.db.Where("(LOWER(name) LIKE LOWER(?) OR LOWER(city) LIKE LOWER(?)) AND is_active = ?",
		searchPattern, searchPattern, true).
		Order("name ASC").
		Find(&pantries).Error
	return pantries, err
}

// UpdateActiveStatus updates the active status of a pantry
func (r *PantryRepository) UpdateActiveStatus(id uuid.UUID, isActive bool) error {
	return r.db.Model(&models.Pantry{}).Where("id = ?", id).
		Update("is_active", isActive).Error
}
