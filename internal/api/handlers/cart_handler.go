package handlers

import (
	"net/http"

	"github.com/byte4bite/byte4bite/internal/repositories"
	"github.com/byte4bite/byte4bite/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CartHandler handles cart-related endpoints
type CartHandler struct {
	cartService  *services.CartService
	orderRepo    *repositories.OrderRepository
}

// NewCartHandler creates a new cart handler
func NewCartHandler(cartService *services.CartService, orderRepo *repositories.OrderRepository) *CartHandler {
	return &CartHandler{
		cartService: cartService,
		orderRepo:   orderRepo,
	}
}

// GetCurrentCart retrieves the current active cart for the user
// GET /api/v1/carts/current
func (h *CartHandler) GetCurrentCart(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	cart, err := h.cartService.GetCurrentCart(userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get cart"})
		return
	}

	if cart == nil {
		c.JSON(http.StatusOK, gin.H{
			"cart":  nil,
			"items": []interface{}{},
			"count": 0,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"cart":  cart,
		"items": cart.Items,
		"count": len(cart.Items),
	})
}

// AddItem adds an item to the cart
// POST /api/v1/carts/items
func (h *CartHandler) AddItem(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// For now, using a default pantry ID - should come from user context or request
	pantryID := uuid.MustParse("00000000-0000-0000-0000-000000000000")
	if pantryIDValue, exists := c.Get("pantry_id"); exists {
		pantryID = pantryIDValue.(uuid.UUID)
	}

	var req services.AddItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cart, err := h.cartService.AddItem(userID.(uuid.UUID), pantryID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cart)
}

// UpdateItemQuantity updates the quantity of an item in the cart
// PUT /api/v1/carts/items/:id
func (h *CartHandler) UpdateItemQuantity(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	cartItemID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cart item ID"})
		return
	}

	var req services.UpdateCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cart, err := h.cartService.UpdateItemQuantity(userID.(uuid.UUID), cartItemID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cart)
}

// RemoveItem removes an item from the cart
// DELETE /api/v1/carts/items/:id
func (h *CartHandler) RemoveItem(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	cartItemID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cart item ID"})
		return
	}

	cart, err := h.cartService.RemoveItem(userID.(uuid.UUID), cartItemID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cart)
}

// ClearCart removes all items from the cart
// DELETE /api/v1/carts/current
func (h *CartHandler) ClearCart(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if err := h.cartService.ClearCart(userID.(uuid.UUID)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cart cleared successfully"})
}

// Checkout converts the cart to an order
// POST /api/v1/carts/checkout
func (h *CartHandler) Checkout(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req struct {
		Notes string `json:"notes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		// Notes are optional, so ignore binding errors
		req.Notes = ""
	}

	order, err := h.cartService.Checkout(userID.(uuid.UUID), req.Notes)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save the order
	if err := h.orderRepo.Create(order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	// Reload order with associations
	order, err = h.orderRepo.FindByID(order.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve order"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Order created successfully",
		"order":   order,
	})
}
