package bookreviews

import (
	"net/http"
	"something/internal/bookreviews/application/update"

	"github.com/gin-gonic/gin"
)

// PatchController ...
func PatchController(us update.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		var param urlParameter
		if err := c.ShouldBindUri(&param); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userID, ok := c.Get("user_id")
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Something wrong happened, try again later ...",
			})
			return
		}

		var request update.BookReviewCommand
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		request.ID = param.ID
		request.UserID = userID.(string)

		err := us.UpdateBookReviewByID(&request)
		if err != nil {
			if err.Error() == "unauthorized" {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": err.Error(),
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Something wrong happened, try again later ...",
			})
			return
		}
		c.Status(http.StatusOK)
		return
	}
}
