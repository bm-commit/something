package books

import (
	"something/internal/books/application/create"
	"something/internal/books/application/delete"
	"something/internal/books/application/find"
	"something/internal/books/application/update"

	m "something/cmd/something/backend/controller/middlewares"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes ...
func RegisterRoutes(
	finder find.Service,
	creator create.Service,
	update update.Service,
	deletor delete.Service,
	accessSecret string,
	router *gin.Engine) {
	booksRouter := router.Group("/books")
	{
		booksRouter.GET("/", GetBooksController(finder))
		booksRouter.GET("/:id", GetBookController(finder))
		booksRouter.PUT("/:id", m.TokenAuthStaffMiddleware(accessSecret), PutController(creator))
		booksRouter.PATCH("/:id", m.TokenAuthStaffMiddleware(accessSecret), PatchController(update))
		booksRouter.DELETE("/:id", m.TokenAuthStaffMiddleware(accessSecret), DeleteBookController(deletor))
	}
}
