package main

import (
	"log"
	"os"

	"example/gorest/config"
	"example/gorest/controllers"
	"example/gorest/routes"
	"example/gorest/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/joho/godotenv"
)

func setupRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		c.Accepts("application/json") // "application/json"

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"message": "This is root endpoint",
		})
	})

	api := app.Group("/api").Group("/v1")

	routes.MovieRoute(api)
	routes.GenreRoute(api)
	routes.UserRoute(api)
}

func main() {
	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error when loading .env")
		}
	}

	app := fiber.New()

	app.Use(cors.New())
	app.Use(logger.New())

	api := app.Group("/api").Group("/v1")
	routes.AuthRoute(api) // Route without authorization

	// JWT Middleware
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
		AuthScheme: "Bearer",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusUnauthorized

			if e, ok := err.(*fiber.Error); ok {
				// Override status code if fiber.Error type
				code = e.Code
			}

			return utils.FailResponse(c, code, "Unauthorized", err)
		},
	}))

	// cfg := fiber.Config{}
	// cfg.ErrorHandler = func(c *Ctx, err error) error {
	// 	code := StatusInternalServerError
	// 	if e, ok := err.(*Error); ok {
	// 		code = e.Code
	// 	}
	// 	c.Set(HeaderContentType, MIMETextPlainCharsetUTF8)
	// 	return c.Status(code).SendString(err.Error())
	// }
	// app := fiber.New(cfg)

	setupRoutes(app) // Route wit authorization

	app.Get("/api/v1/auth/restricted", controllers.Restricted)

	config.ConnectDB()

	port := os.Getenv("PORT")
	err := app.Listen(":" + port)

	if err != nil {
		log.Fatal("Error app failed to start")
		panic(err)
	}
}
