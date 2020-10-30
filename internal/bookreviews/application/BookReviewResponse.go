package application

import (
	"time"

	"something/internal/bookreviews/domain"
)

// BookReviewResponse ...
type BookReviewResponse struct {
	ID        string    `json:"id"`
	Text      string    `json:"text"`
	Rating    int       `json:"rating"`
	BookID    string    `json:"book_id"`
	UserID    string    `json:"user_id"`
	CreatedOn time.Time `json:"created_on"`
}

// NewBookReviewResponse ...
func NewBookReviewResponse(bookReview *domain.BookReview) *BookReviewResponse {
	return &BookReviewResponse{
		ID:        bookReview.ID,
		Text:      bookReview.Text,
		Rating:    bookReview.Rating,
		BookID:    bookReview.BookID,
		UserID:    bookReview.UserID,
		CreatedOn: bookReview.CreatedOn,
	}
}

// NewReviewsResponse ...
func NewReviewsResponse(bookReviews []*domain.BookReview) []*BookReviewResponse {
	var bookReviewsResponse []*BookReviewResponse
	for _, review := range bookReviews {
		bookReviewsResponse = append(bookReviewsResponse, NewBookReviewResponse(review))
	}
	return bookReviewsResponse
}
