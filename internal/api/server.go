package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"hng_stage_two_task/config"
	"hng_stage_two_task/internal/api/rest"
	"hng_stage_two_task/internal/api/rest/handlers"
	"hng_stage_two_task/internal/domain"
	"hng_stage_two_task/internal/helper"
	"log"
)

func StartServer(config config.AppConfig) {
	app := fiber.New()

	// Setup CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:80, http://127.0.0.1:7000", // Specify your allowed origins
		AllowCredentials: true,
	}))

	// connect the ORM
	db, err := gorm.Open(postgres.Open(config.Dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Database connection failed with error: %v\n", err)
	}

	// run migration
	db.AutoMigrate(&domain.User{}, &domain.Organisation{})

	auth := helper.SetupAuth(config.AppSecret)

	restHandler := &rest.RestHandler{
		App:    app,
		DB:     db,
		Auth:   auth,
		Config: config,
	}

	setupRoutes(restHandler)

	app.Listen(config.ServerPort)
}

func setupRoutes(rh *rest.RestHandler) {
	handlers.SetupUserRoutes(rh)
}
