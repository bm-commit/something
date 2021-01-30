package domain

// BookRepository ...
type BookRepository interface {
	Find(*BookCriteria) ([]*Book, error)
	FindByID(string) (*Book, error)
	Update(*Book) error
	Save(*Book) error
	Delete(string) error
}
