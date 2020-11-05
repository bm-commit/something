package users

import (
	"net/http"
	"something/internal/users/application/create"

	"github.com/gin-gonic/gin"
)

// RegisterController ...
func RegisterController(creator create.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		var param urlParameter
		if err := c.ShouldBindUri(&param); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var request create.UserCommand
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := request.Validate(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		request.ID = param.ID

		err := creator.CreateUser(&request)
		if err != nil {
			if err.Error() == "email already in use" ||
				err.Error() == "username already in use" ||
				err.Error() == "user id already exists" {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Something wrong happened, try again later ...",
			})
			return
		}
		c.Status(http.StatusCreated)
		return
	}
}
