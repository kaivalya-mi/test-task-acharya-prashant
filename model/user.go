package model

import (
	"time"

	uuid "github.com/satori/go.uuid"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// User represents the user table in the database
type User struct {
	ID        uuid.UUID `gorm:"type:varchar(50);primaryKey" json:"id"`
	FirstName string    `gorm:"varchar(30)" json:"first_name" validate:"required"`
	LastName  string    `gorm:"varchar(30)" json:"last_name" validate:"required"`
	Email     string    `gorm:"unique;not null" json:"email"`
	Password  string    `gorm:"not null" json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName returns the table name for the User model
func (u *User) TableName() string {
	return "users"
}

func (u *User) TimeStamp() {
	u.CreatedAt = time.Now()
}
