package delete

import (
	"errors"
	"something/internal/books/domain"
)

// Service ...
type Service interface {
	DeleteBookByID(id string) error
}

type service struct {
	repository domain.BookRepository
}

// NewService ...
func NewService(repository domain.BookRepository) Service {
	return &service{repository: repository}
}

func (s *service) DeleteBookByID(id string) error {
	review, _ := s.repository.FindByID(id)
	if review == nil {
		return errors.New("book not found")
	}
	err := s.repository.Delete(id)
	return err
}
