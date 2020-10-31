package userfollow

import (
	"net/http"
	"something/internal/userfollow/application/followers"

	"github.com/gin-gonic/gin"
)

// UnfollowController ...
func UnfollowController(uc followers.Service) func(c *gin.Context) {
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
		if param.ID == userID.(string) {
			c.Status(http.StatusBadRequest)
			return
		}
		err := uc.Unfollow(userID.(string), param.ID)
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
		c.Status(http.StatusOK)
		return
	}
}
