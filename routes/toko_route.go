package routes

import (
	"evermos-mini/controllers"
	"evermos-mini/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupTokoRoutes(app *fiber.App) {
	app.Get("/toko", controllers.GetAllToko)
	toko := app.Group("/toko", middleware.JWTProtected)
	toko.Get("/my", controllers.GetMyToko)
	toko.Put("/:id_toko", controllers.UpdateToko)

	app.Get("/toko/:id_toko", controllers.GetTokoByID)
}
