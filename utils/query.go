package utils

import (
	"context"
	"fmt"
	"log"
	"math"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Result struct {
	Total int64
	Page  int
	Last  float64
	Limit int64
}

type Dynamic struct {
	Value interface{}
}

func Search(text string, opt []string) bson.M {
	field := bson.M{
		"$regex": primitive.Regex{
			Pattern: text,
			Options: "i",
		},
	}

	elements := make([]bson.M, 0)

	for _, v := range opt {
		element := bson.M{v: field}
		elements = append(elements, element)
	}

	filter := bson.M{
		"$or": elements,
	}

	return filter
}

func Paginate(c *fiber.Ctx, collect *mongo.Collection, filter interface{}, sorts [2]string) (*mongo.Cursor, Result, context.Context, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limitVal, _ := strconv.Atoi(c.Query("limit", "10"))
	var limit int64 = int64(limitVal)
	total, _ := collect.CountDocuments(ctx, filter)

	last := math.Ceil(float64(total/limit)) + 1

	if (last < 1 && total > 0) || (limit == total) {
		last = 1
	}

	if limit == 1 {
		last -= 1
	}

	sortVal := -1
	if sorts[1] == "asc" {
		sortVal = 1
	}

	opt := []bson.M{
		{
			"$match": filter,
		},
		{
			"$lookup": bson.M{
				"from":         "genres",
				"localField":   "genre",
				"foreignField": "_id",
				"as":           "genre",
			},
		},
		{
			"$skip": (int64(page) - 1) * limit,
		},
		{
			"$limit": limit,
		},
		{
			"$sort": bson.M{sorts[0]: sortVal},
		},
	}

	cursor, err := collect.Aggregate(ctx, opt)

	if err != nil {
		return nil, Result{}, ctx, err
	}

	result := Result{
		Total: total,
		Page:  page,
		Last:  last,
		Limit: limit,
	}

	return cursor, result, ctx, nil
}

func Detailed(c *fiber.Ctx, collection *mongo.Collection, match bson.M, ref bson.M) (interface{}, string) {
	var result bson.M

	err := collection.FindOne(context.TODO(), match).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, "error no document"
		}
	}

	if err != nil {
		log.Println("err::", err)
		return nil, "error in aggregate"
	}

	return result, ""
}

func Upload(c *fiber.Ctx, nameFile string) (fiber.Map, string) {
	file, err := c.FormFile(nameFile)

	if err != nil {
		return nil, "fail get the file"
	}

	// Get buffer from file
	buffer, err := file.Open()
	if err != nil {
		return nil, "buffer file failed"
	}
	defer buffer.Close()

	fileName := file.Filename
	contentType := file.Header["Content-Type"][0]
	fileZise := file.Size
	dir := "temp"
	filePath := fmt.Sprintf("./%s/%s", dir, fileName)

	err = c.SaveFile(file, filePath)

	if err != nil {
		return nil, "cannot save the file"
	}

	return fiber.Map{
		"name":         fileName,
		"size":         fileZise,
		"type":         contentType,
		"path":         dir,
		"relativePath": filePath,
	}, ""
}
