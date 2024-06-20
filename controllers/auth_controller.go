package controllers

import (
	"go-fiber-auth/database"
	"go-fiber-auth/models/entity"
	"go-fiber-auth/models/request"
	"go-fiber-auth/utils"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
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

func Login(ctx *fiber.Ctx) error {
	// Parse request body
	req := new(request.LoginRequest)
	if err := ctx.BodyParser(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	// Check if user exists
	var user entity.User
	err := database.DB.Where("email = ?", req.Email).First(&user).Error
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Invalid email or password",
			"error":   "User not found",
		})
	}

	// Check password
	if !utils.VerifyPassword(req.Password, user.Password) {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Invalid email or password",
		})
	}

	// Generate JWT token
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
		"iat":     time.Now().Unix(),
	}

	token, err := utils.GenerateToken(&claims)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to generate token",
			"error":   err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Successfully logged in",
		"data": fiber.Map{
			"user":  user,
			"token": token,
		},
	})
}

func SendVerificationEmail(ctx *fiber.Ctx) error {
	// Parse request body
	req := new(request.EmailVerificationRequest)
	if err := ctx.BodyParser(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	// Check if user exists
	var user entity.User
	err := database.DB.Where("email = ?", req.Email).First(&user).Error
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Invalid email or password",
			"error":   "User not found",
		})
	}

	// Generate OTP
	otp := utils.GenerateOTP()
	otpExpiresAt := time.Now().Add(time.Minute * 10)

	// Save OTP to database
	err = database.DB.Model(&user).Updates(map[string]interface{}{
		"otp":            otp,
		"otp_expires_at": otpExpiresAt,
	}).Error
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to send OTP",
			"error":   err.Error(),
		})
	}

	// Send email
	subject := "Email Verification"
	template := "views/emails/verification.html"
	err = utils.SendEmail(user.Email, subject, template, fiber.Map{
		"Name": user.Name,
		"Otp":  otp,
	})
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to send email",
			"error":   err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "OTP code has been sent to email",
	})
}

func VerifyEmail(ctx *fiber.Ctx) error {
	// Parse request body
	req := new(request.VerifyEmailRequest)
	if err := ctx.BodyParser(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	// Check if user exists
	var user entity.User
	err := database.DB.Where("email = ?", req.Email).First(&user).Error
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Invalid email or password",
			"error":   "User not found",
		})
	}

	// Check if OTP is valid
	if user.OTP == nil || *user.OTP != req.OTP || time.Now().After(*user.OTPExpiresAt) {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Invalid OTP code",
		})
	}

	// Update user's email verification status
	err = database.DB.Model(&user).Updates(map[string]interface{}{
		"email_verified_at": time.Now(),
		"otp":               nil,
		"otp_expires_at":    nil,
	}).Error
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to verify email",
			"error":   err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Email has been verified",
	})
}

func GetMyData(ctx *fiber.Ctx) error {
	// Get user_id claims from context
	userID := ctx.Locals("user_id")
	log.Println(userID)

	if userID == nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized",
			"error":   "User ID not found in context",
		})
	}

	// Find user from database
	var user entity.User
	err := database.DB.First(&user, userID).Error
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "User not found",
			"error":   err.Error(),
		})
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"message": "Successfully retrieved user data",
		"data": fiber.Map{
			"user": user,
		},
	})
}
