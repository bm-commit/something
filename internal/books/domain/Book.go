package domain

import "time"

// Book ...
type Book struct {
	ID          string
	Title       string
	Description string
	Author      string
	Genre       string
	Pages       int
	CreatedOn   time.Time
}

// NewBook ...
func NewBook(id, title, description, author, genre string, pages int) (*Book, error) {
	return &Book{
		ID:          id,
		Title:       title,
		Description: description,
		Author:      author,
		Genre:       genre,
		Pages:       pages,
		CreatedOn:   time.Now().UTC(),
	}, nil
}
