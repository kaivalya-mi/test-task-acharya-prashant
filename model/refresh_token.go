package model

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// UserRefreshToken represents the refresh token stored in the database
type UserRefreshToken struct {
	ID           uuid.UUID `gorm:"type:varchar(50);primaryKey" json:"id"`
	UserID       uuid.UUID `gorm:"type:uuid;index;not null;" json:"user_id"`
	RefreshToken string    `gorm:"type:text;unique;not null;" json:"refresh_token"`
	ExpiresAt    time.Time `gorm:"not null;" json:"expires_at"`
	CreatedAt    time.Time `gorm:"autoCreateTime;" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime;" json:"updated_at"`
}

// TableName returns the table name for the UserRefreshToken model
func (u *UserRefreshToken) TableName() string {
	return "user_refresh_tokens"
}
