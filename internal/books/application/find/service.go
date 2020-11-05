package find

import (
	"something/internal/books/application"
	"something/internal/books/domain"
)

// Service ...
type Service interface {
	FindBooks() ([]*application.BookResponse, error)
	FindBookByID(id string) (*application.BookResponse, error)
}

type service struct {
	repository domain.BookRepository
}

// NewService ...
func NewService(repository domain.BookRepository) Service {
	return &service{repository: repository}
}

func (s *service) FindBooks() ([]*application.BookResponse, error) {
	books, err := s.repository.Find()
	if err != nil {
		return nil, err
	}
	return application.NewBooksResponse(books), nil
}

func (s *service) FindBookByID(id string) (*application.BookResponse, error) {
	book, err := s.repository.FindByID(id)
	if err != nil {
		return nil, err
	}
	return application.NewBookResponse(book), nil
}
