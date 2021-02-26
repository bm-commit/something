package find

import (
	"something/internal/bookreviews/application"
	"something/internal/bookreviews/domain"
)

// Service ...
type Service interface {
	FindBookReviews(bookID string) ([]*application.BookReviewResponse, error)
	FindBookReviewByID(id string) (*application.BookReviewResponse, error)
	FindReviews(criteria *Criteria) ([]*application.BookRatingResponse, error)
}

type service struct {
	repository domain.BookReviewRepository
}

// NewService ...
func NewService(repository domain.BookReviewRepository) Service {
	return &service{repository: repository}
}

func (s *service) FindBookReviews(bookID string) ([]*application.BookReviewResponse, error) {
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

func (s *service) FindReviews(criteria *Criteria) ([]*application.BookRatingResponse, error) {
	newBookReviewCriteria := domain.NewBookReviewCriteria(criteria.Sort)
	bookReviews, err := s.repository.FindReviews(newBookReviewCriteria)
	if err != nil {
		return nil, err
	}
	return application.NewReviewShortResponse(bookReviews), nil
}
