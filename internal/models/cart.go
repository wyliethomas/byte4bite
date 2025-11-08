package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CartStatus represents the status of a cart
type CartStatus string

const (
	CartStatusActive    CartStatus = "active"
	CartStatusSubmitted CartStatus = "submitted"
	CartStatusCancelled CartStatus = "cancelled"
)

// Cart represents a user's shopping cart
type Cart struct {
	ID        uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	User      User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	PantryID  uuid.UUID  `gorm:"type:uuid;not null" json:"pantry_id"`
	Pantry    Pantry     `gorm:"foreignKey:PantryID" json:"pantry,omitempty"`
	Status    CartStatus `gorm:"type:varchar(20);not null;default:'active'" json:"status"`
	Items     []CartItem `gorm:"foreignKey:CartID" json:"items,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (c *Cart) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

// CartItem represents an item in a shopping cart
type CartItem struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	CartID    uuid.UUID `gorm:"type:uuid;not null" json:"cart_id"`
	Cart      Cart      `gorm:"foreignKey:CartID" json:"cart,omitempty"`
	ItemID    uuid.UUID `gorm:"type:uuid;not null" json:"item_id"`
	Item      Item      `gorm:"foreignKey:ItemID" json:"item,omitempty"`
	Quantity  int       `gorm:"not null;default:1" json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (ci *CartItem) BeforeCreate(tx *gorm.DB) error {
	if ci.ID == uuid.Nil {
		ci.ID = uuid.New()
	}
	return nil
}
