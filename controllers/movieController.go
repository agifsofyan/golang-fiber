package controllers

import (
	"context"
	"example/gorest/models"
	"example/gorest/utils"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetMovies(c *fiber.Ctx) error {
	var collection = models.MovieTable()
	var movies []models.Movie

	filter := bson.M{}
	sorts := [2]string{"created_at", "desc"}

	if s := c.Query("s"); s != "" {
		fields := []string{"title", "description"}
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
		return utils.FailResponse(c, fiber.StatusNotFound, "Movie Not Found", err)
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var movie models.Movie
		cursor.Decode(&movie)
		movies = append(movies, movie)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":      movies,
		"total":     result.Total,
		"page":      result.Page,
		"last_page": result.Last,
		"limit":     result.Limit,
	})
}

func GetMovie(c *fiber.Ctx) error {
	var collection = models.MovieTable()
	var movie models.Movie
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	objId, err := primitive.ObjectIDFromHex(c.Params("id"))
	findResult := collection.FindOne(ctx, bson.M{"_id": objId})

	if err := findResult.Err(); err != nil {
		return utils.FailResponse(c, fiber.StatusNotFound, "Movie Not Found", err)
	}

	err = findResult.Decode(&movie)
	if err != nil {
		return utils.FailResponse(c, fiber.StatusNotFound, "Movie Not Found", err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    movie,
	})
}

func AddMovie(c *fiber.Ctx) error {
	var collection = models.MovieTable()
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	movie := new(models.InsertMovie)

	if err := c.BodyParser(movie); err != nil {
		log.Println(err)
		return utils.FailResponse(c, fiber.StatusBadRequest, "Failed to parse body", err)
	}

	movie.Slug = utils.Slugify(movie.Title)
	movie.CreatedAt = time.Now()
	result, err := collection.InsertOne(ctx, movie)
	if err != nil {
		return utils.FailResponse(c, fiber.StatusInternalServerError, "Movie failed to insert", err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data":    result,
		"success": true,
		"message": "Movie inserted successfully",
	})
}

func UpdateMovie(c *fiber.Ctx) error {
	var collection = models.MovieTable()
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	movie := new(models.Movie)

	if err := c.BodyParser(movie); err != nil {
		log.Println(err)
		return utils.FailResponse(c, fiber.StatusBadRequest, "Failed to parse body", err)
	}

	objId, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return utils.FailResponse(c, fiber.StatusBadRequest, "Movie not found", err)
	}

	update := bson.M{
		"$set": movie,
	}

	_, err = collection.UpdateOne(ctx, bson.M{"_id": objId}, update)

	if err != nil {
		return utils.FailResponse(c, fiber.StatusInternalServerError, "Movie failed to insert", err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Movie updated successfully",
	})
}

func DeleteMovie(c *fiber.Ctx) error {
	var collection = models.MovieTable()
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	objId, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return utils.FailResponse(c, fiber.StatusBadRequest, "Movie not found", err)
	}

	_, err = collection.DeleteOne(ctx, bson.M{"_id": objId})
	if err != nil {
		return utils.FailResponse(c, fiber.StatusInternalServerError, "Movie failed to delete", err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  true,
		"message": "Movie deleted successfully",
	})
}
