package bookreviews

import (
	"net/http"
	"something/internal/bookreviews/application/delete"

	"github.com/gin-gonic/gin"
)

// DeleteBookReviewController ...
func DeleteBookReviewController(delete delete.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		var param urlParameter
		if err := c.ShouldBindUri(&param); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := delete.DeleteBookReviewByID(param.ID)
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
		c.Status(http.StatusNoContent)
		return
	}
}
