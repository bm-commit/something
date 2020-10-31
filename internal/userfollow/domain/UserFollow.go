package domain

import (
	"time"
)

// UserFollow ...
type UserFollow struct {
	From      string
	To        string
	CreatedOn time.Time
}

// NewUserFollow ...
func NewUserFollow(from, to string) (*UserFollow, error) {
	return &UserFollow{
		From:      from,
		To:        to,
		CreatedOn: time.Now(),
	}, nil
}
