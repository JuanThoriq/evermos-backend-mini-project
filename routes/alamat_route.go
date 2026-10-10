package routes

import (
	"evermos-mini/controllers"
	"evermos-mini/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupAlamatRoutes(router fiber.Router) {
	alamat := router.Group("/user/alamat", middleware.JWTProtected)

	alamat.Get("/", controllers.GetAllAlamat)
	alamat.Get("/:id", controllers.GetAlamatByID)
	alamat.Post("/", controllers.CreateAlamat)
	alamat.Put("/:id", controllers.UpdateAlamat)
	alamat.Delete("/:id", controllers.DeleteAlamat)
}
