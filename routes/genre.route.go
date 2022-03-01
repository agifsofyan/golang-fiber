package routes

import (
	"golang-fiber/controllers"

	"github.com/gofiber/fiber/v2"
)

func GenreRoute(route fiber.Router) {
	routes := route.Group("genres")
	routes.Get("/", controllers.Index)
	routes.Get("/:id", controllers.Detail)
	routes.Post("/", controllers.Add)
}
