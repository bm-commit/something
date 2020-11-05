package domain

// BookReviewRepository ...
type BookReviewRepository interface {
	Find(string) ([]*BookReview, error)
	FindByID(string) (*BookReview, error)
	Update(*BookReview) error
	Save(*BookReview) error
	Delete(string) error
}
