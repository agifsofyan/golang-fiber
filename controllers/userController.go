package controllers

import (
	"context"
	"example/gorest/models"
	"example/gorest/utils"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func Register(c *fiber.Ctx) error {
	var collection = models.UserTable()
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	user := new(models.User)

	if err := c.BodyParser(user); err != nil {
		log.Println(err)
		return utils.FailResponse(c, fiber.StatusBadRequest, "Failed to parse body", err)
	}

	hash, _ := utils.Generate(user.Pass)
	user.Pass = hash
	user.CreatedAt = time.Now()
	user.IsAdmin = false

	result, err := collection.InsertOne(ctx, user)
	if err != nil {
		return utils.FailResponse(c, fiber.StatusInternalServerError, "user failed to insert", err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data":    result,
		"success": true,
		"message": "user inserted successfully",
	})
}

func Me(c *fiber.Ctx) error {
	user, code, msg, err := utils.ShowMe(c)

	if err != nil {
		return utils.FailResponse(c, code, msg, err)
	}

	return utils.SuccessResponse(c, user)
}

func ByEmail(c *fiber.Ctx, email string) error {
	var collection = models.UserTable()

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	findResult := collection.FindOne(ctx, bson.M{"email": email})
	if err := findResult.Err(); err != nil {
		return utils.FailResponse(c, fiber.StatusNotFound, "user Not Found", err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    findResult,
	})
}
