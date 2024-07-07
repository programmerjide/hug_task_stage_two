package helper

import (
	"github.com/gofiber/fiber/v2"
	"hng_stage_two_task/internal/dto"
)

func RespondWithError(ctx *fiber.Ctx, status int, code string, message string) error {
	response := dto.DefaultApiResponse{
		BaseResponse: dto.BaseResponse[any]{
			Status:  code,
			Message: message,
			Data:    nil,
		},
	}
	return ctx.Status(status).JSON(response)
}
