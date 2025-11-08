package services

import (
	"github.com/byte4bite/byte4bite/internal/models"
	"github.com/byte4bite/byte4bite/internal/repositories"
	"github.com/google/uuid"
)

// CategoryService handles category business logic
type CategoryService struct {
	categoryRepo *repositories.CategoryRepository
}

// NewCategoryService creates a new category service
func NewCategoryService(categoryRepo *repositories.CategoryRepository) *CategoryService {
	return &CategoryService{
		categoryRepo: categoryRepo,
	}
}

// CreateCategoryRequest represents a category creation request
type CreateCategoryRequest struct {
	Name        string     `json:"name" binding:"required"`
	Description string     `json:"description"`
	PantryID    uuid.UUID  `json:"pantry_id" binding:"required"`
}

// UpdateCategoryRequest represents a category update request
type UpdateCategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// CreateCategory creates a new category
func (s *CategoryService) CreateCategory(req *CreateCategoryRequest) (*models.Category, error) {
	category := &models.Category{
		Name:        req.Name,
		Description: req.Description,
		PantryID:    req.PantryID,
	}

	if err := s.categoryRepo.Create(category); err != nil {
		return nil, err
	}

	return category, nil
}

// GetCategory retrieves a category by ID
func (s *CategoryService) GetCategory(id uuid.UUID) (*models.Category, error) {
	return s.categoryRepo.FindByID(id)
}

// GetCategoriesByPantry retrieves all categories for a pantry
func (s *CategoryService) GetCategoriesByPantry(pantryID uuid.UUID) ([]models.Category, error) {
	return s.categoryRepo.FindByPantryID(pantryID)
}

// UpdateCategory updates a category
func (s *CategoryService) UpdateCategory(id uuid.UUID, req *UpdateCategoryRequest) (*models.Category, error) {
	category, err := s.categoryRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		category.Name = req.Name
	}
	if req.Description != "" {
		category.Description = req.Description
	}

	if err := s.categoryRepo.Update(category); err != nil {
		return nil, err
	}

	return category, nil
}

// DeleteCategory deletes a category
func (s *CategoryService) DeleteCategory(id uuid.UUID) error {
	// Check if category exists
	_, err := s.categoryRepo.FindByID(id)
	if err != nil {
		return err
	}

	return s.categoryRepo.Delete(id)
}

// ListCategories lists categories with pagination
func (s *CategoryService) ListCategories(pantryID *uuid.UUID, page, pageSize int) ([]models.Category, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	categories, err := s.categoryRepo.List(pantryID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.categoryRepo.Count(pantryID)
	if err != nil {
		return nil, 0, err
	}

	return categories, total, nil
}
