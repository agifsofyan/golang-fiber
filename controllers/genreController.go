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

func Index(c *fiber.Ctx) error {
	var collection = models.GenreTable()
	var list []models.Genre

	filter := bson.M{}
	sorts := [2]string{"created_at", "desc"}

	if s := c.Query("s"); s != "" {
		fields := []string{"name"}
		filter = utils.Search(s, fields)
	}

	if sb := c.Query("sortby"); sb != "" {
		sorts[0] = sb
	}

	if sv := c.Query("sortval"); sv != "" {
		sorts[1] = sv
	}

	cursor, result, ctx, err := utils.Paginate(c, collection, filter, sorts)

	if err != nil {
		return utils.FailResponse(c, fiber.StatusNotFound, "Genre Not Found", err)
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var genre models.Genre
		cursor.Decode(&genre)
		list = append(list, genre)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":      list,
		"total":     result.Total,
		"page":      result.Page,
		"last_page": result.Last,
		"limit":     result.Limit,
	})
}

func Add(c *fiber.Ctx) error {
	var collection = models.GenreTable()
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	genre := new(models.Genre)

	if err := c.BodyParser(genre); err != nil {
		log.Println(err)
		return utils.FailResponse(c, fiber.StatusBadRequest, "Failed to parse body", err)
	}

	genre.Slug = utils.Slugify(genre.Name)
	result, err := collection.InsertOne(ctx, genre)
	if err != nil {
		return utils.FailResponse(c, fiber.StatusInternalServerError, "Genre failed to insert", err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data":    result,
		"success": true,
		"message": "Genre inserted successfully",
	})
}
