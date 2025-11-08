package handlers

import (
	"net/http"

	"github.com/byte4bite/byte4bite/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ItemHandler handles item-related endpoints
type ItemHandler struct {
	itemService *services.ItemService
}

// NewItemHandler creates a new item handler
func NewItemHandler(itemService *services.ItemService) *ItemHandler {
	return &ItemHandler{
		itemService: itemService,
	}
}

// CreateItem creates a new item
// POST /api/v1/admin/items
func (h *ItemHandler) CreateItem(c *gin.Context) {
	var req services.CreateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := h.itemService.CreateItem(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create item"})
		return
	}

	c.JSON(http.StatusCreated, item)
}

// GetItem retrieves an item by ID
// GET /api/v1/admin/items/:id
func (h *ItemHandler) GetItem(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	item, err := h.itemService.GetItem(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, item)
}

// ListItems lists all items with filtering and pagination
// GET /api/v1/admin/items
func (h *ItemHandler) ListItems(c *gin.Context) {
	var req services.ListItemsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	items, total, err := h.itemService.ListItems(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list items"})
		return
	}

	page := req.Page
	if page < 1 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize < 1 {
		pageSize = 20
	}

	c.JSON(http.StatusOK, gin.H{
		"data":        items,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
	})
}

// UpdateItem updates an item
// PUT /api/v1/admin/items/:id
func (h *ItemHandler) UpdateItem(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	var req services.UpdateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := h.itemService.UpdateItem(id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update item"})
		return
	}

	c.JSON(http.StatusOK, item)
}

// DeleteItem deletes an item
// DELETE /api/v1/admin/items/:id
func (h *ItemHandler) DeleteItem(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	if err := h.itemService.DeleteItem(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item deleted successfully"})
}

// UpdateItemQuantity updates the quantity of an item
// PATCH /api/v1/admin/items/:id/quantity
func (h *ItemHandler) UpdateItemQuantity(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	var req struct {
		Quantity int `json:"quantity" binding:"required,min=0"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.itemService.UpdateItemQuantity(id, req.Quantity); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update quantity"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Quantity updated successfully"})
}

// GetLowStockItems retrieves items that are low on stock
// GET /api/v1/admin/items/low-stock
func (h *ItemHandler) GetLowStockItems(c *gin.Context) {
	var pantryID *uuid.UUID
	if pantryIDParam := c.Query("pantry_id"); pantryIDParam != "" {
		id, err := uuid.Parse(pantryIDParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pantry ID"})
			return
		}
		pantryID = &id
	}

	items, err := h.itemService.GetLowStockItems(pantryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get low stock items"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  items,
		"count": len(items),
	})
}

// ListItemsPublic lists available items for regular users
// GET /api/v1/items
func (h *ItemHandler) ListItemsPublic(c *gin.Context) {
	var req services.ListItemsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Force available filter for public endpoint
	available := true
	req.Available = &available

	items, total, err := h.itemService.ListItems(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list items"})
		return
	}

	page := req.Page
	if page < 1 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize < 1 {
		pageSize = 20
	}

	c.JSON(http.StatusOK, gin.H{
		"data":        items,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
	})
}
