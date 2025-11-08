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

// FindAll finds all orders with optional filtering
func (r *OrderRepository) FindAll(status *models.OrderStatus, limit, offset int) ([]models.Order, error) {
	var orders []models.Order
	query := r.db.Preload("Cart.Items.Item").Preload("User").Preload("Pantry").Preload("AssignedTo")

	if status != nil {
		query = query.Where("status = ?", *status)
	}

	err := query.Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&orders).Error
	return orders, err
}

// CountAll counts all orders with optional status filter
func (r *OrderRepository) CountAll(status *models.OrderStatus) (int64, error) {
	var count int64
	query := r.db.Model(&models.Order{})

	if status != nil {
		query = query.Where("status = ?", *status)
	}

	err := query.Count(&count).Error
	return count, err
}

// UpdateStatus updates the status of an order
func (r *OrderRepository) UpdateStatus(id uuid.UUID, status models.OrderStatus) error {
	return r.db.Model(&models.Order{}).Where("id = ?", id).
		Update("status", status).Error
}

// AssignToStaff assigns an order to a staff member
func (r *OrderRepository) AssignToStaff(id uuid.UUID, staffID uuid.UUID) error {
	return r.db.Model(&models.Order{}).Where("id = ?", id).
		Update("assigned_to_id", staffID).Error
}

// Delete soft deletes an order (for cancellation)
func (r *OrderRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Order{}, "id = ?", id).Error
}
