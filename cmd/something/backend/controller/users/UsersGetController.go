package users

import (
	"net/http"
	"something/internal/users/application/find"

	"github.com/gin-gonic/gin"
)

// GetUsersController ...
func GetUsersController(finder find.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		username := c.Query("username")
		if username != "" {
			user, err := finder.FindUserByUsername(username)
			if err != nil {
				if err.Error() == "username not found" {
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

		users, err := finder.FindUsers()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Something wrong happened, try again later ...",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"data": users,
		})
		return
	}
}
