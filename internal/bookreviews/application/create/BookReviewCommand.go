package create

import validation "github.com/go-ozzo/ozzo-validation"

// BookReviewCommand ...
type BookReviewCommand struct {
	ID     string  `json:"id"`
	Text   string  `json:"text"`
	Rating float64 `json:"rating"`
	BookID string  `json:"book_id"`
	UserID string  `json:"user_id"`
}

// Validate ...
func (b BookReviewCommand) Validate() error {
	return validation.ValidateStruct(&b,
		validation.Field(&b.Text, validation.Required, validation.Length(1, 250)),
		validation.Field(&b.Rating, validation.Required, validation.Min(0.5), validation.Max(5.0)),
	)
}
