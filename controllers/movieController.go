package controllers

import (
	"context"
	"fmt"
	"golang-fiber/gorest/models"
	"golang-fiber/gorest/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// @Summary      Show all movie
// @Tags         movies
// @Accept       json
// @Produce      json
// @Param        s    query     string  false  "search by s"
// @Param        page    query     string  false  "number page from limited data"
// @Param        limit    query     string  false  "count rendered data"
// @Param        sortby    query     string  false  "key sort the data"
// @Param        sortval    query     string  false  "value sort the data"
// @Security 	 ApiKeyAuth
// @Router       /movies [get]
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
		return utils.FailResponse(c, fiber.StatusNotFound, "Movie Not Found")
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var movie models.Movie
		cursor.Decode(&movie)
		movies = append(movies, movie)
	}

	return utils.SuccessResponse(c, fiber.Map{
		"success":   true,
		"code":      fiber.StatusOK,
		"message":   "Success get data",
		"data":      movies,
		"total":     result.Total,
		"page":      result.Page,
		"last_page": result.Last,
		"limit":     result.Limit,
	}, false)
}

// @Summary      Show an movie
// @Description  get string by ID
// @Tags         movies
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Movie ID"
// @Security 	 ApiKeyAuth
// @Router       /movies/{id} [get]
func GetMovie(c *fiber.Ctx) error {
	var collection = models.MovieTable()

	id := c.Params("id")
	lookup := bson.M{"collection": "genres", "field": "genre"}

	findResult, msg, code := utils.FindById(collection, id, lookup)

	if code != 200 {
		return utils.FailResponse(c, code, msg)
	}

	return utils.SuccessResponse(c, fiber.Map{
		"success": true,
		"message": "Success get data",
		"data":    findResult,
	}, false)
}

// @Summary      Add new movie
// @Tags         movies
// @Accept       json
// @Produce      json
// @Param        movie  body    models.SwaggerInsertMovie  true  "Add new movie"
// @Security 	 ApiKeyAuth
// @Router       /movies [post]
func AddMovie(c *fiber.Ctx) error {
	var collection = models.MovieTable()
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	movie := new(models.InsertMovie)

	if err := c.BodyParser(movie); err != nil {
		return utils.FailResponse(c, fiber.StatusBadRequest, "Failed to parse body")
	}

	type Request struct {
		File string `json:"file"`
	}

	var body Request

	err := c.BodyParser(&body)
	if err := c.BodyParser(movie); err != nil {
		return utils.FailResponse(c, fiber.StatusBadRequest, "Failed to parse body")
	}

	file := fmt.Sprint(body) // convert to string

	err, path, filename := utils.WriteBase64ToFile(file)
	if err != nil {
		return utils.FailResponse(c, 400, string(err.Error()))
	}

	thumb, err := utils.CreateThumbnail(path, filename)
	if err != nil {
		return utils.FailResponse(c, 400, string(err.Error()))
	}

	img := make(map[string]interface{})
	img["path"] = fmt.Sprintf("%s/%s", path, filename)
	img["thumbnail"] = thumb

	movie.Slug = utils.Slugify(movie.Title)
	movie.Img = img
	movie.CreatedAt = time.Now()

	result, err := collection.InsertOne(ctx, movie)
	if err != nil {
		return utils.FailResponse(c, fiber.StatusInternalServerError, "Movie failed to insert")
	}

	return utils.SuccessResponse(c, fiber.Map{
		"success":    true,
		"message":    "Movie inserted successfully",
		"insertedID": result.InsertedID,
	}, true)
}

// @Summary      Update movie
// @Tags         movies
// @Accept       json
// @Produce      json
// @Param        id  path    string  true  "Movie ID"
// @Param        movie  body    models.SwaggerInsertMovie  true  "Update movie"
// @Security 	 ApiKeyAuth
// @Router       /movies/{id} [put]
func UpdateMovie(c *fiber.Ctx) error {
	var collection = models.MovieTable()
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	movie := new(models.Movie)

	if err := c.BodyParser(movie); err != nil {
		return utils.FailResponse(c, fiber.StatusBadRequest, "Failed to parse body")
	}

	objId, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return utils.FailResponse(c, fiber.StatusBadRequest, "Movie not found")
	}

	update := bson.M{
		"$set": movie,
	}

	_, err = collection.UpdateOne(ctx, bson.M{"_id": objId}, update)

	if err != nil {
		return utils.FailResponse(c, fiber.StatusInternalServerError, "Movie failed to insert")
	}

	return utils.SuccessResponse(c, fiber.Map{
		"success": true,
		"message": "Movie updated successfully",
	}, false)
}

// @Summary      Update movie
// @Tags         movies
// @Accept       json
// @Produce      json
// @Param        id  path    string  true  "Movie ID"
// @Security 	 ApiKeyAuth
// @Router       /movies/{id} [delete]
func DeleteMovie(c *fiber.Ctx) error {
	var collection = models.MovieTable()
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	objId, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return utils.FailResponse(c, fiber.StatusBadRequest, "Movie not found")
	}

	_, err = collection.DeleteOne(ctx, bson.M{"_id": objId})
	if err != nil {
		return utils.FailResponse(c, fiber.StatusInternalServerError, "Movie failed to delete")
	}

	return utils.SuccessResponse(c, fiber.Map{
		"status":  true,
		"message": "Movie deleted successfully",
	}, false)
}
