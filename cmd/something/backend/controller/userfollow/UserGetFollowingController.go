package userfollow

import (
	"net/http"
	"something/internal/userfollow/application"
	"something/internal/userfollow/application/find"
	userFind "something/internal/users/application/find"

	"github.com/gin-gonic/gin"
)

// GetFollowingController ...
func GetFollowingController(uc find.Service, userFinder userFind.Service) func(c *gin.Context) {
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
		following, err := uc.Following(param.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Something wrong happened, try again later ...",
			})
			return
		}

		followingLong := getFollowingLong(following, userFinder)

		c.JSON(http.StatusOK, gin.H{
			"data": followingLong,
		})
		return
	}
}

// TODO Refactor
func getFollowingLong(following []*application.UserFollowResponse, userFinder userFind.Service) []*application.UserFollowResponseLong {
	followingLong := []*application.UserFollowResponseLong{}
	for _, following := range following {
		user, err := userFinder.FindUserByID(following.To)
		if err == nil {
			f := &application.UserFollowResponseLong{
				ID:       user.ID,
				Name:     user.Name,
				Username: user.Username,
				FollowAt: following.CreatedOn,
			}
			followingLong = append(followingLong, f)
		}
	}
	return followingLong
}
