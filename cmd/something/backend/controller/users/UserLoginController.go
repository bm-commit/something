package users

import (
	"net/http"
	"something/internal/users/application/login"
	jwt "something/pkg/redisjwt"

	"github.com/gin-gonic/gin"
)

// LoginController ...
func LoginController(usecase login.Service, tokenParams *jwt.TokenParams) func(c *gin.Context) {
	return func(c *gin.Context) {

		var request login.Command
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := request.Validate(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, err := usecase.Login(&request)
		if err != nil {
			if err.Error() == "email not found" {
				c.JSON(http.StatusNotFound, gin.H{
					"error": err.Error(),
				})
				return
			} else if err.Error() == "invalid email or password" {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": err.Error(),
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Something wrong happened, try again later ...",
			})
			return
		}
		ts, err := jwt.CreateToken(user.ID, tokenParams)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, err.Error())
			return
		}
		tokens := map[string]string{
			"access_token":  ts.AccessToken,
			"refresh_token": ts.RefreshToken,
		}

		c.JSON(http.StatusOK, gin.H{
			"data":   user,
			"tokens": tokens,
		})
		return
	}
}
