package routes

import (
	"evermos-mini/controllers"

	"github.com/gofiber/fiber/v2"
)

func SetupAuthRoutes(app *fiber.App) {
	app.Post("/auth/login", controllers.Login)
	app.Post("/auth/register", controllers.Register)
}
