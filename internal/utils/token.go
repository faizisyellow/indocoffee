package utils

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/google/uuid"
)

func GenerateTokenUuid() string {

	plaintoken := uuid.New().String()

	hash := sha256.Sum256([]byte(plaintoken))

	return hex.EncodeToString(hash[:])

}
