package main

import (
	"log"
	"os"
	"something/cmd/something/backend/controller/healthcheck"
	"something/config"
	"something/pkg/crypto"
	jwt "something/pkg/redisjwt"
	"time"

	"something/cmd/something/backend/controller/bookreviews"
	"something/internal/bookreviews/application/create"
	"something/internal/bookreviews/application/delete"
	"something/internal/bookreviews/application/find"
	"something/internal/bookreviews/application/update"
	"something/internal/bookreviews/infraestructure/persistence"

	"something/cmd/something/backend/controller/books"
	bookCreate "something/internal/books/application/create"
	bookDelete "something/internal/books/application/delete"
	bookFinder "something/internal/books/application/find"
	bookUpdate "something/internal/books/application/update"
	bookPersistance "something/internal/books/infraestructure/persistence"

	"something/cmd/something/backend/controller/users"
	userCreate "something/internal/users/application/create"
	userDelete "something/internal/users/application/delete"
	userFinder "something/internal/users/application/find"
	"something/internal/users/application/login"
	userUpdate "something/internal/users/application/update"
	userPersistance "something/internal/users/infraestructure/persistence"

	"something/cmd/something/backend/controller/userfollow"
	userFollowFinder "something/internal/userfollow/application/find"
	userFollow "something/internal/userfollow/application/followers"
	userFollowPersistance "something/internal/userfollow/infraestructure/persistence"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file")
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	gin.SetMode(gin.ReleaseMode)

	log.Println("This is something app!")

	setupServer().Run("localhost:" + port)
	return
}

func setupServer() *gin.Engine {

	tokenParams := &jwt.TokenParams{
		AccessSecret:  os.Getenv("ACCESS_SECRET"),
		RefreshSecret: os.Getenv("REFRESH_SECRET"),
		AccessTime:    time.Hour * 24,
		RefreshTime:   time.Hour * 24 * 7,
	}

	router := gin.Default()

	router.Use(gin.Recovery())
	router.Use(cors.Default())

	// init database
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	client := config.ConnectBD(dbUser, dbPass, dbHost)

	database := os.Getenv("DB_NAME")
	dbClient := client.Database(database)

	// init crypto
	cryptoRepo := crypto.NewBcrypt()

	// Books
	inMemoryBookRepo := bookPersistance.NewMongoBookRepository(dbClient)
	bookFind := bookFinder.NewService(inMemoryBookRepo)
	bookCreator := bookCreate.NewService(inMemoryBookRepo)
	bookDeletor := bookDelete.NewService(inMemoryBookRepo)
	bookUpdater := bookUpdate.NewService(inMemoryBookRepo)
	books.RegisterRoutes(bookFind, bookCreator, bookUpdater, bookDeletor, tokenParams.AccessSecret, router)

	// Book reviews
	inMemoryBookReviewRepo := persistence.NewMongoBookReviewRepository(dbClient)
	bookReviewFinder := find.NewService(inMemoryBookReviewRepo)
	bookReviewCreator := create.NewService(inMemoryBookReviewRepo)
	bookReviewUpdater := update.NewService(inMemoryBookReviewRepo)
	bookReviewDelete := delete.NewService(inMemoryBookReviewRepo)
	bookreviews.RegisterRoutes(bookReviewFinder, bookFind, bookReviewCreator, bookReviewUpdater, bookReviewDelete, tokenParams.AccessSecret, router)

	// Users
	inMemoryUserRepo := userPersistance.NewMongoUsersRepository(dbClient)
	userFind := userFinder.NewService(inMemoryUserRepo)
	userCreator := userCreate.NewService(inMemoryUserRepo, cryptoRepo)
	userUpdater := userUpdate.NewService(inMemoryUserRepo)
	userDeletor := userDelete.NewService(inMemoryUserRepo)
	authLogin := login.NewService(inMemoryUserRepo, cryptoRepo)
	users.RegisterRoutes(userFind, bookFind, userCreator, userUpdater, userDeletor, authLogin, tokenParams, router)

	// Users followers
	inMemoryUserFollowRepo := userFollowPersistance.NewMongoUserFollowRepository(dbClient)
	userFollowFind := userFollowFinder.NewService(inMemoryUserFollowRepo)
	userFollower := userFollow.NewService(inMemoryUserFollowRepo)
	userfollow.RegisterRoutes(userFollowFind, userFind, userFollower, tokenParams, router)

	// Health-check
	healthcheck.RegisterRoutes(router)

	return router
}
