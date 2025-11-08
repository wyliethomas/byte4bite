package repositories

import (
	"errors"

	"github.com/byte4bite/byte4bite/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CategoryRepository handles database operations for categories
type CategoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository creates a new category repository
func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

// Create creates a new category
func (r *CategoryRepository) Create(category *models.Category) error {
	return r.db.Create(category).Error
}

// FindByID finds a category by ID
func (r *CategoryRepository) FindByID(id uuid.UUID) (*models.Category, error) {
	var category models.Category
	err := r.db.Preload("Pantry").First(&category, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, err
	}
	return &category, nil
}

// FindByPantryID finds all categories for a pantry
func (r *CategoryRepository) FindByPantryID(pantryID uuid.UUID) ([]models.Category, error) {
	var categories []models.Category
	err := r.db.Where("pantry_id = ?", pantryID).Find(&categories).Error
	return categories, err
}

// Update updates a category
func (r *CategoryRepository) Update(category *models.Category) error {
	return r.db.Save(category).Error
}

// Delete deletes a category
func (r *CategoryRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Category{}, "id = ?", id).Error
}

// List returns a list of categories with pagination
func (r *CategoryRepository) List(pantryID *uuid.UUID, limit, offset int) ([]models.Category, error) {
	var categories []models.Category
	query := r.db.Preload("Pantry")

	if pantryID != nil {
		query = query.Where("pantry_id = ?", *pantryID)
	}

	err := query.Limit(limit).Offset(offset).Find(&categories).Error
	return categories, err
}

// Count returns the total count of categories
func (r *CategoryRepository) Count(pantryID *uuid.UUID) (int64, error) {
	var count int64
	query := r.db.Model(&models.Category{})

	if pantryID != nil {
		query = query.Where("pantry_id = ?", *pantryID)
	}

	err := query.Count(&count).Error
	return count, err
}
