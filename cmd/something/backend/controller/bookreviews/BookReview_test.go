package bookreviews

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"something/internal/bookreviews/application/create"
	"something/internal/bookreviews/application/delete"
	"something/internal/bookreviews/application/find"
	"something/internal/bookreviews/application/update"
	"something/internal/bookreviews/domain"
	"something/internal/bookreviews/infraestructure/persistence"
	bookDomain "something/internal/books/domain"
	bookPersistance "something/internal/books/infraestructure/persistence"
	jwt "something/pkg/redisjwt"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var tokenParams *jwt.TokenParams = &jwt.TokenParams{
	AccessSecret:  "secure-access-token",
	RefreshSecret: "secure-refresh-token",
	AccessTime:    time.Minute * 1,
	RefreshTime:   time.Minute * 1,
}

const bookID = "c9d6e6f0-27d9-47d2-851e-bb42f72565ed"
const userID = "c015f5ce-3b42-44c8-8b82-f011b23b989a"

func setupServer(
	bookReviewRepo domain.BookReviewRepository,
	bookRepo bookDomain.BookRepository) *gin.Engine {
	router := gin.Default()
	finder := find.NewService(bookReviewRepo, bookRepo)
	updater := update.NewService(bookReviewRepo)
	creator := create.NewService(bookReviewRepo)
	deletor := delete.NewService(bookReviewRepo)
	RegisterRoutes(finder, creator, updater, deletor, tokenParams.AccessSecret, router)
	return router
}

func TestBookReviewCheck(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Book Review Suite")
}

