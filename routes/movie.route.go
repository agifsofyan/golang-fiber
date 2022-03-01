package routes

import (
	"golang-fiber/gorest/controllers"

	"github.com/gofiber/fiber/v2"
)

func MovieRoute(route fiber.Router) {
	routes := route.Group("movies")
	routes.Get("/", controllers.GetMovies)
	routes.Get("/:id", controllers.GetMovie)
	routes.Post("/", controllers.AddMovie)
	routes.Put("/:id", controllers.UpdateMovie)
	routes.Delete("/:id", controllers.DeleteMovie)
}
