package utils

import (
	"github.com/gofiber/fiber/v2"
)

type HTTPError struct {
	Success bool   `json:"success" example:false`
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"status bad request"`
}

type PaginateResponse struct {
	Paginate map[string]string `json:"paginate"`
}

type HTTPSuccess struct {
	Success bool              `json:"success" example:true`
	Code    int               `json:"code" example:"200"`
	Message string            `json:"message" example:"success get data"`
	Data    map[string]string `json:"data"`
}

type HTTPSuccessLogin struct {
	AccessToken string `json:"access_token" example: AscdcSdwsddsdsdsnlk.dsdscscscs.wdwdwcwc`
}

func FailResponse(c *fiber.Ctx, code int, msg string) error {
	response := HTTPError{
		Success: false,
		Code:    code,
		Message: msg,
	}

	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	return c.Status(code).JSON(response)
}

func SuccessResponse(c *fiber.Ctx, data fiber.Map, inserted bool) error {
	code := fiber.StatusOK
	if inserted == true {
		code = fiber.StatusCreated
	}

	data["code"] = code

	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	return c.Status(code).JSON(data)
}

func SuccessLoginResponse(c *fiber.Ctx, token string) error {
	response := HTTPSuccessLogin{
		AccessToken: token,
	}

	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	return c.Status(200).JSON(response)
}
