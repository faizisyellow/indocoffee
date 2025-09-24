package utils

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/google/uuid"
)

type Token interface {
	Generate() string
}

type UUID struct {
	Plaintoken string
}

func (u UUID) Generate() string {

	if u.Plaintoken == "" {
		u.Plaintoken = uuid.New().String()
	}

	hash := sha256.Sum256([]byte(u.Plaintoken))

	return hex.EncodeToString(hash[:])

}
