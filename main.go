package main

import (
	"log"
	"os"

	"golang-fiber/gorest/config"
	"golang-fiber/gorest/controllers"
	"golang-fiber/gorest/routes"
	"golang-fiber/gorest/utils"

	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/joho/godotenv"

	_ "golang-fiber/gorest/docs"
)

// @title           Swagger MOvies CRUD API
// @version         2.0
// @description     Rest APIs golang - fiber.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey  ApiKeyAuth
// @in                          header
// @name                        Authorization
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

	freeRoute(app) // without Auth

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

			return utils.FailResponse(c, code, "Unauthorized")
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

	restrictRoute(app) // Route wit authorization

	app.Get("/api/v1/auth/restricted", controllers.Restricted)

	config.ConnectDB()

	port := os.Getenv("PORT")
	err := app.Listen(":" + port)

	if err != nil {
		log.Fatal("Error app failed to start")
		panic(err)
	}
}

func freeRoute(app *fiber.App) {
	app.Get("/swagger/*", swagger.HandlerDefault) // default

	api := app.Group("/api").Group("/v1")

	routes.AuthRoute(api) // Route without authorization
	routes.FileRoute(api)
}

func restrictRoute(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		c.Accepts("application/json") // "application/json"

		return utils.SuccessResponse(c, fiber.Map{
			"success": true,
			"message": "This is root endpoint",
		}, false)
	})

	api := app.Group("/api").Group("/v1")

	routes.MovieRoute(api)
	routes.GenreRoute(api)
	routes.UserRoute(api)
}