var _ = Describe("Server", func() {
	var server *httptest.Server
	var bookRepo bookDomain.BookRepository
	var bookReviewRepo domain.BookReviewRepository

	BeforeEach(func() {
		bookRepo = bookPersistance.NewInMemoryBookRepository()
		defaultBook, _ := bookDomain.NewBook(bookID, "title", "description", "author", "genre", 1)
		bookRepo.Save(defaultBook)
		bookReviewRepo = persistence.NewInMemoryBookReviewsRepository()
		server = httptest.NewServer(setupServer(bookReviewRepo, bookRepo))
	})

	AfterEach(func() {
		server.Close()
	})

	Context("When GET request is sent to /books/:id/reviews", func() {
		It("Returns null data if not exists reviews", func() {
			resp, err := http.Get(server.URL + "/books/" + bookID + "/reviews")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusOK))

			body, err := ioutil.ReadAll(resp.Body)
			defer resp.Body.Close()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(string(body)).To(MatchJSON(`{"data":null}`))
		})
		It("Returns an existing books review", func() {
			newBookReview, _ := domain.NewBookReview("1", "abc", 1, bookID, userID)
			bookReviewRepo.Save(newBookReview)

			resp, err := http.Get(server.URL + "/books/" + bookID + "/reviews")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusOK))

			body, err := ioutil.ReadAll(resp.Body)
			defer resp.Body.Close()
			Expect(err).ShouldNot(HaveOccurred())

			Expect(string(body)).To(MatchJSON(`
			{
				"data":
					[
						{
							"id":"` + newBookReview.ID + `",
							"text":"` + newBookReview.Text + `",
							"rating":` + strconv.Itoa(newBookReview.Rating) + ` ,
							"book_id":"` + newBookReview.BookID + `",
							"user_id":"` + newBookReview.UserID + `",
							"created_on":"` + newBookReview.CreatedOn.Format("2006-01-02T15:04:05.999999999Z07:00") + `"
						}
					]
			}`))
		})
	})
	Context("When GET request by ID is sent to /book/reviews/:review_id", func() {
		It("Returns an existing book review by id", func() {
			newBookReview, _ := domain.NewBookReview("c0b369a0-8de4-417d-a905-c33644c2907d", "abc", 1, bookID, userID)
			bookReviewRepo.Save(newBookReview)

			resp, err := http.Get(
				server.URL + "/book/reviews/c0b369a0-8de4-417d-a905-c33644c2907d")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusOK))

			body, err := ioutil.ReadAll(resp.Body)
			defer resp.Body.Close()
			Expect(err).ShouldNot(HaveOccurred())

			Expect(string(body)).To(MatchJSON(`
			{
				"data":
					{
						"id":"` + newBookReview.ID + `",
						"text":"` + newBookReview.Text + `",
						"rating":` + strconv.Itoa(newBookReview.Rating) + `,
						"book_id":"` + newBookReview.BookID + `",
						"user_id":"` + newBookReview.UserID + `",
						"created_on":"` + newBookReview.CreatedOn.Format("2006-01-02T15:04:05.999999999Z07:00") + `"
					}
			}`))
		})
		It("Returns an 404 status code in non existing id", func() {
			resp, err := http.Get(
				server.URL + "/book/reviews/c0b369a0-8de4-417d-a905-c33644c2907d")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusNotFound))

			body, err := ioutil.ReadAll(resp.Body)
			defer resp.Body.Close()
			Expect(err).ShouldNot(HaveOccurred())

			Expect(string(body)).To(MatchJSON(`{"error":"book review not found"}`))
		})
	})
	Context("When PUT request by ID is sent to /books/:id/reviews/:review_id", func() {
		It("Create a new books review", func() {
			reviewID := "c0b369a0-8de4-417d-a905-c33644c2907d"
			bookReview := map[string]interface{}{
				"text":   "abc",
				"rating": 1,
			}
			jsonReq, err := json.Marshal(bookReview)

			generateAuth, err := jwt.CreateToken(userID, "default", tokenParams)
			Expect(err).ShouldNot(HaveOccurred())

			req, err := http.NewRequest(
				http.MethodPut,
				server.URL+"/books/"+bookID+"/reviews/"+reviewID,
				bytes.NewBuffer(jsonReq))
			req.Header.Set("Content-Type", "application/json; charset=utf-8")
			req.Header.Set("Authorization", "Bearer "+generateAuth.AccessToken)
			client := &http.Client{}

			resp, err := client.Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusCreated))

			createdReview, _ := bookReviewRepo.FindByID(reviewID)
			Expect(createdReview).ShouldNot(BeNil())
		})
		It("Returns an 400 status code with an invalid uuid", func() {
			generateAuth, err := jwt.CreateToken(userID, "default", tokenParams)
			Expect(err).ShouldNot(HaveOccurred())

			req, err := http.NewRequest(
				http.MethodPut,
				server.URL+"/books/"+bookID+"/reviews/1",
				nil)
			req.Header.Set("Authorization", "Bearer "+generateAuth.AccessToken)
			client := &http.Client{}
			resp, err := client.Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusBadRequest))
		})
	})
	Context("When PATCH request by ID is sent to /book/reviews/:review_id", func() {
		It("modify an existing review", func() {
			newBookReview, _ := domain.NewBookReview("47bb4bed-e1ee-413a-85ed-2cc4c598e562", "abc", 1, bookID, userID)
			bookReviewRepo.Save(newBookReview)

			fieldsToModify := map[string]interface{}{
				"text": "lorem ipsum",
			}
			jsonReq, err := json.Marshal(fieldsToModify)

			generateAuth, err := jwt.CreateToken(userID, "default", tokenParams)
			Expect(err).ShouldNot(HaveOccurred())

			req, err := http.NewRequest(
				http.MethodPatch,
				server.URL+"/book/reviews/"+newBookReview.ID,
				bytes.NewBuffer(jsonReq),
			)
			req.Header.Set("Content-Type", "application/json; charset=utf-8")
			req.Header.Set("Authorization", "Bearer "+generateAuth.AccessToken)
			client := &http.Client{}

			resp, err := client.Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusOK))

			updatedBookReview, _ := domain.NewBookReview(
				newBookReview.ID,
				fieldsToModify["text"].(string),
				newBookReview.Rating,
				bookID,
				userID,
			)
			updatedBookReview.CreatedOn = newBookReview.CreatedOn

			bookReview, _ := bookReviewRepo.FindByID(newBookReview.ID)
			Expect(bookReview).Should(BeEquivalentTo(updatedBookReview))
		})
		It("Returns an 401 status code with not review owner", func() {
			newBookReview, _ := domain.NewBookReview("47bb4bed-e1ee-413a-85ed-2cc4c598e562", "abc", 1, bookID, userID)
			bookReviewRepo.Save(newBookReview)

			fieldsToModify := map[string]interface{}{
				"text": "lorem ipsum yep",
			}
			jsonReq, err := json.Marshal(fieldsToModify)

			generateAuth, err := jwt.CreateToken("55a5cd53-6d6d-46f1-9eb0-689435c269f0", "default", tokenParams)
			Expect(err).ShouldNot(HaveOccurred())

			req, err := http.NewRequest(
				http.MethodPatch,
				server.URL+"/book/reviews/"+newBookReview.ID,
				bytes.NewBuffer(jsonReq),
			)
			req.Header.Set("Content-Type", "application/json; charset=utf-8")
			req.Header.Set("Authorization", "Bearer "+generateAuth.AccessToken)
			client := &http.Client{}

			resp, err := client.Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusUnauthorized))

			body, err := ioutil.ReadAll(resp.Body)
			defer resp.Body.Close()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(string(body)).To(MatchJSON(`{"error":"unauthorized"}`))
		})
	})
	Context("When DELETE request by ID is sent to /book/reviews/:review_id", func() {
		It("delete an existing book review", func() {
			newBookReview, _ := domain.NewBookReview("f73cbfc4-1971-49d6-8964-d696b4e2e220", "abc", 1, bookID, userID)
			bookReviewRepo.Save(newBookReview)

			generateAuth, err := jwt.CreateToken(userID, "staff", tokenParams)
			Expect(err).ShouldNot(HaveOccurred())
			req, err := http.NewRequest(
				http.MethodDelete,
				server.URL+"/book/reviews/f73cbfc4-1971-49d6-8964-d696b4e2e220",
				nil)
			req.Header.Set("Authorization", "Bearer "+generateAuth.AccessToken)
			client := &http.Client{}
			resp, err := client.Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusNoContent))

			bookReview, _ := bookReviewRepo.FindByID(newBookReview.ID)
			Expect(bookReview).Should(BeNil())
		})
		It("return an 404 status code in non existing bookReview", func() {
			generateAuth, err := jwt.CreateToken(userID, "staff", tokenParams)
			Expect(err).ShouldNot(HaveOccurred())

			req, err := http.NewRequest(
				http.MethodDelete,
				server.URL+"/book/reviews/427bfa5b-9144-4f1c-8069-b42307192d65",
				nil,
			)
			req.Header.Set("Authorization", "Bearer "+generateAuth.AccessToken)
			client := &http.Client{}
			resp, err := client.Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusNotFound))
		})
		It("return an 401 status code in not admin user", func() {
			bookReviewID := "f73cbfc4-1971-49d6-8964-d696b4e2e220"
			newBookReview, _ := domain.NewBookReview(bookReviewID, "abc", 1, bookID, userID)
			bookReviewRepo.Save(newBookReview)

			generateAuth, err := jwt.CreateToken(userID, "default", tokenParams)
			Expect(err).ShouldNot(HaveOccurred())

			req, err := http.NewRequest(
				http.MethodDelete,
				server.URL+"/book/reviews/"+bookReviewID,
				nil,
			)
			req.Header.Set("Authorization", "Bearer "+generateAuth.AccessToken)
			client := &http.Client{}
			resp, err := client.Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusUnauthorized))
		})
	})
})
