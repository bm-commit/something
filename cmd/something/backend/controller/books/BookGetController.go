package books

import (
	"net/http"
	bookReviewFinder "something/internal/bookreviews/application/find"
	"something/internal/books/application/find"
	"something/internal/helpers"

	"github.com/gin-gonic/gin"
)

type urlParameter struct {
	ID string `uri:"id" binding:"required,uuid"`
}

// GetBookController ...
func GetBookController(finder find.Service, reviewFinder bookReviewFinder.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		var param urlParameter
		if err := c.ShouldBindUri(&param); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		book, err := finder.FindBookByID(param.ID)
		if err != nil {
			if err.Error() == "book not found" {
				c.JSON(http.StatusNotFound, gin.H{
					"error": err.Error(),
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Something wrong happened, try again later ...",
			})
			return
		}
		bookReviews, err := reviewFinder.FindBookReviews(book.ID)
		if err == nil {
			book.Rating = helpers.GetBookRating(bookReviews)
		}
		c.JSON(http.StatusOK, gin.H{
			"data": book,
		})
		return
	}
}
