package userfollow

import (
	"net/http"
	"something/internal/userfollow/application/find"
	userFind "something/internal/users/application/find"

	"github.com/gin-gonic/gin"
)

// GetFollowersController ...
func GetFollowersController(uc find.Service, userFinder userFind.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		var param urlParameter
		if err := c.ShouldBindUri(&param); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		_, err := userFinder.FindUserByID(param.ID)
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

		followers, err := uc.Followers(param.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Something wrong happened, try again later ...",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": followers,
		})
		return
	}
}
