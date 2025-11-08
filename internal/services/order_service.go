package services

import (
	"errors"
	"time"

	"github.com/byte4bite/byte4bite/internal/models"
	"github.com/byte4bite/byte4bite/internal/repositories"
	"github.com/google/uuid"
)

// OrderService handles business logic for orders
type OrderService struct {
	orderRepo *repositories.OrderRepository
	itemRepo  *repositories.ItemRepository
}

// NewOrderService creates a new order service
func NewOrderService(orderRepo *repositories.OrderRepository, itemRepo *repositories.ItemRepository) *OrderService {
	return &OrderService{
		orderRepo: orderRepo,
		itemRepo:  itemRepo,
	}
}

// GetOrderRequest represents the request to get orders
type GetOrdersRequest struct {
	Status   *models.OrderStatus
	Page     int
	PageSize int
}

// GetOrdersResponse represents the response containing orders
type GetOrdersResponse struct {
	Orders []models.Order `json:"orders"`
	Total  int64          `json:"total"`
	Page   int            `json:"page"`
	Pages  int            `json:"pages"`
}

// UpdateStatusRequest represents a request to update order status
type UpdateStatusRequest struct {
	Status models.OrderStatus `json:"status" binding:"required"`
}

// AssignStaffRequest represents a request to assign staff to an order
type AssignStaffRequest struct {
	StaffID uuid.UUID `json:"staff_id" binding:"required"`
}

// GetOrders returns a list of orders
func (s *OrderService) GetOrders(userID uuid.UUID, isAdmin bool, req GetOrdersRequest) (*GetOrdersResponse, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 100 {
		req.PageSize = 20
	}

	offset := (req.Page - 1) * req.PageSize

	var orders []models.Order
	var total int64
	var err error

	if isAdmin {
		// Admin can see all orders
		orders, err = s.orderRepo.FindAll(req.Status, req.PageSize, offset)
		if err != nil {
			return nil, err
		}
		total, err = s.orderRepo.CountAll(req.Status)
	} else {
		// Users can only see their own orders
		orders, err = s.orderRepo.FindByUserID(userID, req.PageSize, offset)
		if err != nil {
			return nil, err
		}
		total, err = s.orderRepo.CountByUserID(userID)
	}

	if err != nil {
		return nil, err
	}

	pages := int(total) / req.PageSize
	if int(total)%req.PageSize != 0 {
		pages++
	}

	return &GetOrdersResponse{
		Orders: orders,
		Total:  total,
		Page:   req.Page,
		Pages:  pages,
	}, nil
}

// GetOrder returns a single order by ID
func (s *OrderService) GetOrder(orderID uuid.UUID, userID uuid.UUID, isAdmin bool) (*models.Order, error) {
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return nil, err
	}

	// Check permissions: users can only view their own orders
	if !isAdmin && order.UserID != userID {
		return nil, errors.New("unauthorized to view this order")
	}

	return order, nil
}

// UpdateOrderStatus updates the status of an order with validation
func (s *OrderService) UpdateOrderStatus(orderID uuid.UUID, newStatus models.OrderStatus) error {
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return err
	}

	// Validate status transition
	if !isValidStatusTransition(order.Status, newStatus) {
		return errors.New("invalid status transition")
	}

	// Update the status
	order.Status = newStatus

	// Set timestamps based on status
	now := time.Now()
	switch newStatus {
	case models.OrderStatusReady:
		order.ReadyAt = &now
	case models.OrderStatusPickedUp:
		order.PickedUpAt = &now
	}

	return s.orderRepo.Update(order)
}

// AssignStaff assigns an order to a staff member
func (s *OrderService) AssignStaff(orderID uuid.UUID, staffID uuid.UUID) error {
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return err
	}

	// Can't assign cancelled or picked up orders
	if order.Status == models.OrderStatusCancelled || order.Status == models.OrderStatusPickedUp {
		return errors.New("cannot assign staff to cancelled or completed orders")
	}

	order.AssignedToID = &staffID
	return s.orderRepo.Update(order)
}

// CancelOrder cancels an order
func (s *OrderService) CancelOrder(orderID uuid.UUID, userID uuid.UUID, isAdmin bool) error {
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return err
	}

	// Check permissions
	if !isAdmin && order.UserID != userID {
		return errors.New("unauthorized to cancel this order")
	}

	// Can only cancel pending or preparing orders
	if order.Status != models.OrderStatusPending && order.Status != models.OrderStatusPreparing {
		return errors.New("can only cancel pending or preparing orders")
	}

	// Restore inventory when cancelling
	if len(order.Cart.Items) > 0 {
		for _, cartItem := range order.Cart.Items {
			item, err := s.itemRepo.FindByID(cartItem.ItemID)
			if err != nil {
				continue // Skip if item not found
			}
			item.Quantity += cartItem.Quantity
			s.itemRepo.Update(item)
		}
	}

	order.Status = models.OrderStatusCancelled
	return s.orderRepo.Update(order)
}

// isValidStatusTransition checks if a status transition is valid
func isValidStatusTransition(from, to models.OrderStatus) bool {
	validTransitions := map[models.OrderStatus][]models.OrderStatus{
		models.OrderStatusPending: {
			models.OrderStatusPreparing,
			models.OrderStatusCancelled,
		},
		models.OrderStatusPreparing: {
			models.OrderStatusReady,
			models.OrderStatusCancelled,
		},
		models.OrderStatusReady: {
			models.OrderStatusPickedUp,
			models.OrderStatusCancelled,
		},
		models.OrderStatusPickedUp: {}, // Final state
		models.OrderStatusCancelled: {}, // Final state
	}

	allowedTransitions, exists := validTransitions[from]
	if !exists {
		return false
	}

	for _, allowed := range allowedTransitions {
		if allowed == to {
			return true
		}
	}

	return false
}
