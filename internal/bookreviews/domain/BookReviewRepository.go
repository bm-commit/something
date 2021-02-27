package domain

// BookReviewRepository ...
type BookReviewRepository interface {
	Find(string) ([]*BookReview, error)
	FindByID(string) (*BookReview, error)
	FindReviews(*BookReviewCriteria) ([]*BookReviewShort, error)
	Update(*BookReview) error
	Save(*BookReview) error
	Delete(string) error
}
