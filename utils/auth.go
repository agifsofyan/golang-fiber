package utils

import (
	"example/gorest/models"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserInfo struct {
	UUID  primitive.ObjectID `json:"uuid,omitempty" bson:"uuid,omitempty"`
	Name  string             `json:"name,omitempty" bson:"name,omitempty"`
	Email string             `json:"email,omitempty" bson:"email,omitempty"`
}

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type CustomClaims struct {
	jwt.StandardClaims
	UserInfo
}

var expired time.Duration = time.Hour * 72

func GenerateToken(user models.User) (map[string]string, error) {
	claims := CustomClaims{
		jwt.StandardClaims{
			Issuer:    os.Getenv("APPLICATION_NAME"),
			ExpiresAt: time.Now().Add(expired).Unix(),
		},
		UserInfo{
			user.ID,
			user.Name,
			user.Email,
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"access_token": signedToken,
	}, nil
}

func ExtractToken(c *fiber.Ctx, keyword string) string {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	return claims[keyword].(string)
}
