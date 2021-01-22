package helpers

import (
	"math"

	"something/internal/bookreviews/application"
)

// Round value to specific unit
func Round(x, unit float64) float64 {
	return math.Round(x/unit) * unit
}

func GetBookRating(reviews []*application.BookReviewResponse) float64 {
	var sumRating float64 = 0
	if len(reviews) == 0 {
		return sumRating
	}
	for _, review := range reviews {
		sumRating += review.Rating
	}
	return Round(sumRating/float64(len(reviews)), 0.5)
}
