package update

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

// UserInterestsCommand ...
type UserInterestsCommand struct {
	UserID string `json:"user_id"`
	BookID string `json:"book_id"`
	Status string `json:"status"`
}

// Validate ...
func (u UserInterestsCommand) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Status, validation.Required, validation.In("reading", "pending", "done")),
	)
}
