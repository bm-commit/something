package persistence

import (
	"errors"
	"something/internal/books/domain"
)

type repository struct {
	books map[string]*domain.Book
}

var (
	bookInstance *repository
)

// NewInMemoryBookRepository ...
func NewInMemoryBookRepository() domain.BookRepository {
	bookInstance = &repository{
		books: make(map[string]*domain.Book),
	}
	return bookInstance
}

func (r *repository) Find() ([]*domain.Book, error) {
	var books []*domain.Book
	for _, book := range r.books {
		books = append(books, book)
	}
	return books, nil
}

func (r *repository) FindByID(id string) (*domain.Book, error) {
	book, ok := r.books[id]
	if !ok {
		return nil, errors.New("book not found")
	}
	return book, nil
}

func (r *repository) Update(book *domain.Book) error {
	r.books[book.ID] = book
	return nil
}

func (r *repository) Save(book *domain.Book) error {
	r.books[book.ID] = book
	return nil
}

func (r *repository) Delete(id string) error {
	_, ok := r.books[id]
	if ok {
		delete(r.books, id)
	}
	return nil
}
