package controllers

import (
	"context"
	"example/gorest/models"
	"example/gorest/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

// @Summary      Show all genre
// @Tags         genres
// @Accept       json
// @Produce      json
// @Param        s    query     string  false  "search by s"
// @Param        page    query     string  false  "number page from limited data"
// @Param        limit    query     string  false  "count rendered data"
// @Param        sortby    query     string  false  "key sort the data"
// @Param        sortval    query     string  false  "value sort the data"
// @Security 	 ApiKeyAuth
// @Router       /genres [get]
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
		return utils.FailResponse(c, fiber.StatusNotFound, "Genre Not Found")
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var genre models.Genre
		cursor.Decode(&genre)
		list = append(list, genre)
	}

	response := fiber.Map{
		"success":   true,
		"code":      fiber.StatusOK,
		"message":   "Success get data",
		"data":      list,
		"total":     result.Total,
		"page":      result.Page,
		"last_page": result.Last,
		"limit":     result.Limit,
	}

	return utils.SuccessResponse(c, response, false)
}

// @Summary      Add new genre
// @Tags         genres
// @Accept       json
// @Produce      json
// @Param        genre  body    models.SwaggerInsertGenre  true  "Add new genre"
// @Security 	 ApiKeyAuth
// @Router       /genres [post]
func Add(c *fiber.Ctx) error {
	var collection = models.GenreTable()
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	genre := new(models.Genre)

	if err := c.BodyParser(genre); err != nil {
		return utils.FailResponse(c, fiber.StatusBadRequest, "Failed to parse body")
	}

	genre.Slug = utils.Slugify(genre.Name)
	result, err := collection.InsertOne(ctx, genre)
	if err != nil {
		return utils.FailResponse(c, fiber.StatusInternalServerError, "Genre failed to insert")
	}

	return utils.SuccessResponse(c, fiber.Map{
		"success":    true,
		"message":    "Genre inserted successfully",
		"insertedID": result.InsertedID,
	}, true)
}

// @Summary      Show an genre
// @Description  get string by ID
// @Tags         genres
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Genre ID"
// @Security 	 ApiKeyAuth
// @Router       /genres/{id} [get]
func Detail(c *fiber.Ctx) error {
	var collection = models.GenreTable()

	id := c.Params("id")

	findResult, msg, code := utils.FindById(collection, id, bson.M{})

	if code != 200 {
		return utils.FailResponse(c, code, msg)
	}

	return utils.SuccessResponse(c, fiber.Map{
		"success": true,
		"message": "Success get data",
		"data":    findResult,
	}, false)
}
