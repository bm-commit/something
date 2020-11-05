package users

import (
	"net/http"
	bookFind "something/internal/books/application/find"
	"something/internal/users/application/update"

	"github.com/gin-gonic/gin"
)

type urlBookParameter struct {
	ID string `uri:"book_id" binding:"required,uuid"`
}

// InterestsPatchController ...
func InterestsPatchController(us update.Service, bookFinder bookFind.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		var param urlBookParameter
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

		userID, ok := c.Get("user_id")
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Something wrong happened, try again later ...",
			})
			return
		}

		var request update.UserInterestsCommand
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := request.Validate(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		request.UserID = userID.(string)
		request.BookID = param.ID

		err = us.UpdateUserInterests(&request)
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
