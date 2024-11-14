package model

import (
	"test-task/shared/database"

	_ "gorm.io/driver/postgres"
)

func AutoMigrate() {
	conn := database.NewConnection()

	conn.GetDB().AutoMigrate(
		// For auto migrate database tables, need to add model below
		&User{},
		&UserRefreshToken{},
	)

}
