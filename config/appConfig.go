package config

import (
	"errors"
	"github.com/joho/godotenv"
	"hng_stage_two_task/internal/utils"
	"os"
)

type AppConfig struct {
	ServerPort string
	Dsn        string
	AppSecret  string
}

func SetupEnv() (cfg AppConfig, err error) {

	if os.Getenv("APP_ENV") == "dev" {
		godotenv.Load()
	}

	httpPort := os.Getenv("HTTP_PORT")

	if len(httpPort) < 1 {
		return AppConfig{}, errors.New("env variables not found")
	}

	Dsn := os.Getenv("DSN")
	if len(Dsn) < 1 {
		return AppConfig{}, errors.New("env variables not found")
	}

	appSecret := os.Getenv("APP_SECRET")
	if utils.IsEmpty(appSecret) {
		return AppConfig{}, errors.New("app secret not found")
	}

	return AppConfig{ServerPort: httpPort, Dsn: Dsn, AppSecret: appSecret}, nil
}
