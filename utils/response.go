package utils

import (
	"github.com/gofiber/fiber/v2"
)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func FailResponse(c *fiber.Ctx, status int, msg string, err error) error {
	response := fiber.Map{
		"success": false,
		"message": msg,
		"error":   err,
	}

	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	return c.Status(status).JSON(response)
}

func SuccessResponse(c *fiber.Ctx, result interface{}) error {
	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	return c.Status(fiber.StatusOK).JSON(result)
}
