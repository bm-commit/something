package delete

import (
	"errors"
	"something/internal/bookreviews/domain"
)

// Service ...
type Service interface {
	DeleteBookReviewByID(id string) error
}

type service struct {
	repository domain.BookReviewRepository
}

// NewService ...
func NewService(repository domain.BookReviewRepository) Service {
	return &service{repository: repository}
}

func (s *service) DeleteBookReviewByID(id string) error {
	bookReview, _ := s.repository.FindByID(id)
	if bookReview == nil {
		return errors.New("book review not found")
	}
	err := s.repository.Delete(id)
	return err
}
