package domain

// BookReviewCriteria ...
type BookReviewCriteria struct {
	Sort int64
}

// NewBookReviewCriteria ...
func NewBookReviewCriteria(sort int) *BookReviewCriteria {
	return &BookReviewCriteria{
		Sort: int64(sort),
	}
}
