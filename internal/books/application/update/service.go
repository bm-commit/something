package update

import (
	"encoding/json"
	"errors"
	"something/internal/books/application"
	"something/internal/books/domain"
	"strings"
)

// Service ...
type Service interface {
	UpdateBookByID(*application.BookCommand) error
}

type service struct {
	repository domain.BookRepository
}

// NewService ...
func NewService(repository domain.BookRepository) Service {
	return &service{repository: repository}
}

func (s *service) UpdateBookByID(book *application.BookCommand) error {
	existingBook, _ := s.repository.FindByID(book.ID)
	if existingBook == nil {
		return errors.New("book not found")
	}

	out, err := json.Marshal(book)
	if err != nil {
		return err
	}
	// Merge new book fields with existing book fields
	err = json.NewDecoder(strings.NewReader(string(out))).Decode(existingBook)
	if err != nil {
		return err
	}

	updatedBook, err := domain.NewBook(
		existingBook.ID, existingBook.Title, existingBook.Description,
		existingBook.Author, existingBook.Genre, existingBook.Pages)
	if err != nil {
		return err
	}
	updatedBook.CreatedOn = existingBook.CreatedOn

	err = s.repository.Update(updatedBook)
	return err
}
