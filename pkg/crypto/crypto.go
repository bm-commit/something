package crypto

import (
	"golang.org/x/crypto/bcrypt"
)

const defaultSalt = 12

// Crypto ...
type Crypto interface {
	Hash(string) (string, error)
	CompareHashAndText(text, hash string) bool
}

type repository struct {
}

// NewBcrypt ...
func NewBcrypt() Crypto {
	return &repository{}
}

// Hash ...
func (r *repository) Hash(text string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(text), defaultSalt)
	return string(bytes), err
}

// CheckPasswordHash ...
func (r *repository) CompareHashAndText(text, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(text))
	return err == nil
}
