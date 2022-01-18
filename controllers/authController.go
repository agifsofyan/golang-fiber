package controllers

import (
	"context"
	"example/gorest/models"
	"example/gorest/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
)

func Login(c *fiber.Ctx) error {
	var collection = models.UserTable()
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	type request struct {
		Email string `json:"email"`
		Pass  string `json:"pass"`
	}

	var body request
	var user models.User

	err := c.BodyParser(&body)
	if err != nil {
		return utils.FailResponse(c, fiber.StatusBadRequest, "Cannot parse JSON", err)
	}

	// Throws Unauthorized error
	if body.Email == "" || body.Pass == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "email & password wajib"})
	}

	findResult := collection.FindOne(ctx, bson.M{"email": body.Email})

	if err := findResult.Err(); err != nil {
		return utils.FailResponse(c, fiber.StatusNotFound, "User Not Found", err)
	}

	errno := findResult.Decode(&user)
	if errno != nil {
		return utils.FailResponse(c, fiber.StatusNotFound, "Error decode user", err)
	}

	compare := utils.Compare(user.Pass, body.Pass)

	if compare != true {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Password not match"})
	}

	token, err := utils.GenerateToken(user)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to generate token"})
	}

	return c.JSON(token)
}

func Accessible(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotAcceptable).JSON(fiber.Map{"message": "Cannot Access"})
}

func Restricted(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Welcome " + name})
}
