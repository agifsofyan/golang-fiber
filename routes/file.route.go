package routes

import (
	"golang-fiber/controllers"

	"github.com/gofiber/fiber/v2"
)

func FileRoute(route fiber.Router) {
	routes := route.Group("files")
	// routes.Get("/:id", controllers.Detail)
	routes.Post("/", controllers.FileEncode)
	routes.Post("/decode", controllers.FileDecode)
}
