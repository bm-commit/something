package userfollow

import (
	"net/http"
	"something/internal/userfollow/application/find"

	"github.com/gin-gonic/gin"
)

// GetFollowingController ...
func GetFollowingController(uc find.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		var param urlParameter
		if err := c.ShouldBindUri(&param); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		following, err := uc.Following(param.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Something wrong happened, try again later ...",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": following,
		})
		return
	}
}
