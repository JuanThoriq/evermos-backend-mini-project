package routes

import (
	"evermos-mini/controllers"
	"github.com/gofiber/fiber/v2"
)

func SetupProvCityRoutes(app *fiber.App) {
	app.Get("/provcity/listprovincies", controllers.GetListProvincies)
	app.Get("/provcity/listcities/:prov_id", controllers.GetListCities)
	app.Get("/provcity/detailprovince/:prov_id", controllers.GetDetailProvince)
	app.Get("/provcity/detailcity/:city_id", controllers.GetDetailCity)
}
