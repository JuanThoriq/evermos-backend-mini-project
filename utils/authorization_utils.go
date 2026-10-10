package utils

import "github.com/gofiber/fiber/v2"

func AuthorizeOwner(c *fiber.Ctx, ownerID uint) bool {
	userID := c.Locals("user_id").(uint)
	return userID == ownerID
}
