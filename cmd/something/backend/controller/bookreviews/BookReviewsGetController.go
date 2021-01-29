package bookreviews

import (
	"net/http"
	"something/internal/bookreviews/application"
	"something/internal/bookreviews/application/find"
	bookFind "something/internal/books/application/find"
	userFind "something/internal/users/application/find"

	"github.com/gin-gonic/gin"
)

// bookURLParameter ...
type bookURLParameter struct {
	ID string `uri:"id" binding:"required,uuid"`
}

// GetBookReviewsController ...
func GetBookReviewsController(
	finder find.Service,
	bookFinder bookFind.Service,
	userFinder userFind.Service,
) func(c *gin.Context) {
	return func(c *gin.Context) {
		var param bookURLParameter
		if err := c.ShouldBindUri(&param); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		_, err := bookFinder.FindBookByID(param.ID)
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

		bookReviews, err := finder.FindBookReviews(param.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Something wrong happened, try again later ...",
			})
			return
		}
		getUserInfoReview(bookReviews, userFinder)
		c.JSON(http.StatusOK, gin.H{
			"data": bookReviews,
		})
		return
	}
}

func getUserInfoReview(reviews []*application.BookReviewResponse, userFinder userFind.Service) []*application.BookReviewResponse {
	for _, review := range reviews {
		user, err := userFinder.FindUserByID(review.User.ID)
		if err == nil {
			review.User.Name = user.Name
			review.User.Username = user.Username
		}
	}
	return reviews
}
