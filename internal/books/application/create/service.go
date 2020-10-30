package create

import (
	"errors"
	"something/internal/books/application"
	"something/internal/books/domain"
)

// Service ...
type Service interface {
	CreateBook(*application.BookCommand) error
}

type service struct {
	repository domain.BookRepository
}

// NewService ...
func NewService(repository domain.BookRepository) Service {
	return &service{repository: repository}
}

func (s *service) CreateBook(command *application.BookCommand) error {
	existingBookID, _ := s.repository.FindByID(command.ID)
	if existingBookID != nil {
		return errors.New("book id already exists")
	}
	book, err := domain.NewBook(command.ID, command.Title, command.Description,
		command.Author, command.Genre, command.Pages)
	if err != nil {
		return err
	}
	err = s.repository.Save(book)
	if err != nil {
		return err
	}
	return nil
}
