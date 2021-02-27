package domain

// BookCriteria ...
type BookCriteria struct {
	Page    int64
	PerPage int64
	Query   string
	Genre   string
	Author  string
}

// NewBookCriteria ...
func NewBookCriteria(page, perPage int, query, genre, author string) *BookCriteria {
	return &BookCriteria{
		Page:    int64(page),
		PerPage: int64(perPage),
		Query:   query,
		Genre:   genre,
		Author:  author,
	}
}
