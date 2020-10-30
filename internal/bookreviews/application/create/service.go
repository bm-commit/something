package create

import (
	"errors"
	"something/internal/bookreviews/domain"
)

// Service ...
type Service interface {
	CreateBookReview(*BookReviewCommand) error
}

type service struct {
	repository domain.BookReviewRepository
}

// NewService ...
func NewService(repository domain.BookReviewRepository) Service {
	return &service{repository: repository}
}

func (s *service) CreateBookReview(command *BookReviewCommand) error {
	bookReview, err := domain.NewBookReview(
		command.ID, command.Text, command.Rating, command.BookID, command.UserID)
	if err != nil {
		return err
	}

	existingReviewID, _ := s.repository.FindByID(command.ID)
	if existingReviewID != nil {
		return errors.New("book review id already exists")
	}

	err = s.repository.Save(bookReview)
	if err != nil {
		return err
	}
	return nil
}
