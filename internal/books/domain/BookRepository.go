package domain

// BookRepository ...
type BookRepository interface {
	Find() ([]*Book, error)
	FindByID(string) (*Book, error)
	Update(*Book) error
	Save(*Book) error
	Delete(string) error
}
