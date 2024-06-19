package controllers

import (
	"go-fiber-auth/database"
	"go-fiber-auth/models/entity"
	"go-fiber-auth/models/request"
	"go-fiber-auth/utils"

	"github.com/gofiber/fiber/v2"
)

func Register(ctx *fiber.Ctx) error {
	// Parse request body
	req := new(request.RegisterRequest)
	if err := ctx.BodyParser(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	// Create new user
	user := entity.User{
		Name:  req.Name,
		Email: req.Email,
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to hash password",
			"error":   err.Error(),
		})
	}
	user.Password = hashedPassword

	// Save user to database
	if err := database.DB.Create(&user).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to create user",
			"error":   err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Successfully registered",
	})
}

