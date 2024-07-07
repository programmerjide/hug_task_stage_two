package test

import (
	"github.com/gofiber/fiber/v2"
	"log"
	"os"
	"testing"

	"hng_stage_two_task/config"
	"hng_stage_two_task/internal/api"
)

var app *fiber.App

func TestMain(m *testing.M) {
	// Setup test environment
	cfg, err := config.SetupEnv()
	if err != nil {
		log.Fatalf("config file is not loaded properly %v\n", err)
	}

	// Start Fiber server for testing
	app = api.StartServer(cfg)

	// Run tests
	code := m.Run()

	// Shutdown Fiber server after tests
	app.Shutdown()
	os.Exit(code)
}
