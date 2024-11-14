package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// HashedPassword() is Password generator
func HashedPassword(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("ERROR : ", err.Error())
		return ""
	}
	return string(hashedPassword)
}
