package userfollow

import (
	m "something/cmd/something/backend/controller/middlewares"
	"something/internal/userfollow/application/find"
	"something/internal/userfollow/application/followers"
	jwt "something/pkg/redisjwt"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes ...
func RegisterRoutes(
	finder find.Service,
	follow followers.Service,
	tokenParams *jwt.TokenParams,
	router *gin.Engine) {
	router.GET("/users/:id/followers", GetFollowersController(finder))
	router.GET("/users/:id/following", GetFollowingController(finder))
	router.POST("/user/follow/:id", m.TokenAuthMiddleware(tokenParams.AccessSecret), FollowController(follow))
	router.POST("/user/unfollow/:id", m.TokenAuthMiddleware(tokenParams.AccessSecret), UnfollowController(follow))
}
