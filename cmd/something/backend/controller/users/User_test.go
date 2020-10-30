package users

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"something/internal/users/application/create"
	"something/internal/users/application/delete"
	"something/internal/users/application/find"
	"something/internal/users/application/login"
	"something/internal/users/application/update"
	"something/internal/users/domain"
	"something/internal/users/infraestructure/persistence"
	"something/pkg/crypto"
	jwt "something/pkg/redisjwt"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var tokenParams *jwt.TokenParams = &jwt.TokenParams{
	AccessSecret:  "2TkA87mUUU2pT1j2anRmF72sO",
	RefreshSecret: "xW7xXMWtDv5sDTxEwVFZitjBt",
	AccessTime:    time.Minute * 1,
	RefreshTime:   time.Minute * 1,
}

func TestUserCheck(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "User Suite")
}

func setupServer(
	userRepo domain.UserRepository,
	crypto crypto.Crypto) *gin.Engine {
	router := gin.Default()
	finder := find.NewService(userRepo)
	creator := create.NewService(userRepo, crypto)
	updater := update.NewService(userRepo)
	deleter := delete.NewService(userRepo)
	authLogin := login.NewService(userRepo, crypto)
	RegisterRoutes(finder, creator, updater, deleter, authLogin, tokenParams, router)
	return router
}

