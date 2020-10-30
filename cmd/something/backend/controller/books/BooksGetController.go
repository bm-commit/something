package books

import (
	"net/http"
	"something/internal/books/application/find"

	"github.com/gin-gonic/gin"
)

// GetBooksController ...
func GetBooksController(finder find.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		books, err := finder.FindBooks()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Something wrong happened, try again later ...",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"data": books,
		})
		return
	}
}
