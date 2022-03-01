package models

import (
	"golang-fiber/gorest/config"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Genre struct {
	ID   primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name string             `json:"name,omitempty" bson:"name,omitempty"`
	Slug string             `json:"slug,omitempty" bson:"slug,omitempty"`
}

type SwaggerInsertGenre struct {
	Name string `json:"name,omitempty" bson:"name,omitempty"`
}

func GenreTable() *mongo.Collection {
	return config.MI.DB.Collection("genres")
}
