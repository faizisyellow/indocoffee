package models

import (
	"time"

	"github.com/faizisyellow/indocoffee/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  Hashed    `json:"-"`
	IsActive  *bool     `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

type Hashed struct {
	Text       *string
	HashedText []byte
}

// Parse parses password to bcrypt hash.
// Parsed password will be saved to HashedText
// and, keep the original plaintext password to text.
func (h *Hashed) ParseFromPassword(password string) error {

	hashed, err := utils.GeneratePassword(password)
	if err != nil {
		return err
	}

	h.Text = &password
	h.HashedText = hashed

	return nil
}

// ComparePassword compares plaintext password to its hashed password.
// Returns nil on success, or an error on failure.
func (h *Hashed) ComparePassword(password string) error {

	return bcrypt.CompareHashAndPassword(h.HashedText, []byte(password))
}
