package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Pantry represents a community pantry
type Pantry struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name         string    `gorm:"not null" json:"name"`
	Address      string    `gorm:"not null" json:"address"`
	City         string    `gorm:"not null" json:"city"`
	State        string    `gorm:"not null" json:"state"`
	ZipCode      string    `gorm:"not null" json:"zip_code"`
	ContactEmail string    `gorm:"not null" json:"contact_email"`
	ContactPhone string    `json:"contact_phone"`
	IsActive     bool      `gorm:"default:true" json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (p *Pantry) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}
