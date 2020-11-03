package bookreviews

import (
	m "something/cmd/something/backend/controller/middlewares"
	"something/internal/bookreviews/application/create"
	"something/internal/bookreviews/application/delete"
	"something/internal/bookreviews/application/find"
	"something/internal/bookreviews/application/update"
	bookFind "something/internal/books/application/find"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes ...
func RegisterRoutes(
	finder find.Service,
	bookFinder bookFind.Service,
	creator create.Service,
	updater update.Service,
	delete delete.Service,
	accessSecret string, router *gin.Engine) {
	router.GET("/books/:id/reviews", GetBookReviewsController(finder, bookFinder))
	router.GET("/book/reviews/:review_id", GetBookReviewController(finder))
	router.PATCH("/book/reviews/:review_id", m.TokenAuthMiddleware(accessSecret), PatchController(updater))
	router.PUT("/books/:id/reviews/:review_id", m.TokenAuthMiddleware(accessSecret), PutController(creator))
	router.DELETE("/book/reviews/:review_id", m.TokenAuthStaffMiddleware(accessSecret), DeleteBookReviewController(delete))
}
