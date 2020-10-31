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
	Role      string
	CreatedOn time.Time
}

// NewUser ...
func NewUser(id, name, username, email, password string) (*User, error) {
	return &User{
		ID:        id,
		Name:      name,
		Username:  strings.TrimSpace(strings.ToLower(username)),
		Email:     strings.TrimSpace(strings.ToLower(email)),
		Password:  password,
		Role:      "default",
		CreatedOn: time.Now(),
	}, nil
}
