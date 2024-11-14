package crypto

import (
	"test-task/shared/utils/middleware"
	"time"

	uuid "github.com/satori/go.uuid" // Importing the satori UUID package
)

type UserTokenData struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// GenerateAuthToken creates a JWT token with the user's UUID and creation timestamp
func GenerateAuthToken(userID uuid.UUID, email string, createdAt time.Time) (string, error) {
	tokenData := &UserTokenData{
		ID:        userID.String(), // Store UUID as a string
		Email:     email,
		CreatedAt: createdAt,
	}

	token, err := middleware.GenerateToken(tokenData)
	if err != nil {
		return "", err
	}
	return token, nil
}
