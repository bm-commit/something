package domain

import "time"

// BookReview ...
type BookReview struct {
	ID        string
	Text      string
	Rating    float64
	BookID    string
	UserID    string
	CreatedOn time.Time
}

// NewBookReview ...
func NewBookReview(id, text string, rating float64, bookID, userID string) (*BookReview, error) {
	return &BookReview{
		ID:        id,
		Text:      text,
		Rating:    rating,
		BookID:    bookID,
		UserID:    userID,
		CreatedOn: time.Now().UTC(),
	}, nil
}
