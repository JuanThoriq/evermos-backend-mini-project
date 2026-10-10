package routes

import (
	"evermos-mini/controllers"
	"evermos-mini/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupProductRoutes(app *fiber.App) {
	// Public
	app.Get("/product", controllers.GetAllProduct)
	app.Get("/product/:id", controllers.GetProductByID)

	// Protected Product (Hanya user login)
	product := app.Group("/product", middleware.JWTProtected)
	product.Post("/", controllers.CreateProduct)
	product.Put("/:id", controllers.UpdateProduct)
	product.Delete("/:id", controllers.DeleteProduct)
}
