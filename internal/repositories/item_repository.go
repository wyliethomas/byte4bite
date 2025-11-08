package repositories

import (
	"errors"
	"strings"

	"github.com/byte4bite/byte4bite/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ItemRepository handles database operations for items
type ItemRepository struct {
	db *gorm.DB
}

// NewItemRepository creates a new item repository
func NewItemRepository(db *gorm.DB) *ItemRepository {
	return &ItemRepository{db: db}
}

// ItemFilter represents filtering options for items
type ItemFilter struct {
	PantryID   *uuid.UUID
	CategoryID *uuid.UUID
	Search     string
	Available  *bool
	LowStock   bool
}

// Create creates a new item
func (r *ItemRepository) Create(item *models.Item) error {
	return r.db.Create(item).Error
}

// FindByID finds an item by ID
func (r *ItemRepository) FindByID(id uuid.UUID) (*models.Item, error) {
	var item models.Item
	err := r.db.Preload("Category").Preload("Pantry").First(&item, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("item not found")
		}
		return nil, err
	}
	return &item, nil
}

// Update updates an item
func (r *ItemRepository) Update(item *models.Item) error {
	return r.db.Save(item).Error
}

// Delete deletes an item
func (r *ItemRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Item{}, "id = ?", id).Error
}

// List returns a list of items with filtering and pagination
func (r *ItemRepository) List(filter ItemFilter, limit, offset int) ([]models.Item, error) {
	var items []models.Item
	query := r.db.Preload("Category").Preload("Pantry")

	// Apply filters
	query = r.applyFilters(query, filter)

	err := query.Limit(limit).Offset(offset).Find(&items).Error
	return items, err
}

// Count returns the total count of items matching the filter
func (r *ItemRepository) Count(filter ItemFilter) (int64, error) {
	var count int64
	query := r.db.Model(&models.Item{})

	// Apply filters
	query = r.applyFilters(query, filter)

	err := query.Count(&count).Error
	return count, err
}

// FindLowStock finds items that are low on stock
func (r *ItemRepository) FindLowStock(pantryID *uuid.UUID) ([]models.Item, error) {
	var items []models.Item
	query := r.db.Preload("Category").Preload("Pantry").
		Where("quantity <= low_stock_threshold")

	if pantryID != nil {
		query = query.Where("pantry_id = ?", *pantryID)
	}

	err := query.Find(&items).Error
	return items, err
}

// UpdateQuantity updates the quantity of an item
func (r *ItemRepository) UpdateQuantity(id uuid.UUID, quantity int) error {
	return r.db.Model(&models.Item{}).Where("id = ?", id).Update("quantity", quantity).Error
}

// AdjustQuantity adjusts the quantity of an item by a delta (can be negative)
func (r *ItemRepository) AdjustQuantity(id uuid.UUID, delta int) error {
	return r.db.Model(&models.Item{}).Where("id = ?", id).
		UpdateColumn("quantity", gorm.Expr("quantity + ?", delta)).Error
}

// applyFilters applies filtering conditions to a query
func (r *ItemRepository) applyFilters(query *gorm.DB, filter ItemFilter) *gorm.DB {
	if filter.PantryID != nil {
		query = query.Where("pantry_id = ?", *filter.PantryID)
	}

	if filter.CategoryID != nil {
		query = query.Where("category_id = ?", *filter.CategoryID)
	}

	if filter.Search != "" {
		searchTerm := "%" + strings.ToLower(filter.Search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(description) LIKE ?", searchTerm, searchTerm)
	}

	if filter.Available != nil {
		query = query.Where("is_available = ?", *filter.Available)
	}

	if filter.LowStock {
		query = query.Where("quantity <= low_stock_threshold")
	}

	return query
}
