package routes

import (
	"evermos-mini/controllers"
	"evermos-mini/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupTrxRoutes(app *fiber.App) {
	trx := app.Group("/trx", middleware.JWTProtected)

	trx.Get("/", controllers.GetAllTrx)
	trx.Get("/:id", controllers.GetTrxByID)
	trx.Post("/", controllers.CreateTrx)
}
