package application

import (
	"time"

	"something/internal/bookreviews/domain"
)

// BookReviewResponse ...
type BookReviewResponse struct {
	ID        string  `json:"id"`
	Text      string  `json:"text"`
	Rating    float64 `json:"rating"`
	BookID    string  `json:"book_id"`
	User      `json:"user"`
	CreatedOn time.Time `json:"created_on"`
}

// BookRatingResponse ...
type BookRatingResponse struct {
	BookID string  `json:"book_id"`
	Title  string  `json:"title"`
	Author string  `json:"author"`
	Rating float64 `json:"rating"`
	Total  int     `json:"total"`
}

// User ...
type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

// NewBookReviewResponse ...
func NewBookReviewResponse(bookReview *domain.BookReview) *BookReviewResponse {
	return &BookReviewResponse{
		ID:        bookReview.ID,
		Text:      bookReview.Text,
		Rating:    bookReview.Rating,
		BookID:    bookReview.BookID,
		User:      User{ID: bookReview.UserID},
		CreatedOn: bookReview.CreatedOn,
	}
}

// NewReviewsResponse ...
func NewReviewsResponse(bookReviews []*domain.BookReview) []*BookReviewResponse {
	bookReviewsResponse := []*BookReviewResponse{}
	for _, review := range bookReviews {
		bookReviewsResponse = append(bookReviewsResponse, NewBookReviewResponse(review))
	}
	return bookReviewsResponse
}

// NewBookReviewShortResponse ...
func NewBookReviewShortResponse(bookReview *domain.BookReviewShort) *BookRatingResponse {
	return &BookRatingResponse{
		Rating: bookReview.Rating,
		BookID: bookReview.ID,
		Total:  bookReview.Total,
	}
}

// NewReviewShortResponse ...
func NewReviewShortResponse(bookReviews []*domain.BookReviewShort) []*BookRatingResponse {
	bookReviewsResponse := []*BookRatingResponse{}
	for _, review := range bookReviews {
		bookReviewsResponse = append(bookReviewsResponse, NewBookReviewShortResponse(review))
	}
	return bookReviewsResponse
}
