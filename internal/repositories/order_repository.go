package repositories

import (
	"errors"

	"github.com/byte4bite/byte4bite/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// OrderRepository handles database operations for orders
type OrderRepository struct {
	db *gorm.DB
}

// NewOrderRepository creates a new order repository
func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

// Create creates a new order
func (r *OrderRepository) Create(order *models.Order) error {
	return r.db.Create(order).Error
}

// FindByID finds an order by ID
func (r *OrderRepository) FindByID(id uuid.UUID) (*models.Order, error) {
	var order models.Order
	err := r.db.Preload("Cart.Items.Item").Preload("User").Preload("Pantry").
		Preload("AssignedTo").First(&order, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("order not found")
		}
		return nil, err
	}
	return &order, nil
}

// Update updates an order
func (r *OrderRepository) Update(order *models.Order) error {
	return r.db.Save(order).Error
}

// FindByUserID finds all orders for a user
func (r *OrderRepository) FindByUserID(userID uuid.UUID, limit, offset int) ([]models.Order, error) {
	var orders []models.Order
	err := r.db.Preload("Cart.Items.Item").Preload("Pantry").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&orders).Error
	return orders, err
}

// FindByPantryID finds all orders for a pantry
func (r *OrderRepository) FindByPantryID(pantryID uuid.UUID, limit, offset int) ([]models.Order, error) {
	var orders []models.Order
	err := r.db.Preload("Cart.Items.Item").Preload("User").
		Where("pantry_id = ?", pantryID).
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&orders).Error
	return orders, err
}

// CountByUserID counts orders for a user
func (r *OrderRepository) CountByUserID(userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&models.Order{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

// CountByPantryID counts orders for a pantry
func (r *OrderRepository) CountByPantryID(pantryID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&models.Order{}).Where("pantry_id = ?", pantryID).Count(&count).Error
	return count, err
}
