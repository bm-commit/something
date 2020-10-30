package healthcheck

import "github.com/gin-gonic/gin"

// RegisterRoutes ...
func RegisterRoutes(router *gin.Engine) {
	router.GET("/health-check", GetController())
}
