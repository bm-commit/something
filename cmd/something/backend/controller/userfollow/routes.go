package userfollow

import (
	m "something/cmd/something/backend/controller/middlewares"
	"something/internal/userfollow/application/find"
	"something/internal/userfollow/application/followers"
	userFind "something/internal/users/application/find"
	jwt "something/pkg/redisjwt"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes ...
func RegisterRoutes(
	finder find.Service,
	userFinder userFind.Service,
	follow followers.Service,
	tokenParams *jwt.TokenParams,
	router *gin.Engine) {
	router.GET("/users/:id/followers", GetFollowersController(finder, userFinder))
	router.GET("/users/:id/following", GetFollowingController(finder, userFinder))
	router.POST("/user/follow/:id", m.TokenAuthMiddleware(tokenParams.AccessSecret), FollowController(follow, userFinder))
	router.POST("/user/unfollow/:id", m.TokenAuthMiddleware(tokenParams.AccessSecret), UnfollowController(follow, userFinder))
}
