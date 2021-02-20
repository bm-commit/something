package update

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

// UserCommand ...
type UserCommand struct {
	ID       string `json:"id"`
	Name     string `json:"name,omitempty"`
	Username string `json:"username,omitempty"`
}

// Validate ...
func (u UserCommand) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Name, validation.Length(1, 45)),
		validation.Field(&u.Username, validation.Length(1, 45)),
	)
}
