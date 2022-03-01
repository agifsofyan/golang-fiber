package routes

import (
	"golang-fiber/controllers"

	"github.com/gofiber/fiber/v2"
)

func UserRoute(route fiber.Router) {
	routes := route.Group("user")
	routes.Get("/me", controllers.Me)
}
