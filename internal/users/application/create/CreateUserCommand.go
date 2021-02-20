package create

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

// UserCommand ...
type UserCommand struct {
	ID       string `json:"id"`
	Name     string `json:"name,omitempty"`
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

// Validate ...
func (u UserCommand) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Name, validation.Required, validation.Length(1, 45)),
		validation.Field(&u.Username, validation.Required, validation.Length(1, 45)),
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Password, validation.Required, validation.Length(8, 64)),
	)
}
