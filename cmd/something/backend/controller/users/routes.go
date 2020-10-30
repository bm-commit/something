package users

import (
	m "something/cmd/something/backend/controller/middlewares"
	"something/internal/users/application/create"
	"something/internal/users/application/delete"
	"something/internal/users/application/find"
	"something/internal/users/application/login"
	"something/internal/users/application/update"
	jwt "something/pkg/redisjwt"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes ...
func RegisterRoutes(
	finder find.Service,
	creator create.Service,
	updater update.Service,
	deleter delete.Service,
	login login.Service,
	tokenParams *jwt.TokenParams,
	router *gin.Engine) {
	usersRouter := router.Group("/users")
	{
		usersRouter.GET("/", GetUsersController(finder))
		usersRouter.GET("/:id", GetUserController(finder))
		usersRouter.PUT("/:id", RegisterController(creator))
		usersRouter.PATCH("/:id", m.TokenAuthMiddleware(tokenParams.AccessSecret), PatchController(updater))
		usersRouter.DELETE("/:id", m.TokenAuthMiddleware(tokenParams.AccessSecret), DeleteUserController(deleter))
	}
	router.POST("/login", LoginController(login, tokenParams))
}
