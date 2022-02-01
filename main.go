package main

import (
	"log"
	"os"

	"github.com/TesyarRAz/testes/domain/repository"
	"github.com/TesyarRAz/testes/infrastructure/network"
	"github.com/TesyarRAz/testes/infrastructure/persistence"
	"github.com/TesyarRAz/testes/infrastructure/service"
	"github.com/TesyarRAz/testes/interfaces"
	"github.com/TesyarRAz/testes/interfaces/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func init() {
	//To load our environmental variables.
	if err := godotenv.Load(); err != nil {
		log.Println("no env gotten")
	}
}

func main() {
	var (
		// database details
		dbDriver   = os.Getenv("DB_DRIVER")
		dbHost     = os.Getenv("DB_HOST")
		dbPassword = os.Getenv("DB_PASSWORD")
		dbUser     = os.Getenv("DB_USER")
		dbName     = os.Getenv("DB_NAME")
		dbPort     = os.Getenv("DB_PORT")

		// redis details
		redisHost     = os.Getenv("REDIS_HOST")
		redisPort     = os.Getenv("REDIS_PORT")
		redisPassword = os.Getenv("REDIS_PASSWORD")

		// network setup
		database *network.Database
		redis    *network.Redis

		// repository setup
		userRepository repository.UserRepository

		// service setup
		authService  *service.AuthService
		tokenService *service.TokenService

		err error
	)

	if database, err = network.NewDatabase(dbDriver, dbHost, dbPort, dbName, dbUser, dbPassword); err != nil {
		panic(err)
	}

	if err = database.AutoMigrate(); err != nil {
		panic(err)
	}

	redis = network.NewRedis(redisHost, redisPort, redisPassword)

	userRepository = persistence.NewUserRepository(database.Client)

	authService = service.NewAuthService(redis.Client)
	tokenService = service.NewTokenService()

	users := interfaces.NewUsers(userRepository, authService, tokenService)
	auth := interfaces.NewAuth(userRepository, authService, tokenService)

	app := fiber.New(fiber.Config{
		AppName: "Tes Tes",
	})

	app.Use(middleware.CORSMiddleware())

	app.Get("/users", users.Get)

	app.Post("/login", auth.Login)
	app.Post("/register", auth.Register)
	app.Post("/logout", auth.Logout)
	app.Post("/refresh", auth.Refresh)

	port := os.Getenv("APP_PORT")

	if port == "" {
		port = "3000"
	}

	log.Fatal(app.Listen(":" + port))
}
