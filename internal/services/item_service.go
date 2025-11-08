package services

import (
	"github.com/byte4bite/byte4bite/internal/models"
	"github.com/byte4bite/byte4bite/internal/repositories"
	"github.com/google/uuid"
)

// ItemService handles item business logic
type ItemService struct {
	itemRepo *repositories.ItemRepository
}

// NewItemService creates a new item service
func NewItemService(itemRepo *repositories.ItemRepository) *ItemService {
	return &ItemService{
		itemRepo: itemRepo,
	}
}

// CreateItemRequest represents an item creation request
type CreateItemRequest struct {
	Name              string     `json:"name" binding:"required"`
	Description       string     `json:"description"`
	CategoryID        uuid.UUID  `json:"category_id" binding:"required"`
	PantryID          uuid.UUID  `json:"pantry_id" binding:"required"`
	Quantity          int        `json:"quantity" binding:"required,min=0"`
	LowStockThreshold int        `json:"low_stock_threshold" binding:"required,min=0"`
	Unit              string     `json:"unit" binding:"required"`
	ImageURL          string     `json:"image_url"`
	IsAvailable       bool       `json:"is_available"`
}

// UpdateItemRequest represents an item update request
type UpdateItemRequest struct {
	Name              *string    `json:"name"`
	Description       *string    `json:"description"`
	CategoryID        *uuid.UUID `json:"category_id"`
	Quantity          *int       `json:"quantity"`
	LowStockThreshold *int       `json:"low_stock_threshold"`
	Unit              *string    `json:"unit"`
	ImageURL          *string    `json:"image_url"`
	IsAvailable       *bool      `json:"is_available"`
}

// ListItemsRequest represents a request to list items with filters
type ListItemsRequest struct {
	PantryID   *uuid.UUID `form:"pantry_id"`
	CategoryID *uuid.UUID `form:"category_id"`
	Search     string     `form:"search"`
	Available  *bool      `form:"available"`
	LowStock   bool       `form:"low_stock"`
	Page       int        `form:"page"`
	PageSize   int        `form:"page_size"`
}

// CreateItem creates a new item
func (s *ItemService) CreateItem(req *CreateItemRequest) (*models.Item, error) {
	item := &models.Item{
		Name:              req.Name,
		Description:       req.Description,
		CategoryID:        req.CategoryID,
		PantryID:          req.PantryID,
		Quantity:          req.Quantity,
		LowStockThreshold: req.LowStockThreshold,
		Unit:              req.Unit,
		ImageURL:          req.ImageURL,
		IsAvailable:       req.IsAvailable,
	}

	if err := s.itemRepo.Create(item); err != nil {
		return nil, err
	}

	// Reload to get associations
	return s.itemRepo.FindByID(item.ID)
}

// GetItem retrieves an item by ID
func (s *ItemService) GetItem(id uuid.UUID) (*models.Item, error) {
	return s.itemRepo.FindByID(id)
}

// UpdateItem updates an item
func (s *ItemService) UpdateItem(id uuid.UUID, req *UpdateItemRequest) (*models.Item, error) {
	item, err := s.itemRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Name != nil {
		item.Name = *req.Name
	}
	if req.Description != nil {
		item.Description = *req.Description
	}
	if req.CategoryID != nil {
		item.CategoryID = *req.CategoryID
	}
	if req.Quantity != nil {
		item.Quantity = *req.Quantity
	}
	if req.LowStockThreshold != nil {
		item.LowStockThreshold = *req.LowStockThreshold
	}
	if req.Unit != nil {
		item.Unit = *req.Unit
	}
	if req.ImageURL != nil {
		item.ImageURL = *req.ImageURL
	}
	if req.IsAvailable != nil {
		item.IsAvailable = *req.IsAvailable
	}

	if err := s.itemRepo.Update(item); err != nil {
		return nil, err
	}

	return s.itemRepo.FindByID(id)
}

// DeleteItem deletes an item
func (s *ItemService) DeleteItem(id uuid.UUID) error {
	// Check if item exists
	_, err := s.itemRepo.FindByID(id)
	if err != nil {
		return err
	}

	return s.itemRepo.Delete(id)
}

// ListItems lists items with filtering and pagination
func (s *ItemService) ListItems(req *ListItemsRequest) ([]models.Item, int64, error) {
	page := req.Page
	if page < 1 {
		page = 1
	}

	pageSize := req.PageSize
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100 // Max page size
	}

	offset := (page - 1) * pageSize

	filter := repositories.ItemFilter{
		PantryID:   req.PantryID,
		CategoryID: req.CategoryID,
		Search:     req.Search,
		Available:  req.Available,
		LowStock:   req.LowStock,
	}

	items, err := s.itemRepo.List(filter, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.itemRepo.Count(filter)
	if err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

// GetLowStockItems retrieves items that are low on stock
func (s *ItemService) GetLowStockItems(pantryID *uuid.UUID) ([]models.Item, error) {
	return s.itemRepo.FindLowStock(pantryID)
}

// UpdateItemQuantity updates the quantity of an item
func (s *ItemService) UpdateItemQuantity(id uuid.UUID, quantity int) error {
	// Verify item exists
	_, err := s.itemRepo.FindByID(id)
	if err != nil {
		return err
	}

	return s.itemRepo.UpdateQuantity(id, quantity)
}

// AdjustItemQuantity adjusts the quantity of an item by a delta
func (s *ItemService) AdjustItemQuantity(id uuid.UUID, delta int) error {
	// Verify item exists
	item, err := s.itemRepo.FindByID(id)
	if err != nil {
		return err
	}

	// Check that adjustment won't result in negative quantity
	if item.Quantity+delta < 0 {
		return err
	}

	return s.itemRepo.AdjustQuantity(id, delta)
}
