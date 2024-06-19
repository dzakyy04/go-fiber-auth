package routes

import (
	"go-fiber-auth/controllers"

	"github.com/gofiber/fiber/v2"
)

func RouteInit(route *fiber.App) {
	route.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	// Auth
	route.Post("/api/register", controllers.Register)
}
