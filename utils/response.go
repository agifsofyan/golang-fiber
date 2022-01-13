package utils

import (
	"github.com/gofiber/fiber/v2"
)

func FailResponse(c *fiber.Ctx, status int, msg string, err error) error {
	response := fiber.Map{
		"success": false,
		"message": msg,
		"error":   err,
	}

	return c.Status(status).JSON(response)
}

func SuccessResponse(c *fiber.Ctx, result interface{}) error {
	return c.Status(fiber.StatusOK).JSON(result)
}
