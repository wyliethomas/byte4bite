package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// NotificationType represents the type of notification
type NotificationType string

const (
	NotificationTypeEmail NotificationType = "email"
	NotificationTypeSMS   NotificationType = "sms"
)

// Notification represents a notification sent to a user
type Notification struct {
	ID        uuid.UUID        `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID        `gorm:"type:uuid;not null" json:"user_id"`
	User      User             `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Type      NotificationType `gorm:"type:varchar(20);not null" json:"type"`
	Subject   string           `json:"subject"`
	Message   string           `gorm:"not null" json:"message"`
	Sent      bool             `gorm:"default:false" json:"sent"`
	SentAt    *time.Time       `json:"sent_at"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (n *Notification) BeforeCreate(tx *gorm.DB) error {
	if n.ID == uuid.Nil {
		n.ID = uuid.New()
	}
	return nil
}
