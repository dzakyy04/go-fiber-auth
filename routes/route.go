package routes

import "github.com/gofiber/fiber/v2"

func RouteInit(route *fiber.App) {
	route.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
}
