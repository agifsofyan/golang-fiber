package utils

import (
	"context"
	"golang-fiber/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Payload(c *fiber.Ctx, keyword string) jwt.MapClaims {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	return claims
}

func ShowMe(c *fiber.Ctx) (user models.User, code int, msg string, err error) {
	var collection = models.UserTable()
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	claim := Payload(c, "uuid")
	uuid := claim["uuid"].(string)
	userId, _ := primitive.ObjectIDFromHex(uuid)
	findResult := collection.FindOne(ctx, bson.M{"_id": userId})

	if err := findResult.Err(); err != nil {
		return user, fiber.StatusBadGateway, "User Not Found", err
	}

	errno := findResult.Decode(&user)
	if errno != nil {
		return user, fiber.StatusBadGateway, "Error decode user", err
	}

	return user, fiber.StatusOK, "Success", nil
}
