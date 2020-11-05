package update

import (
	"encoding/json"
	"errors"
	"strings"

	"something/internal/bookreviews/domain"
)

// Service ...
type Service interface {
	UpdateBookReviewByID(*BookReviewCommand) error
}

type service struct {
	repository domain.BookReviewRepository
}

// NewService ...
func NewService(repository domain.BookReviewRepository) Service {
	return &service{repository: repository}
}

func (s *service) UpdateBookReviewByID(bookReview *BookReviewCommand) error {
	existingBookReview, _ := s.repository.FindByID(bookReview.ID)
	if existingBookReview == nil {
		return errors.New("book review not found")
	}
	if existingBookReview.UserID != bookReview.UserID {
		return errors.New("unauthorized")
	}

	out, err := json.Marshal(bookReview)
	if err != nil {
		return err
	}
	// Merge new book fields with existing book fields
	err = json.NewDecoder(strings.NewReader(string(out))).Decode(existingBookReview)
	if err != nil {
		return err
	}

	updatedBookReview, err := domain.NewBookReview(
		existingBookReview.ID, existingBookReview.Text, existingBookReview.Rating,
		existingBookReview.BookID, existingBookReview.UserID)
	if err != nil {
		return err
	}
	updatedBookReview.CreatedOn = existingBookReview.CreatedOn

	err = s.repository.Update(updatedBookReview)
	return err
}
