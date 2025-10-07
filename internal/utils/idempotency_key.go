package utils

import (
	"github.com/google/uuid"
)

func ValidateIdempotencyKey(key string) bool {
	parsed, err := uuid.Parse(key)
	if err != nil || parsed.Version() != 4 {
		return false
	}
	return true
}
