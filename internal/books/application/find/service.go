package find

import (
	"something/internal/books/application"
	"something/internal/books/domain"
)

// PAGE Default pagination page
const PAGE int = 1

// PERPAGE Default page size (the number of items to return per page).
const PERPAGE int = 50

// Service ...
type Service interface {
	FindBooks(criteria *Criteria) ([]*application.BookResponse, error)
	FindBookByID(id string) (*application.BookResponse, error)
}

type service struct {
	repository domain.BookRepository
}

// NewService ...
func NewService(repository domain.BookRepository) Service {
	return &service{repository: repository}
}

func (s *service) FindBooks(criteria *Criteria) ([]*application.BookResponse, error) {

	if criteria.Page == 0 {
		criteria.Page = PAGE
	}
	if criteria.PerPage == 0 || criteria.PerPage > 1000 {
		criteria.PerPage = PERPAGE
	}

	newBookCriteria := domain.NewBookCriteria(
		criteria.Page, criteria.PerPage, criteria.Query,
		criteria.Genre, criteria.Author,
	)

	books, err := s.repository.Find(newBookCriteria)
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
