package find

import (
	"something/internal/bookreviews/application"
	"something/internal/bookreviews/domain"
	bookDomain "something/internal/books/domain"
)

// Service ...
type Service interface {
	FindBookReviews(bookID string) ([]*application.BookReviewResponse, error)
	FindBookReviewByID(id string) (*application.BookReviewResponse, error)
}

type service struct {
	bookRepository bookDomain.BookRepository
	repository     domain.BookReviewRepository
}

// NewService ...
func NewService(repository domain.BookReviewRepository, bookRepository bookDomain.BookRepository) Service {
	return &service{repository: repository, bookRepository: bookRepository}
}

func (s *service) FindBookReviews(bookID string) ([]*application.BookReviewResponse, error) {
	book, err := s.bookRepository.FindByID(bookID)
	if book == nil {
		return nil, err
	}

	bookReviews, err := s.repository.Find(bookID)
	if err != nil {
		return nil, err
	}
	return application.NewReviewsResponse(bookReviews), nil
}

func (s *service) FindBookReviewByID(id string) (*application.BookReviewResponse, error) {
	bookReview, err := s.repository.FindByID(id)
	if err != nil {
		return nil, err
	}
	return application.NewBookReviewResponse(bookReview), nil
}
