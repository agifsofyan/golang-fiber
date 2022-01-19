package models

import (
	"example/gorest/config"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type InsertMovie struct {
	ID          primitive.ObjectID   `json:"_id,omitempty" bson:"_id,omitempty"`
	Title       string               `json:"title,omitempty" bson:"title,omitempty"`
	Slug        string               `json:"slug,omitempty" bson:"slug,omitempty"`
	Genre       []primitive.ObjectID `json:"genre,omitempty" bson:"genre,omitempty"`
	Description string               `json:"description,omitempty" bson:"description,omitempty"`
	CreatedAt   time.Time            `json:"created_at,omitempty" bson:"created_at,omitempty"`
}

type Movie struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title       string             `json:"title,omitempty" bson:"title,omitempty"`
	Slug        string             `json:"slug,omitempty" bson:"slug,omitempty"`
	Genre       []Genre            `json:"genre,omitempty" bson:"genre,omitempty"`
	Description string             `json:"description,omitempty" bson:"description,omitempty"`
	CreatedAt   time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
}

type SwaggerInsertMovie struct {
	Title       string               `json:"title,omitempty" bson:"title,omitempty"`
	Genre       []primitive.ObjectID `json:"genre,omitempty" bson:"genre,omitempty"`
	Description string               `json:"description,omitempty" bson:"description,omitempty"`
}

func MovieTable() *mongo.Collection {
	return config.MI.DB.Collection("movies")
}
