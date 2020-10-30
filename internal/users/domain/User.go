package domain

import (
	"strings"
	"time"
)

// User ...
type User struct {
	ID        string
	Name      string
	Username  string
	Email     string
	Password  string
	IsAdmin   bool
	CreatedOn time.Time
}

// NewUser ...
func NewUser(id, name, username, email, password string) (*User, error) {
	return &User{
		ID:        id,
		Name:      name,
		Username:  username, // strings.TrimSpace(strings.ToLower(username))
		Email:     strings.ToLower(email),
		Password:  password,
		IsAdmin:   false,
		CreatedOn: time.Now(),
	}, nil
}
