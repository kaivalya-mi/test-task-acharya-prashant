package v1ORM

import (
	"database/sql"
	"test-task/model"
	"test-task/shared/database"
	"test-task/shared/log"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

type IUserRepository interface {
	CreateUser(conn database.IConnection, request *model.User) error
	GetUserByEmail(conn database.IConnection, email string) (*model.User, error)
	GetUserById(conn database.IConnection, userID uuid.UUID) (*model.User, error)
}

type userRepo struct {
	DB *sql.DB
}

func NewUserWriter() IUserRepository {
	return &userRepo{}
}

func (ar *userRepo) CreateUser(conn database.IConnection, request *model.User) error {
	log.GetLog().Info("INFO : ", "User Repo Called(CreateUser).")

	result := conn.GetDB().Create(request)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (ar *userRepo) GetUserByEmail(conn database.IConnection, email string) (*model.User, error) {
	log.GetLog().Info("INFO:", "User Repo Called (GetUserByEmail).")

	var user model.User
	err := conn.GetDB().Where("email = ?", email).First(&user).Error

	// Handle error
	if err != nil {
		// If error is "no rows found", return nil without logging
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		log.GetLog().Info("ERROR (query):", err.Error())
		return nil, err
	}

	// Return the user if found
	return &user, nil
}

func (ar *userRepo) GetUserById(conn database.IConnection, userID uuid.UUID) (*model.User, error) {
	log.GetLog().Info("INFO:", "User Repo Called (GetUserByEmail).")

	var user model.User

	// Use First to find by primary key ID
	result := conn.GetDB().First(&user, "id = ?", userID)
	// Handle error
	if result.Error != nil {
		// Return the user if found
		return &user, result.Error
	}
	return &user, nil
}
