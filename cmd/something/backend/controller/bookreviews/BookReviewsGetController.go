package bookreviews

import (
	"net/http"
	"something/internal/bookreviews/application/find"

	"github.com/gin-gonic/gin"
)

// bookURLParameter ...
type bookURLParameter struct {
	ID string `uri:"id" binding:"required,uuid"`
}

// GetBookReviewsController ...
func GetBookReviewsController(finder find.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		var param bookURLParameter
		if err := c.ShouldBindUri(&param); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		bookReviews, err := finder.FindBookReviews(param.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Something wrong happened, try again later ...",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"data": bookReviews,
		})
		return
	}
}
