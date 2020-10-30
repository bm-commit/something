package bookreviews

import (
	"net/http"

	"something/internal/bookreviews/application/create"

	"github.com/gin-gonic/gin"
)

type urlParameters struct {
	BookID   string `uri:"id" binding:"required,uuid"`
	ReviewID string `uri:"review_id" binding:"required,uuid"`
}

// PutController ...
func PutController(creator create.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		var param urlParameters
		if err := c.ShouldBindUri(&param); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var request create.BookReviewCommand
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := request.Validate(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userID, ok := c.Get("user_id")
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Something wrong happened, try again later ...",
			})
			return
		}

		request.ID = param.ReviewID
		request.BookID = param.BookID
		request.UserID = userID.(string)

		err := creator.CreateBookReview(&request)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Something wrong happened, try again later ...",
			})
			return
		}
		c.Status(http.StatusCreated)
		return
	}
}
