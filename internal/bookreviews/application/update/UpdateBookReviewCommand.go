package update

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

// BookReviewCommand ...
type BookReviewCommand struct {
	ID     string `json:"id"`
	Text   string `json:"text"`
	UserID string `json:"user_id"`
}

// Validate ...
func (b BookReviewCommand) Validate() error {
	return validation.ValidateStruct(&b,
		validation.Field(&b.Text, validation.Required, validation.Length(1, 250)),
	)
}
