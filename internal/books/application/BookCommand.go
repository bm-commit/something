package application

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

// BookCommand ...
type BookCommand struct {
	ID          string `json:"id"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Author      string `json:"author,omitempty"`
	Genre       string `json:"genre,omitempty"`
	Pages       int    `json:"pages,omitempty"`
}

// Validate ...
func (b BookCommand) Validate() error {
	return validation.ValidateStruct(&b,
		validation.Field(&b.Title, validation.Required, validation.Length(1, 75)),
		validation.Field(&b.Description, validation.Required, validation.Length(1, 1500)),
		validation.Field(&b.Author, validation.Length(1, 75)),
		validation.Field(&b.Genre, validation.Required, validation.Length(1, 150)),
		validation.Field(&b.Pages, validation.Min(1)),
	)
}
