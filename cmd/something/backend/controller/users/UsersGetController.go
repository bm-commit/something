package users

import (
	"net/http"
	"strings"

	bookFinder "something/internal/books/application/find"
	"something/internal/users/application/find"

	"github.com/gin-gonic/gin"
)

// GetUsersController ...
func GetUsersController(finder find.Service, bFinder bookFinder.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		username := c.Query("username")
		if username != "" {
			user, err := finder.FindUserByUsername(strings.ToLower(username))
			if err != nil {
				if err.Error() == "username not found" {
					c.JSON(http.StatusNotFound, gin.H{
						"error": err.Error(),
					})
					return
				}
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Something wrong happened, try again later ...",
				})
				return
			}
			interests := classifyBookInterests(user.Interests, bFinder)
			c.JSON(http.StatusOK, gin.H{
				"data":      user,
				"interests": interests,
			})
			return
		}

		users, err := finder.FindUsers()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Something wrong happened, try again later ...",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"data": users,
		})
		return
	}
}

// TODO Refactor code below to classify user book interests

type bookShort struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

func classifyBookInterests(interests map[string]string, finder bookFinder.Service) map[string][]*bookShort {

	bookInterests := map[string][]*bookShort{}

	pending := []*bookShort{}
	reading := []*bookShort{}
	done := []*bookShort{}

	for bookID, status := range interests {
		book := &bookShort{}
		bookResponse, err := finder.FindBookByID(bookID)
		if err != nil {
			continue
		}
		book.ID = bookResponse.ID
		book.Title = bookResponse.Title
		book.Author = bookResponse.Author

		switch status {
		case "pending":
			pending = append(pending, book)
			break
		case "reading":
			reading = append(reading, book)
			break
		case "done":
			done = append(done, book)
			break
		}
	}

	bookInterests["pending"] = pending
	bookInterests["reading"] = reading
	bookInterests["done"] = done

	return bookInterests
}
