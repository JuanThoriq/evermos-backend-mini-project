package routes

import (
	"evermos-mini/controllers"
	"evermos-mini/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(app *fiber.App) {
	user := app.Group("/user", middleware.JWTProtected)

	// Profil
	user.Get("/", controllers.GetMyProfile)
	user.Put("/", controllers.UpdateMyProfile)

	// Alamat
	user.Get("/alamat", controllers.GetAllAlamat)
	user.Get("/alamat/:id", controllers.GetAlamatByID)
	user.Post("/alamat", controllers.CreateAlamat)
	user.Put("/alamat/:id", controllers.UpdateAlamat)
	user.Delete("/alamat/:id", controllers.DeleteAlamat)
}
