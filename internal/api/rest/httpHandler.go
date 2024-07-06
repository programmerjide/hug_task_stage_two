package rest

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"hng_stage_two_task/config"
	"hng_stage_two_task/internal/helper"
)

type RestHandler struct {
	App    *fiber.App
	DB     *gorm.DB
	Auth   helper.Auth
	Config config.AppConfig
}
