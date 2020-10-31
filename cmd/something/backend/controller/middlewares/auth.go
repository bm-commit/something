package middlewares

import (
	"net/http"
	jwt "something/pkg/redisjwt"

	"github.com/gin-gonic/gin"
)

// TokenAuthMiddleware ...
func TokenAuthMiddleware(accessSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		au, err := jwt.ExtractTokenMetadata(c.Request, accessSecret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			c.Abort()
			return
		}
		c.Set("user_id", au.UserID)
		c.Next()
	}
}

const authorizedRole = "staff"

// TokenAuthStaffMiddleware ...
func TokenAuthStaffMiddleware(accessSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		au, err := jwt.ExtractTokenMetadata(c.Request, accessSecret)
		if err != nil || au.Role != authorizedRole {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			c.Abort()
			return
		}
		c.Set("user_id", au.UserID)
		c.Next()
	}
}
