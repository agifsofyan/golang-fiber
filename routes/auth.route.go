package routes

import (
	"golang-fiber/controllers"

	"github.com/gofiber/fiber/v2"
)

func AuthRoute(route fiber.Router) {
	routes := route.Group("auth")

	// routes.Get("/", controllers.Accessible)
	routes.Post("/", controllers.Login)
	routes.Post("/register", controllers.Register)
	routes.Get("/email/:email", func(c *fiber.Ctx) error {
		email := c.Params("email")

		return controllers.ByEmail(c, email)
	})
}
