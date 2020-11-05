package userfollow

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"something/config"
	"something/internal/userfollow/application/find"
	"something/internal/userfollow/application/followers"
	"something/internal/userfollow/domain"
	"something/internal/userfollow/infraestructure/persistence"
	userFind "something/internal/users/application/find"
	userDomain "something/internal/users/domain"
	userPersistence "something/internal/users/infraestructure/persistence"
	"testing"
	"time"

	jwt "something/pkg/redisjwt"

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

func TestUserFollowCheck(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "User Follow Suite")
}

func setupServer(
	userFollowRepo domain.UserFollowRepository,
	userRepo userDomain.UserRepository) *gin.Engine {
	router := gin.Default()
	userFinder := userFind.NewService(userRepo)
	finder := find.NewService(userFollowRepo)
	follow := followers.NewService(userFollowRepo)
	RegisterRoutes(finder, userFinder, follow, tokenParams, router)
	return router
}

var _ = Describe("Server", func() {
	var server *httptest.Server
	var userFollowRepo domain.UserFollowRepository
	var userRepo userDomain.UserRepository

	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	client := config.ConnectBD(dbUser, dbPass, dbHost)

	database := os.Getenv("TEST_DB_NAME")
	dbClient := client.Database(database)

	BeforeEach(func() {
		userFollowRepo = persistence.NewMongoUserFollowRepository(dbClient)
		userRepo = userPersistence.NewInMemoryUserRepository()
		server = httptest.NewServer(setupServer(userFollowRepo, userRepo))
	})

	AfterEach(func() {
		if err := dbClient.Collection("user_follows").Drop(context.TODO()); err != nil {
			Expect(err).ShouldNot(HaveOccurred())
		}
		server.Close()
	})

	Context("When GET request is sent to /users/:id/followers", func() {
		It("Returns null data if not exists followers", func() {
			newUser, _ := userDomain.NewUser(
				"8d4eb934-8116-4b2f-bd9d-2b6134a6a6f9",
				"dante", "dante06", "dante@gmail.com",
				"dante-secure-password")
			userRepo.Save(newUser)

			resp, err := http.Get(server.URL + "/users/" + newUser.ID + "/followers")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusOK))

			body, err := ioutil.ReadAll(resp.Body)
			defer resp.Body.Close()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(string(body)).To(MatchJSON(`{"data":null}`))
		})
		It("Returns an existing follower", func() {
			newUser, _ := userDomain.NewUser(
				"7e57186c-ab07-4372-8ae1-efea970edd50",
				"dante", "dante06", "dante@gmail.com",
				"dante-secure-password")
			userRepo.Save(newUser)
			userFollow, _ := domain.NewUserFollow(
				"a6e31ea4-af01-4426-a89f-98a14cf2b077",
				newUser.ID)
			userFollowRepo.Follow(userFollow)

			resp, err := http.Get(server.URL + "/users/" + newUser.ID + "/followers")
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
							"from":"` + userFollow.From + `",
							"to":"` + userFollow.To + `",
							"created_on":"` + userFollow.CreatedOn.Format("2006-01-02T15:04:05.999Z07:00") + `"
						}
					]
			}`))
		})
	})

	Context("When GET request is sent to /users/:id/following", func() {
		It("Returns null data if not exists following", func() {
			newUser, _ := userDomain.NewUser(
				"1b995593-812f-411e-ad6f-8bd4ff22fb98",
				"dante", "dante06", "dante@gmail.com",
				"dante-secure-password")
			userRepo.Save(newUser)

			resp, err := http.Get(server.URL + "/users/" + newUser.ID + "/following")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusOK))

			body, err := ioutil.ReadAll(resp.Body)
			defer resp.Body.Close()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(string(body)).To(MatchJSON(`{"data":null}`))
		})
		It("Returns an existing following", func() {
			newUser, _ := userDomain.NewUser(
				"3e2db844-53d9-4b48-8e02-57d1d079da7e",
				"dante", "dante06", "dante@gmail.com",
				"dante-secure-password")
			userRepo.Save(newUser)
			userFollow, _ := domain.NewUserFollow(newUser.ID, "4328edff-5422-46eb-b7d6-2b5bf89cb151")
			userFollowRepo.Follow(userFollow)

			resp, err := http.Get(server.URL + "/users/" + newUser.ID + "/following")
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
							"from":"` + userFollow.From + `",
							"to":"` + userFollow.To + `",
							"created_on":"` + userFollow.CreatedOn.Format("2006-01-02T15:04:05.999Z07:00") + `"
						}
					]
			}`))
		})
	})

	Context("When POST request is sent to /user/follow/:id", func() {
		It("follow an existing user", func() {
			newUser, _ := userDomain.NewUser(
				"552394d5-620c-4b7b-99da-c95fe5e52730",
				"madison", "madison1", "madison@example.com",
				"super-secure-password")
			anotherUser, _ := userDomain.NewUser(
				"03de8950-2a96-453c-bc71-29e7487a55cd",
				"james", "james1", "james@example.com",
				"super-strong-password")
			userRepo.Save(newUser)

			generateAuth, err := jwt.CreateToken(anotherUser.ID, anotherUser.Role, tokenParams)
			req, err := http.NewRequest(
				http.MethodPost,
				server.URL+"/user/follow/"+newUser.ID,
				nil,
			)
			req.Header.Set("Authorization", "Bearer "+generateAuth.AccessToken)
			client := &http.Client{}
			resp, err := client.Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusOK))
		})
		It("return an 404 status code in non existing user", func() {
			nonExistingUserID := "5b9eb022-6445-43d4-8e35-b0e9a6e97275"
			newUser, _ := userDomain.NewUser(
				"68a004c8-e1c1-49c0-a430-66f9cf6fd1ad",
				"dante", "dante06", "dante@gmail.com",
				"dante-secure-password")
			userRepo.Save(newUser)

			generateAuth, err := jwt.CreateToken(newUser.ID, newUser.Role, tokenParams)
			req, err := http.NewRequest(
				http.MethodPost,
				server.URL+"/user/follow/"+nonExistingUserID,
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
			Expect(string(body)).To(MatchJSON(`{"error":"user not found"}`))
		})
	})
	Context("When POST request is sent to /user/unfollow/:id", func() {
		It("unfollow an existing user", func() {
			userID := "e1f4be08-bd64-4e2c-8dd2-4ec0545cce0e"
			newUser, _ := userDomain.NewUser(
				"88d5c1cd-3367-446e-8716-984b6e28d984",
				"dante", "dante06", "dante@gmail.com",
				"dante-secure-password")
			userRepo.Save(newUser)
			userFollow, _ := domain.NewUserFollow(userID, newUser.ID)
			userFollowRepo.Follow(userFollow)

			generateAuth, err := jwt.CreateToken(userID, "default", tokenParams)
			req, err := http.NewRequest(
				http.MethodPost,
				server.URL+"/user/unfollow/"+newUser.ID,
				nil,
			)
			req.Header.Set("Authorization", "Bearer "+generateAuth.AccessToken)
			client := &http.Client{}
			resp, err := client.Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusOK))
		})
		It("return an 404 status code in non existing user", func() {
			nonExistingUserID := "5b9eb022-6445-43d4-8e35-b0e9a6e97275"
			newUser, _ := userDomain.NewUser(
				"68a004c8-e1c1-49c0-a430-66f9cf6fd1ad",
				"dante", "dante06", "dante@gmail.com",
				"dante-secure-password")
			userRepo.Save(newUser)

			generateAuth, err := jwt.CreateToken(newUser.ID, newUser.Role, tokenParams)
			req, err := http.NewRequest(
				http.MethodPost,
				server.URL+"/user/unfollow/"+nonExistingUserID,
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
			Expect(string(body)).To(MatchJSON(`{"error":"user not found"}`))
		})
	})

})
