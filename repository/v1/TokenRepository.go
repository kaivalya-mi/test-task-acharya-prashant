package v1ORM

import (
	"database/sql"
	"fmt"
	"test-task/model"
	"test-task/shared/database"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

type ITokenRepository interface {
	FindTokenData(conn database.IConnection, userID uuid.UUID, token string) (*model.UserRefreshToken, error)
	SaveRefreshToken(conn database.IConnection, refreshToken *model.UserRefreshToken) error
}

type tokenRepo struct {
	DB *sql.DB
}

func NewTokenWriter() ITokenRepository {
	return &tokenRepo{}
}

// FindByUserIDAndToken retrieves a refresh token by userID and token
func (r *tokenRepo) FindTokenData(conn database.IConnection, userID uuid.UUID, token string) (*model.UserRefreshToken, error) {
	var refreshToken model.UserRefreshToken
	if err := conn.GetDB().Where("user_id = ? AND refresh_token = ?", userID, token).First(&refreshToken).Error; err != nil {
		return nil, err
	}
	return &refreshToken, nil
}

// SaveRefreshToken saves or updates a user's refresh token
func (r *tokenRepo) SaveRefreshToken(conn database.IConnection, refreshToken *model.UserRefreshToken) error {

	var existingToken model.UserRefreshToken
	err := conn.GetDB().Where("user_id = ?", refreshToken.UserID).First(&existingToken).Error

	// If the token doesn't exist, create a new one
	if err != nil && err == gorm.ErrRecordNotFound {
		// No token found for the user, create a new one
		refreshToken.ID = uuid.NewV1() // Create a new ID for the refresh token
		err = conn.GetDB().Create(refreshToken).Error
		if err != nil {
			return fmt.Errorf("error creating new refresh token: %v", err)
		}
		return nil
	} else if err != nil {
		// If any other error occurred, return the error
		return fmt.Errorf("error checking for existing refresh token: %v", err)
	}

	// Here, we simply replace the existing refresh token and expiration
	existingToken.RefreshToken = refreshToken.RefreshToken
	existingToken.ExpiresAt = refreshToken.ExpiresAt

	err = conn.GetDB().Save(&existingToken).Error
	if err != nil {
		return fmt.Errorf("error updating refresh token: %v", err)
	}

	return nil
}
