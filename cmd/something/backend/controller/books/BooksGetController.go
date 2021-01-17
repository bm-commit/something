package books

import (
	"net/http"
	bookReview "something/internal/bookreviews/application"
	bookReviewFinder "something/internal/bookreviews/application/find"
	"something/internal/books/application"
	"something/internal/books/application/find"
	"something/internal/helpers"

	"github.com/gin-gonic/gin"
)

// GetBooksController ...
func GetBooksController(finder find.Service, reviewFinder bookReviewFinder.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		books, err := finder.FindBooks()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Something wrong happened, try again later ...",
			})
			return
		}
		addRatingToBooks(books, reviewFinder)
		c.JSON(http.StatusOK, gin.H{
			"data": books,
		})
		return
	}
}

// TODO Refactor
func addRatingToBooks(books []*application.BookResponse, reviewFinder bookReviewFinder.Service) {
	for _, book := range books {
		bookReviews, err := reviewFinder.FindBookReviews(book.ID)
		if err == nil {
			book.Rating = getBookRating(bookReviews)
		}
	}
}

func getBookRating(bookReviews []*bookReview.BookReviewResponse) float64 {
	var sumRating float64 = 0
	if len(bookReviews) == 0 {
		return sumRating
	}
	for _, review := range bookReviews {
		sumRating += review.Rating
	}
	return helpers.Round(sumRating/float64(len(bookReviews)), 0.5)
}
