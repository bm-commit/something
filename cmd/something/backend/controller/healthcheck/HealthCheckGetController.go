package healthcheck

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetController ...
func GetController() func(c *gin.Context) {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	}
}
