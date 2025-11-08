package repositories

import (
	"errors"

	"github.com/byte4bite/byte4bite/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CartRepository handles database operations for carts
type CartRepository struct {
	db *gorm.DB
}

// NewCartRepository creates a new cart repository
func NewCartRepository(db *gorm.DB) *CartRepository {
	return &CartRepository{db: db}
}

// Create creates a new cart
func (r *CartRepository) Create(cart *models.Cart) error {
	return r.db.Create(cart).Error
}

// FindByID finds a cart by ID
func (r *CartRepository) FindByID(id uuid.UUID) (*models.Cart, error) {
	var cart models.Cart
	err := r.db.Preload("Items.Item.Category").Preload("User").Preload("Pantry").
		First(&cart, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("cart not found")
		}
		return nil, err
	}
	return &cart, nil
}

// FindActiveByUserID finds the active cart for a user
func (r *CartRepository) FindActiveByUserID(userID uuid.UUID) (*models.Cart, error) {
	var cart models.Cart
	err := r.db.Preload("Items.Item.Category").Preload("User").Preload("Pantry").
		Where("user_id = ? AND status = ?", userID, models.CartStatusActive).
		First(&cart).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // No active cart is not an error
		}
		return nil, err
	}
	return &cart, nil
}

// Update updates a cart
func (r *CartRepository) Update(cart *models.Cart) error {
	return r.db.Save(cart).Error
}

// Delete deletes a cart
func (r *CartRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Cart{}, "id = ?", id).Error
}

// AddItem adds an item to a cart
func (r *CartRepository) AddItem(cartItem *models.CartItem) error {
	return r.db.Create(cartItem).Error
}

// UpdateItem updates a cart item
func (r *CartRepository) UpdateItem(cartItem *models.CartItem) error {
	return r.db.Save(cartItem).Error
}

// RemoveItem removes an item from a cart
func (r *CartRepository) RemoveItem(cartItemID uuid.UUID) error {
	return r.db.Delete(&models.CartItem{}, "id = ?", cartItemID).Error
}

// FindCartItem finds a specific cart item
func (r *CartRepository) FindCartItem(cartID, itemID uuid.UUID) (*models.CartItem, error) {
	var cartItem models.CartItem
	err := r.db.Where("cart_id = ? AND item_id = ?", cartID, itemID).First(&cartItem).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &cartItem, nil
}

// ClearCart removes all items from a cart
func (r *CartRepository) ClearCart(cartID uuid.UUID) error {
	return r.db.Where("cart_id = ?", cartID).Delete(&models.CartItem{}).Error
}

// GetCartItemCount returns the number of items in a cart
func (r *CartRepository) GetCartItemCount(cartID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&models.CartItem{}).Where("cart_id = ?", cartID).Count(&count).Error
	return count, err
}

// FindCartsByUserID finds all carts for a user
func (r *CartRepository) FindCartsByUserID(userID uuid.UUID, limit, offset int) ([]models.Cart, error) {
	var carts []models.Cart
	err := r.db.Preload("Items.Item").Preload("Pantry").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&carts).Error
	return carts, err
}
