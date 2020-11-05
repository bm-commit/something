package login

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

// Command ...
type Command struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

// Validate ...
func (c Command) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Email, validation.Required, is.Email),
		validation.Field(&c.Password, validation.Required),
	)
}
