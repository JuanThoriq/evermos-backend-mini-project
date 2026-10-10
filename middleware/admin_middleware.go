package middleware

import (
	"evermos-mini/config"
	"evermos-mini/models"

	"github.com/gofiber/fiber/v2"
)

func AdminOnly(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  false,
			"message": "User tidak ditemukan",
		})
	}

	if user.Role != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  false,
			"message": "Hanya admin yang dapat mengakses fitur ini",
		})
	}

	return c.Next()
}
