package persistence

import (
	"errors"

	"something/internal/bookreviews/domain"
)

type repository struct {
	bookReviews map[string]*domain.BookReview
}

var (
	reviewInstance *repository
)

// NewInMemoryBookReviewsRepository ...
func NewInMemoryBookReviewsRepository() domain.BookReviewRepository {
	reviewInstance = &repository{
		bookReviews: make(map[string]*domain.BookReview),
	}
	return reviewInstance
}

func (r *repository) Find(bookID string) ([]*domain.BookReview, error) {
	var bookReviews []*domain.BookReview
	for _, bookReview := range r.bookReviews {
		if bookReview.BookID == bookID {
			bookReviews = append(bookReviews, bookReview)
		}
	}

	return bookReviews, nil
}

func (r *repository) FindByID(id string) (*domain.BookReview, error) {
	bookReview, ok := r.bookReviews[id]
	if !ok {
		return nil, errors.New("book review not found")
	}

	return bookReview, nil
}

func (r *repository) Update(bookReview *domain.BookReview) error {
	r.bookReviews[bookReview.ID] = bookReview
	return nil
}

func (r *repository) Save(bookReview *domain.BookReview) error {
	r.bookReviews[bookReview.ID] = bookReview
	return nil
}

func (r *repository) Delete(id string) error {
	_, ok := r.bookReviews[id]
	if ok {
		delete(r.bookReviews, id)
	}
	return nil
}
