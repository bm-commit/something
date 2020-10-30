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
			c.JSON(http.StatusUnauthorized, "unauthorized")
			c.Abort()
			return
		}
		c.Set("user_id", au.UserID)
		c.Next()
	}
}
