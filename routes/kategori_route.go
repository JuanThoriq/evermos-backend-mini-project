package routes

import (
	"evermos-mini/controllers"
	"evermos-mini/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupKategoriRoutes(app *fiber.App) {
	app.Get("/category", controllers.GetAllKategori)
	app.Get("/category/:id", controllers.GetKategoriByID)

	// admin only
	kategori := app.Group("/category", middleware.JWTProtected, middleware.AdminOnly)
	kategori.Post("/", controllers.CreateKategori)
	kategori.Put("/:id", controllers.UpdateKategori)
	kategori.Delete("/:id", controllers.DeleteKategori)
}
