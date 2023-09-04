package uuid

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrInvalidUUID = errors.New("invalid uuid")
)

func Generate() string {
	return uuid.New().String()
}

func Validate(value string) bool {
	if _, err := uuid.Parse(value); err != nil {
		return false
	}

	return true
}
