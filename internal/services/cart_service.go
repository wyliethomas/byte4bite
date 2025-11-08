package services

import (
	"errors"

	"github.com/byte4bite/byte4bite/internal/models"
	"github.com/byte4bite/byte4bite/internal/repositories"
	"github.com/google/uuid"
)

// CartService handles cart business logic
type CartService struct {
	cartRepo *repositories.CartRepository
	itemRepo *repositories.ItemRepository
}

// NewCartService creates a new cart service
func NewCartService(cartRepo *repositories.CartRepository, itemRepo *repositories.ItemRepository) *CartService {
	return &CartService{
		cartRepo: cartRepo,
		itemRepo: itemRepo,
	}
}

// AddItemRequest represents a request to add an item to cart
type AddItemRequest struct {
	ItemID   uuid.UUID `json:"item_id" binding:"required"`
	Quantity int       `json:"quantity" binding:"required,min=1"`
}

// UpdateCartItemRequest represents a request to update cart item quantity
type UpdateCartItemRequest struct {
	Quantity int `json:"quantity" binding:"required,min=0"`
}

// GetOrCreateCart gets the active cart for a user or creates a new one
func (s *CartService) GetOrCreateCart(userID, pantryID uuid.UUID) (*models.Cart, error) {
	// Try to find active cart
	cart, err := s.cartRepo.FindActiveByUserID(userID)
	if err != nil {
		return nil, err
	}

	// If no active cart exists, create one
	if cart == nil {
		cart = &models.Cart{
			UserID:   userID,
			PantryID: pantryID,
			Status:   models.CartStatusActive,
		}
		if err := s.cartRepo.Create(cart); err != nil {
			return nil, err
		}

		// Reload to get associations
		cart, err = s.cartRepo.FindByID(cart.ID)
		if err != nil {
			return nil, err
		}
	}

	return cart, nil
}

// GetCart retrieves a cart by ID
func (s *CartService) GetCart(cartID uuid.UUID) (*models.Cart, error) {
	return s.cartRepo.FindByID(cartID)
}

// GetCurrentCart gets the current active cart for a user
func (s *CartService) GetCurrentCart(userID uuid.UUID) (*models.Cart, error) {
	return s.cartRepo.FindActiveByUserID(userID)
}

// AddItem adds an item to the cart
func (s *CartService) AddItem(userID, pantryID uuid.UUID, req *AddItemRequest) (*models.Cart, error) {
	// Verify item exists and is available
	item, err := s.itemRepo.FindByID(req.ItemID)
	if err != nil {
		return nil, errors.New("item not found")
	}

	if !item.IsAvailable {
		return nil, errors.New("item is not available")
	}

	if item.Quantity < req.Quantity {
		return nil, errors.New("insufficient quantity available")
	}

	// Get or create cart
	cart, err := s.GetOrCreateCart(userID, pantryID)
	if err != nil {
		return nil, err
	}

	// Check if item already in cart
	existingCartItem, err := s.cartRepo.FindCartItem(cart.ID, req.ItemID)
	if err != nil {
		return nil, err
	}

	if existingCartItem != nil {
		// Update quantity
		newQuantity := existingCartItem.Quantity + req.Quantity
		if item.Quantity < newQuantity {
			return nil, errors.New("insufficient quantity available")
		}
		existingCartItem.Quantity = newQuantity
		if err := s.cartRepo.UpdateItem(existingCartItem); err != nil {
			return nil, err
		}
	} else {
		// Add new item
		cartItem := &models.CartItem{
			CartID:   cart.ID,
			ItemID:   req.ItemID,
			Quantity: req.Quantity,
		}
		if err := s.cartRepo.AddItem(cartItem); err != nil {
			return nil, err
		}
	}

	// Reload cart with updated items
	return s.cartRepo.FindByID(cart.ID)
}

// UpdateItemQuantity updates the quantity of an item in the cart
func (s *CartService) UpdateItemQuantity(userID, cartItemID uuid.UUID, req *UpdateCartItemRequest) (*models.Cart, error) {
	// Get user's active cart
	cart, err := s.cartRepo.FindActiveByUserID(userID)
	if err != nil {
		return nil, err
	}
	if cart == nil {
		return nil, errors.New("no active cart found")
	}

	// Find the cart item
	var cartItem *models.CartItem
	for i := range cart.Items {
		if cart.Items[i].ID == cartItemID {
			cartItem = &cart.Items[i]
			break
		}
	}

	if cartItem == nil {
		return nil, errors.New("cart item not found")
	}

	// If quantity is 0, remove the item
	if req.Quantity == 0 {
		if err := s.cartRepo.RemoveItem(cartItemID); err != nil {
			return nil, err
		}
	} else {
		// Verify item availability
		item, err := s.itemRepo.FindByID(cartItem.ItemID)
		if err != nil {
			return nil, err
		}

		if item.Quantity < req.Quantity {
			return nil, errors.New("insufficient quantity available")
		}

		cartItem.Quantity = req.Quantity
		if err := s.cartRepo.UpdateItem(cartItem); err != nil {
			return nil, err
		}
	}

	// Reload cart
	return s.cartRepo.FindByID(cart.ID)
}

// RemoveItem removes an item from the cart
func (s *CartService) RemoveItem(userID, cartItemID uuid.UUID) (*models.Cart, error) {
	// Get user's active cart
	cart, err := s.cartRepo.FindActiveByUserID(userID)
	if err != nil {
		return nil, err
	}
	if cart == nil {
		return nil, errors.New("no active cart found")
	}

	// Verify cart item belongs to this cart
	var found bool
	for _, item := range cart.Items {
		if item.ID == cartItemID {
			found = true
			break
		}
	}

	if !found {
		return nil, errors.New("cart item not found in your cart")
	}

	if err := s.cartRepo.RemoveItem(cartItemID); err != nil {
		return nil, err
	}

	// Reload cart
	return s.cartRepo.FindByID(cart.ID)
}

// ClearCart removes all items from the cart
func (s *CartService) ClearCart(userID uuid.UUID) error {
	cart, err := s.cartRepo.FindActiveByUserID(userID)
	if err != nil {
		return err
	}
	if cart == nil {
		return errors.New("no active cart found")
	}

	return s.cartRepo.ClearCart(cart.ID)
}

// Checkout converts the cart to an order
func (s *CartService) Checkout(userID uuid.UUID, notes string) (*models.Order, error) {
	// Get user's active cart
	cart, err := s.cartRepo.FindActiveByUserID(userID)
	if err != nil {
		return nil, err
	}
	if cart == nil {
		return nil, errors.New("no active cart found")
	}

	if len(cart.Items) == 0 {
		return nil, errors.New("cart is empty")
	}

	// Verify all items are still available
	for _, cartItem := range cart.Items {
		item, err := s.itemRepo.FindByID(cartItem.ItemID)
		if err != nil {
			return nil, errors.New("item not found: " + cartItem.Item.Name)
		}

		if !item.IsAvailable {
			return nil, errors.New("item no longer available: " + item.Name)
		}

		if item.Quantity < cartItem.Quantity {
			return nil, errors.New("insufficient quantity for: " + item.Name)
		}
	}

	// Create order
	order := &models.Order{
		CartID:   cart.ID,
		UserID:   userID,
		PantryID: cart.PantryID,
		Status:   models.OrderStatusPending,
		Notes:    notes,
	}

	// Update cart status
	cart.Status = models.CartStatusSubmitted
	if err := s.cartRepo.Update(cart); err != nil {
		return nil, err
	}

	return order, nil
}
