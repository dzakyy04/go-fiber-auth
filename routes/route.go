package routes

import (
	"go-fiber-auth/controllers"
	"go-fiber-auth/middleware"

	"github.com/gofiber/fiber/v2"
)

func RouteInit(route *fiber.App) {
	route.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	// Auth
	route.Post("/api/register", controllers.Register)
	route.Post("/api/login", controllers.Login)
	route.Post("/api/email-verification", controllers.SendVerificationEmail)
	route.Post("/api/verify-email", controllers.VerifyEmail)
	route.Get("/api/me", middleware.AuthMiddleware, middleware.VerifiedMiddleware, controllers.GetMyData)
}
