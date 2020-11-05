package bookreviews

import (
	"net/http"
	"something/internal/bookreviews/application/find"

	"github.com/gin-gonic/gin"
)

// urlParameter ...
type urlParameter struct {
	ID string `uri:"review_id" binding:"required,uuid"`
}

// GetBookReviewController ...
func GetBookReviewController(finder find.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		var param urlParameter
		if err := c.ShouldBindUri(&param); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		bookReview, err := finder.FindBookReviewByID(param.ID)
		if err != nil {
			if err.Error() == "book review not found" {
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

		c.JSON(http.StatusOK, gin.H{
			"data": bookReview,
		})
		return
	}
}
