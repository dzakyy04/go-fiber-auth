package middleware

import (
	"go-fiber-auth/database"
	"go-fiber-auth/models/entity"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte(os.Getenv("JWT_KEY"))

func AuthMiddleware(ctx *fiber.Ctx) error {
	// Check token from Authorization header
	authHeader := ctx.Get("Authorization")
	if authHeader == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized",
			"error":   "Authorization header is missing",
		})
	}

	// Split header to get token
	tokenParts := strings.Split(authHeader, "Bearer ")
	if len(tokenParts) != 2 {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized",
			"error":   "Invalid authorization header format",
		})
	}

	tokenString := tokenParts[1]
	if tokenString == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized",
			"error":   "Token is empty",
		})
	}

	// Parse and validate token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.NewError(fiber.StatusUnauthorized, "Unexpected signing method")
		}
		return jwtKey, nil
	})
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized",
			"error":   err.Error(),
		})
	}

	// Get claims from token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized",
			"error":   "Invalid token claims",
		})
	}

	// Get user_id from claims
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized",
			"error":   "Invalid user id in token claims",
		})
	}

	// Find user from database
	var user entity.User
	err = database.DB.First(&user, userID).Error
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "User not found",
			"error":   err.Error(),
		})
	}

	// Save user data into context
	ctx.Locals("user", &user)
	return ctx.Next()
}
