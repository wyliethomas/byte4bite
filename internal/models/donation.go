package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Donation represents a donation to a pantry
type Donation struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	PantryID     uuid.UUID `gorm:"type:uuid;not null" json:"pantry_id"`
	Pantry       Pantry    `gorm:"foreignKey:PantryID" json:"pantry,omitempty"`
	DonorName    string    `gorm:"not null" json:"donor_name"`
	DonorEmail   string    `json:"donor_email"`
	DonorPhone   string    `json:"donor_phone"`
	Amount       *float64  `json:"amount"` // For monetary donations
	Description  string    `gorm:"not null" json:"description"`
	DonationDate time.Time `gorm:"not null" json:"donation_date"`
	ReceiptSent  bool      `gorm:"default:false" json:"receipt_sent"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (d *Donation) BeforeCreate(tx *gorm.DB) error {
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	return nil
}
