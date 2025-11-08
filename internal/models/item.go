package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Item represents an inventory item
type Item struct {
	ID                 uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name               string     `gorm:"not null" json:"name"`
	Description        string     `json:"description"`
	CategoryID         uuid.UUID  `gorm:"type:uuid;not null" json:"category_id"`
	Category           Category   `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	PantryID           uuid.UUID  `gorm:"type:uuid;not null" json:"pantry_id"`
	Pantry             Pantry     `gorm:"foreignKey:PantryID" json:"pantry,omitempty"`
	Quantity           int        `gorm:"not null;default:0" json:"quantity"`
	LowStockThreshold  int        `gorm:"not null;default:10" json:"low_stock_threshold"`
	Unit               string     `gorm:"not null;default:'count'" json:"unit"` // e.g., "lb", "oz", "count"
	ImageURL           string     `json:"image_url"`
	IsAvailable        bool       `gorm:"default:true" json:"is_available"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (i *Item) BeforeCreate(tx *gorm.DB) error {
	if i.ID == uuid.Nil {
		i.ID = uuid.New()
	}
	return nil
}

// IsLowStock checks if the item is below the low stock threshold
func (i *Item) IsLowStock() bool {
	return i.Quantity <= i.LowStockThreshold
}
