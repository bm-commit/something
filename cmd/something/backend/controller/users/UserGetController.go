package users

import (
	"net/http"
	"something/internal/users/application/find"

	"github.com/gin-gonic/gin"
)

type urlParameter struct {
	ID string `uri:"id" binding:"required,uuid"`
}

// GetUserController ...
func GetUserController(finder find.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		var param urlParameter
		if err := c.ShouldBindUri(&param); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		user, err := finder.FindUserByID(param.ID)
		if err != nil {
			if err.Error() == "user not found" {
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
			"data": user,
		})
		return
	}
}
