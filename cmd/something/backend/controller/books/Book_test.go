package books

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"something/config"
	"something/internal/books/application/create"
	"something/internal/books/application/delete"
	"something/internal/books/application/find"
	"something/internal/books/application/update"
	"something/internal/books/domain"
	"something/internal/books/infraestructure/persistence"
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

const userID = "c6facd8d-17f4-43bd-9d90-f4fb024fa2f9"

func setupServer(bookRepo domain.BookRepository) *gin.Engine {
	router := gin.Default()
	finder := find.NewService(bookRepo)
	creator := create.NewService(bookRepo)
	updater := update.NewService(bookRepo)
	deletor := delete.NewService(bookRepo)
	RegisterRoutes(finder, creator, updater, deletor, tokenParams.AccessSecret, router)
	return router
}

func TestBookCheck(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Book Suite")
}

var _ = Describe("Server", func() {
	var server *httptest.Server
	var bookRepo domain.BookRepository

	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	client := config.ConnectBD(dbUser, dbPass, dbHost)

	database := os.Getenv("TEST_DB_NAME")
	dbClient := client.Database(database)

	BeforeEach(func() {
		bookRepo = persistence.NewMongoBookRepository(dbClient)
		server = httptest.NewServer(setupServer(bookRepo))
	})

	AfterEach(func() {
		if err := dbClient.Collection("books").Drop(context.TODO()); err != nil {
			Expect(err).ShouldNot(HaveOccurred())
		}
		server.Close()
	})

	Context("When GET request is sent to /books", func() {
		It("Returns null data if not exists books", func() {
			resp, err := http.Get(server.URL + "/books")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusOK))

			body, err := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(string(body)).To(MatchJSON(`{"data":null}`))
		})

		It("Returns an existing book", func() {
			newBook, _ := domain.NewBook("4c881080-710f-458a-8ec3-058154c47794", "title", "desc", "author", "genre", 1)
			bookRepo.Save(newBook)

			resp, err := http.Get(server.URL + "/books")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusOK))

			body, err := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			Expect(err).ShouldNot(HaveOccurred())

			Expect(string(body)).To(MatchJSON(`
				{
					"data":
						[
							{
								"id":"` + newBook.ID + `",
								"title":"` + newBook.Title + `",
								"description":"` + newBook.Description + `",
								"author":"` + newBook.Author + `",
								"genre":"` + newBook.Genre + `",
								"pages":` + strconv.Itoa(newBook.Pages) + ` ,
								"created_on":"` + newBook.CreatedOn.Format("2006-01-02T15:04:05.999Z07:00") + `"
							}
						]
				}`))
		})
	})
	Context("When GET request by ID is sent to /books/:id", func() {
		It("Returns an existing book by id", func() {
			newBook, _ := domain.NewBook("90cbf21e-f1db-473d-b7b2-6ad77a4ea359", "title", "desc", "author", "genre", 1)
			bookRepo.Save(newBook)

			resp, err := http.Get(server.URL + "/books/" + newBook.ID)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusOK))

			body, err := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			Expect(err).ShouldNot(HaveOccurred())

			Expect(string(body)).To(MatchJSON(`
				{
					"data":
						{
							"id":"` + newBook.ID + `",
							"title":"` + newBook.Title + `",
							"description":"` + newBook.Description + `",
							"author":"` + newBook.Author + `",
							"genre":"` + newBook.Genre + `",
							"pages":` + strconv.Itoa(newBook.Pages) + ` ,
							"created_on":"` + newBook.CreatedOn.Format("2006-01-02T15:04:05.999Z07:00") + `"
						}
				}`))
		})

		It("Returns an 404 status code in non existing id", func() {
			resp, err := http.Get(server.URL + "/books/c0b369a0-8de4-417d-a905-c33644c2907d")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusNotFound))

			body, err := ioutil.ReadAll(resp.Body)
			defer resp.Body.Close()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(string(body)).To(MatchJSON(`{"error":"book not found"}`))
		})
	})
	Context("When PUT request by ID is sent to /books/:id", func() {
		It("Return an 201 status code", func() {
			bookID := "66eb18c5-792b-4897-99e8-86e048132d7b"
			book := map[string]interface{}{
				"title":       "title",
				"description": "description",
				"author":      "author",
				"genre":       "genre",
				"pages":       25,
			}
			jsonReq, err := json.Marshal(book)

			generateAuth, err := jwt.CreateToken(userID, "staff", tokenParams)
			Expect(err).ShouldNot(HaveOccurred())
			req, err := http.NewRequest(
				http.MethodPut,
				server.URL+"/books/"+bookID,
				bytes.NewBuffer(jsonReq))
			req.Header.Set("Content-Type", "application/json; charset=utf-8")
			req.Header.Set("Authorization", "Bearer "+generateAuth.AccessToken)
			client := &http.Client{}

			resp, err := client.Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusCreated))

			createdBook, _ := bookRepo.FindByID(bookID)
			Expect(createdBook).ShouldNot(BeNil())
		})
		It("Returns an 400 status code with an invalid uuid", func() {
			generateAuth, err := jwt.CreateToken(userID, "staff", tokenParams)
			Expect(err).ShouldNot(HaveOccurred())
			req, err := http.NewRequest(http.MethodPut, server.URL+"/books/1", nil)
			req.Header.Set("Authorization", "Bearer "+generateAuth.AccessToken)
			client := &http.Client{}
			resp, err := client.Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusBadRequest))
		})
	})
	Context("When PATCH request by ID is sent to /books/:id", func() {
		It("modify an existing book", func() {
			newBook, _ := domain.NewBook("d14d3e93-4c85-49eb-b6b9-5637c5fcb57c", "title", "desc", "author", "genre", 1)
			bookRepo.Save(newBook)

			fieldsToModify := map[string]interface{}{
				"title":       "title1",
				"description": "description1",
			}
			jsonReq, err := json.Marshal(fieldsToModify)

			generateAuth, err := jwt.CreateToken(userID, "staff", tokenParams)
			Expect(err).ShouldNot(HaveOccurred())

			req, err := http.NewRequest(
				http.MethodPatch,
				server.URL+"/books/"+newBook.ID,
				bytes.NewBuffer(jsonReq),
			)
			req.Header.Set("Content-Type", "application/json; charset=utf-8")
			req.Header.Set("Authorization", "Bearer "+generateAuth.AccessToken)
			client := &http.Client{}

			resp, err := client.Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusOK))

			updatedBook, _ := domain.NewBook(
				newBook.ID,
				fieldsToModify["title"].(string),
				fieldsToModify["description"].(string),
				newBook.Author,
				newBook.Genre,
				newBook.Pages,
			)

			book, _ := bookRepo.FindByID(newBook.ID)
			updatedBook.CreatedOn = book.CreatedOn
			Expect(book).Should(BeEquivalentTo(updatedBook))
		})
	})
	Context("When DELETE request by ID is sent to /books/:id", func() {
		It("delete an existing book", func() {
			newBook, _ := domain.NewBook("567fb602-5533-42a3-8b47-68b474b53e45", "title", "desc", "author", "genre", 1)
			bookRepo.Save(newBook)

			generateAuth, err := jwt.CreateToken(userID, "staff", tokenParams)
			Expect(err).ShouldNot(HaveOccurred())

			req, err := http.NewRequest(http.MethodDelete, server.URL+"/books/567fb602-5533-42a3-8b47-68b474b53e45", nil)
			req.Header.Set("Authorization", "Bearer "+generateAuth.AccessToken)
			client := &http.Client{}
			resp, err := client.Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusNoContent))

			book, _ := bookRepo.FindByID(newBook.ID)
			Expect(book).Should(BeNil())
		})
		It("return an 404 status code in non existing book", func() {
			generateAuth, err := jwt.CreateToken(userID, "staff", tokenParams)
			Expect(err).ShouldNot(HaveOccurred())

			req, err := http.NewRequest(
				http.MethodDelete,
				server.URL+"/books/9b6848af-5e94-44ad-b59c-960c223ee182",
				nil,
			)
			req.Header.Set("Authorization", "Bearer "+generateAuth.AccessToken)
			client := &http.Client{}
			resp, err := client.Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusNotFound))

			body, err := ioutil.ReadAll(resp.Body)
			defer resp.Body.Close()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(string(body)).To(MatchJSON(`{"error":"book not found"}`))
		})
	})
})
