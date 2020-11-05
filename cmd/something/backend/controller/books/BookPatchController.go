package books

import (
	"net/http"
	"something/internal/books/application"
	"something/internal/books/application/update"

	"github.com/gin-gonic/gin"
)

// PatchController ...
func PatchController(update update.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		var param urlParameter
		if err := c.ShouldBindUri(&param); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var request application.BookCommand
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		request.ID = param.ID

		err := update.UpdateBookByID(&request)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Something wrong happened, try again later ...",
			})
			return
		}
		c.Status(http.StatusOK)
		return
	}
}
