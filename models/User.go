package models

import (
	"golang-fiber/config"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name,omitempty" bson:"name,omitempty"`
	Email     string             `json:"email,omitempty" bson:"email,omitempty"`
	Pass      string             `json:"pass,omitempty" bson:"pass,omitempty"`
	CreatedAt time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	IsAdmin   bool               `json:"is_admin,omitempty" bson:"is_admin,omitempty"`
}

type SwaggerLogin struct {
	Email string `json:"email,omitempty" bson:"email,omitempty"`
	Pass  string `json:"pass,omitempty" bson:"pass,omitempty"`
}

type SwaggerRegis struct {
	Name  string `json:"name,omitempty" bson:"name,omitempty"`
	Email string `json:"email,omitempty" bson:"email,omitempty"`
	Pass  string `json:"pass,omitempty" bson:"pass,omitempty"`
}

func UserTable() *mongo.Collection {
	return config.MI.DB.Collection("users")
}
