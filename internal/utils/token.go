package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

type Token interface {
	Generate() string
}

type UUID struct {
	Plaintoken string
}

func (u UUID) Generate() string {

	hash := sha256.Sum256([]byte(u.Plaintoken))

	return hex.EncodeToString(hash[:])

}
