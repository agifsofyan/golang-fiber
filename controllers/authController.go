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

// models.SwaggerLogin godoc
// @Summary      Login
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        auth  body     models.SwaggerLogin  true  "Login"
// @Router       /auth [post]
func Login(c *fiber.Ctx) error {
	var collection = models.UserTable()
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	var body models.SwaggerLogin
	var user models.User

	err := c.BodyParser(&body)
	if err != nil {
		return utils.FailResponse(c, fiber.StatusBadRequest, "Cannot parse JSON")
	}

	// Throws Unauthorized error
	if body.Email == "" || body.Pass == "" {
		return utils.FailResponse(c, fiber.StatusBadRequest, "email & password wajib")
	}

	findResult := collection.FindOne(ctx, bson.M{"email": body.Email})

	if err := findResult.Err(); err != nil {
		return utils.FailResponse(c, fiber.StatusNotFound, "User Not Found")
	}

	errno := findResult.Decode(&user)
	if errno != nil {
		return utils.FailResponse(c, fiber.StatusNotFound, "Error decode user")
	}

	compare := utils.Compare(user.Pass, body.Pass)

	if compare != true {
		return utils.FailResponse(c, fiber.StatusBadRequest, "Password not match")
	}

	token, err := utils.GenerateToken(user)

	if err != nil {
		return utils.FailResponse(c, fiber.StatusBadRequest, "Failed to generate token")
	}

	return utils.SuccessLoginResponse(c, token["access_token"])
}

func Accessible(c *fiber.Ctx) error {
	return utils.FailResponse(c, fiber.StatusNotAcceptable, "Cannot Access")
}

func Restricted(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)

	return utils.SuccessResponse(c, fiber.Map{"message": "Welcome " + name}, false)
}
