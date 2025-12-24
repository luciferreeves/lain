package controllers

import (
	"lain/services"
	"lain/session"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func GetEmailAPI(context *fiber.Ctx) error {
	emailID, err := strconv.ParseUint(context.Params("id"), 10, 32)
	if err != nil {
		return context.Status(400).JSON(fiber.Map{"error": "Invalid email ID"})
	}

	userEmail, err := session.GetSessionEmail(context)
	if err != nil {
		return context.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	email, err := services.GetEmailDetails(userEmail, uint(emailID))
	if err != nil {
		return context.Status(404).JSON(fiber.Map{"error": "Email not found"})
	}

	return context.JSON(email)
}

func ToggleFlagAPI(context *fiber.Ctx) error {
	emailID, err := strconv.ParseUint(context.Params("id"), 10, 32)
	if err != nil {
		return context.Status(400).JSON(fiber.Map{"error": "Invalid email ID"})
	}

	userEmail, err := session.GetSessionEmail(context)
	if err != nil {
		return context.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	isFlagged, err := services.ToggleEmailFlag(userEmail, uint(emailID))
	if err != nil {
		return context.Status(500).JSON(fiber.Map{"error": "Failed to toggle flag"})
	}

	return context.JSON(fiber.Map{"flagged": isFlagged})
}

func MarkEmailAsReadAPI(context *fiber.Ctx) error {
	emailID, err := strconv.ParseUint(context.Params("id"), 10, 32)
	if err != nil {
		return context.Status(400).JSON(fiber.Map{"error": "Invalid email ID"})
	}

	userEmail, err := session.GetSessionEmail(context)
	if err != nil {
		return context.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	err = services.MarkEmailAsRead(userEmail, uint(emailID))
	if err != nil {
		return context.Status(500).JSON(fiber.Map{"error": "Failed to mark as read"})
	}

	return context.JSON(fiber.Map{"success": true})
}
