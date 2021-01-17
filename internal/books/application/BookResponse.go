package application

import (
	"something/internal/books/domain"
	"time"
)

// BookResponse ...
type BookResponse struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Author      string    `json:"author"`
	Genre       string    `json:"genre"`
	Pages       int       `json:"pages"`
	Rating      float64   `json:"rating"`
	CreatedOn   time.Time `json:"created_on"`
}

// NewBookResponse ...
func NewBookResponse(book *domain.Book) *BookResponse {
	return &BookResponse{
		ID:          book.ID,
		Title:       book.Title,
		Description: book.Description,
		Author:      book.Author,
		Genre:       book.Genre,
		Pages:       book.Pages,
		Rating:      0,
		CreatedOn:   book.CreatedOn,
	}
}

// NewBooksResponse ...
func NewBooksResponse(books []*domain.Book) []*BookResponse {
	booksResponse := []*BookResponse{}
	for _, book := range books {
		booksResponse = append(booksResponse, NewBookResponse(book))
	}
	return booksResponse
}
