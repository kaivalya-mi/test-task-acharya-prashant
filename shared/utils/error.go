package utils

import (
	"errors"
	"fmt"
)

func NewError(format string, a ...interface{}) error {
	msg := fmt.Sprintf(format, a...)
	return errors.New(msg)
}
