package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// OrderStatus represents the status of an order
type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusPreparing OrderStatus = "preparing"
	OrderStatusReady     OrderStatus = "ready"
	OrderStatusPickedUp  OrderStatus = "picked_up"
	OrderStatusCancelled OrderStatus = "cancelled"
)

// Order represents a submitted cart order
type Order struct {
	ID           uuid.UUID    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	CartID       uuid.UUID    `gorm:"type:uuid;not null" json:"cart_id"`
	Cart         Cart         `gorm:"foreignKey:CartID" json:"cart,omitempty"`
	UserID       uuid.UUID    `gorm:"type:uuid;not null" json:"user_id"`
	User         User         `gorm:"foreignKey:UserID" json:"user,omitempty"`
	PantryID     uuid.UUID    `gorm:"type:uuid;not null" json:"pantry_id"`
	Pantry       Pantry       `gorm:"foreignKey:PantryID" json:"pantry,omitempty"`
	Status       OrderStatus  `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`
	Notes        string       `json:"notes"`
	AssignedToID *uuid.UUID   `gorm:"type:uuid" json:"assigned_to_id"`
	AssignedTo   *User        `gorm:"foreignKey:AssignedToID" json:"assigned_to,omitempty"`
	SubmittedAt  time.Time    `gorm:"not null" json:"submitted_at"`
	ReadyAt      *time.Time   `json:"ready_at"`
	PickedUpAt   *time.Time   `json:"picked_up_at"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (o *Order) BeforeCreate(tx *gorm.DB) error {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	if o.SubmittedAt.IsZero() {
		o.SubmittedAt = time.Now()
	}
	return nil
}
