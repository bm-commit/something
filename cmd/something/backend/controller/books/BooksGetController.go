package books

import (
	"net/http"
	bookReviewFinder "something/internal/bookreviews/application/find"
	"something/internal/books/application"
	"something/internal/books/application/find"
	"something/internal/helpers"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetBooksController ...
func GetBooksController(finder find.Service, reviewFinder bookReviewFinder.Service) func(c *gin.Context) {
	return func(c *gin.Context) {

		criteria := getQueryParameters(c)

		books, err := finder.FindBooks(criteria)
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

func getQueryParameters(c *gin.Context) *find.Criteria {
	page, _ := strconv.Atoi(c.Query("page"))
	perPage, _ := strconv.Atoi(c.Query("per_page"))

	return &find.Criteria{
		Page:    page,
		PerPage: perPage,
		Query:   c.Query("q"),
		Genre:   c.Query("genre"),
		Author:  c.Query("author"),
	}
}

// TODO Refactor
func addRatingToBooks(books []*application.BookResponse, reviewFinder bookReviewFinder.Service) {
	for _, book := range books {
		bookReviews, err := reviewFinder.FindBookReviews(book.ID)
		if err == nil {
			book.Rating = helpers.GetBookRating(bookReviews)
		}
	}
}
