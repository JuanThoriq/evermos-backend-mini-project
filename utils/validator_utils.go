package utils

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

func ValidateStruct(c *fiber.Ctx, input interface{}) bool {
	if err := validate.Struct(input); err != nil {
		var errors []string
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, err.Field()+" tidak valid atau kosong")
		}
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Validasi gagal",
			"errors":  errors,
		})
		return false
	}
	return true
}