var _ = Describe("Server", func() {
	var server *httptest.Server
	var userRepo domain.UserRepository
	var cryptoRepo crypto.Crypto

	BeforeEach(func() {
		userRepo = persistence.NewInMemoryUserRepository()
		cryptoRepo = crypto.NewBcrypt()
		server = httptest.NewServer(setupServer(userRepo, cryptoRepo))
	})

	AfterEach(func() {
		server.Close()
	})

	Context("When GET request is sent to /users/", func() {
		It("Returns null data if not exists users", func() {
			resp, err := http.Get(server.URL + "/users/")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusOK))

			body, err := ioutil.ReadAll(resp.Body)
			defer resp.Body.Close()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(string(body)).To(MatchJSON(`{"data":null}`))
		})
		It("Returns an existing user", func() {
			newUser, _ := domain.NewUser("6adbcea4-4fd4-45eb-8803-6c8474ac663a", "bob", "bob", "bob@mail.com", "bob123")
			userRepo.Save(newUser)

			resp, err := http.Get(server.URL + "/users/")
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
							"id":"` + newUser.ID + `",
							"name":"` + newUser.Name + `",
							"username":"` + newUser.Username + `",
							"email":"` + newUser.Email + `" ,
							"created_on":"` + newUser.CreatedOn.Format("2006-01-02T15:04:05.999999999Z07:00") + `"
						}
					]
			}`))
		})
	})
	Context("When GET request by ID is sent to /users/:id", func() {
		It("Returns an existing user by id", func() {
			newUser, _ := domain.NewUser("03d0b376-046f-415c-85d5-c4f102645835", "alice", "alice", "alice@mail.com", "alice123")
			userRepo.Save(newUser)

			resp, err := http.Get(
				server.URL + "/users/03d0b376-046f-415c-85d5-c4f102645835")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusOK))

			body, err := ioutil.ReadAll(resp.Body)
			defer resp.Body.Close()
			Expect(err).ShouldNot(HaveOccurred())

			Expect(string(body)).To(MatchJSON(`
			{
				"data":
					{
						"id":"` + newUser.ID + `",
						"name":"` + newUser.Name + `",
						"username":"` + newUser.Username + `",
						"email":"` + newUser.Email + `" ,
						"created_on":"` + newUser.CreatedOn.Format("2006-01-02T15:04:05.999999999Z07:00") + `"
					}
			}`))
		})
		It("Returns an 404 status code in non existing id", func() {
			resp, err := http.Get(
				server.URL + "/users/9e3bea73-3f38-4d02-9e70-fca95154e782")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusNotFound))

			body, err := ioutil.ReadAll(resp.Body)
			defer resp.Body.Close()
			Expect(err).ShouldNot(HaveOccurred())

			Expect(string(body)).To(MatchJSON(`{"error":"user not found"}`))
		})
	})

	Context("When PUT request by ID is sent to /users/:id", func() {
		It("Return an 201 status code", func() {
			user := map[string]interface{}{
				"name":     "john",
				"username": "john1",
				"email":    "john@example.com",
				"password": "super-secure-password",
			}
			jsonReq, err := json.Marshal(user)
			req, err := http.NewRequest(
				http.MethodPut,
				server.URL+"/users/dc4fc484-a281-463c-bdd0-6adfa2167931",
				bytes.NewBuffer(jsonReq))
			req.Header.Set("Content-Type", "application/json; charset=utf-8")
			client := &http.Client{}

			resp, err := client.Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusCreated))

			createdUser, _ := userRepo.FindByID("dc4fc484-a281-463c-bdd0-6adfa2167931")
			Expect(createdUser).ShouldNot(BeNil())
		})
		It("Returns an 404 status code with an invalid uuid", func() {
			req, err := http.NewRequest(http.MethodPut, server.URL+"/users/1", nil)
			client := &http.Client{}
			resp, err := client.Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusBadRequest))
		})
		It("Returns an 400 status code with an existing email", func() {
			newUser, _ := domain.NewUser(
				"a9f6b942-6289-4037-938b-5e19d9262443",
				"tom",
				"tom01",
				"tom@example.com",
				"super-password")
			userRepo.Save(newUser)

			user := map[string]interface{}{
				"name":     "tom",
				"username": "tom05",
				"email":    "tom@example.com",
				"password": "secure-password",
			}
			jsonReq, err := json.Marshal(user)
			req, err := http.NewRequest(
				http.MethodPut,
				server.URL+"/users/975e0aa8-a78e-4cf7-a2b8-80855314f768",
				bytes.NewBuffer(jsonReq))
			req.Header.Set("Content-Type", "application/json; charset=utf-8")
			client := &http.Client{}

			resp, err := client.Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusBadRequest))

			body, err := ioutil.ReadAll(resp.Body)
			defer resp.Body.Close()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(string(body)).To(MatchJSON(`{"error":"email already in use"}`))
		})
		It("Returns an 400 status code with an existing username", func() {
			newUser, _ := domain.NewUser(
				"ac7b7cfb-8448-44df-b12a-f59273682eeb",
				"martin",
				"martin01",
				"martin@example.com",
				"super-master-password")
			userRepo.Save(newUser)

			user := map[string]interface{}{
				"name":     "mart",
				"username": "martin01",
				"email":    "mart@example.com",
				"password": "secure-mega-password",
			}
			jsonReq, err := json.Marshal(user)
			req, err := http.NewRequest(
				http.MethodPut,
				server.URL+"/users/b9bd8896-6aec-4058-82a3-865087560b91",
				bytes.NewBuffer(jsonReq))
			req.Header.Set("Content-Type", "application/json; charset=utf-8")
			client := &http.Client{}

			resp, err := client.Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusBadRequest))

			body, err := ioutil.ReadAll(resp.Body)
			defer resp.Body.Close()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(string(body)).To(MatchJSON(`{"error":"username already in use"}`))
		})
	})

	Context("When PATCH request by ID is sent to /users/:id", func() {
		It("modify an existing user", func() {
			newUser, _ := domain.NewUser(
				"d14d3e93-4c85-49eb-b6b9-5637c5fcb57c",
				"Martin",
				"martin01",
				"martin@example.com",
				"super-ultra-password",
			)
			userRepo.Save(newUser)

			fieldsToModify := map[string]interface{}{
				"name": "Martin Cooper",
			}
			jsonReq, err := json.Marshal(fieldsToModify)

			generateAuth, err := jwt.CreateToken(newUser.ID, tokenParams)
			Expect(err).ShouldNot(HaveOccurred())
			req, err := http.NewRequest(
				http.MethodPatch,
				server.URL+"/users/"+newUser.ID,
				bytes.NewBuffer(jsonReq))
			req.Header.Set("Content-Type", "application/json; charset=utf-8")
			req.Header.Set("Authorization", "Bearer "+generateAuth.AccessToken)
			client := &http.Client{}

			resp, err := client.Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusOK))

			updatedUser, _ := domain.NewUser(
				newUser.ID,
				fieldsToModify["name"].(string),
				newUser.Username,
				newUser.Email,
				newUser.Password,
			)
			updatedUser.CreatedOn = newUser.CreatedOn

			user, _ := userRepo.FindByID(newUser.ID)
			Expect(user).Should(BeEquivalentTo(updatedUser))
		})
	})

	Context("When DELETE request by ID is sent to /users/:id", func() {
		It("delete an existing user", func() {
			newUser, _ := domain.NewUser(
				"552394d5-620c-4b7b-99da-c95fe5e52730",
				"madison", "madison1", "madison@example.com",
				"secret-pass-1")
			userRepo.Save(newUser)

			generateAuth, err := jwt.CreateToken(newUser.ID, tokenParams)
			Expect(err).ShouldNot(HaveOccurred())

			req, err := http.NewRequest(
				http.MethodDelete,
				server.URL+"/users/"+newUser.ID, nil)
			req.Header.Set("Authorization", "Bearer "+generateAuth.AccessToken)
			client := &http.Client{}
			resp, err := client.Do(req)

			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusNoContent))

			user, _ := userRepo.FindByID(newUser.ID)
			Expect(user).Should(BeNil())
		})
		It("return an 404 status code in non existing user", func() {
			userID := "9b6848af-5e94-44ad-b59c-960c223ee182"

			generateAuth, err := jwt.CreateToken(userID, tokenParams)
			Expect(err).ShouldNot(HaveOccurred())

			req, err := http.NewRequest(
				http.MethodDelete,
				server.URL+"/users/"+userID, nil)
			req.Header.Set("Authorization", "Bearer "+generateAuth.AccessToken)
			client := &http.Client{}
			resp, err := client.Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusNotFound))

			body, err := ioutil.ReadAll(resp.Body)
			defer resp.Body.Close()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(string(body)).To(MatchJSON(`{"error":"user not found"}`))
		})
	})

	Context("When POST request is sent to /login", func() {
		It("authenticate existing user", func() {
			hash, _ := cryptoRepo.Hash("secret-pass-1")
			newUser, _ := domain.NewUser(
				"552394d5-620c-4b7b-99da-c95fe5e52730",
				"madison", "madison1", "madison@example.com",
				hash)
			userRepo.Save(newUser)

			loginFields := map[string]interface{}{
				"email":    "madison@example.com",
				"password": "secret-pass-1",
			}
			jsonReq, err := json.Marshal(loginFields)
			req, err := http.NewRequest(
				http.MethodPost,
				server.URL+"/login", bytes.NewBuffer(jsonReq))
			client := &http.Client{}
			resp, err := client.Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusOK))

		})
		It("return an 404 status code in non existing email", func() {
			hash, _ := cryptoRepo.Hash("pep-secret-pass")
			newUser, _ := domain.NewUser(
				"68a004c8-e1c1-49c0-a430-66f9cf6fd1ad",
				"Pep", "peep", "pep@gmail.com",
				hash)
			userRepo.Save(newUser)

			loginFields := map[string]interface{}{
				"email":    "madison@example.com",
				"password": "madison-11",
			}
			jsonReq, err := json.Marshal(loginFields)
			req, err := http.NewRequest(
				http.MethodPost,
				server.URL+"/login", bytes.NewBuffer(jsonReq))
			client := &http.Client{}
			resp, err := client.Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusNotFound))

			body, err := ioutil.ReadAll(resp.Body)
			defer resp.Body.Close()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(string(body)).To(MatchJSON(`{"error":"email not found"}`))
		})
		It("return an 401 status code in invalid email/password", func() {
			hash, _ := cryptoRepo.Hash("pep-secret-pass")
			newUser, _ := domain.NewUser(
				"68a004c8-e1c1-49c0-a430-66f9cf6fd1ad",
				"Pep", "peep", "pep@gmail.com",
				hash)
			userRepo.Save(newUser)

			loginFields := map[string]interface{}{
				"email":    "pep@gmail.com",
				"password": "pep-secret-pass-10",
			}
			jsonReq, err := json.Marshal(loginFields)
			req, err := http.NewRequest(
				http.MethodPost,
				server.URL+"/login", bytes.NewBuffer(jsonReq))
			client := &http.Client{}
			resp, err := client.Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusUnauthorized))
		})
	})
})
