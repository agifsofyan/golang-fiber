package utils

import (
	"context"
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

	sortVal := -1
	if sorts[1] == "desc" {
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
