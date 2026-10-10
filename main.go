package main

import (
	"evermos-mini/config"
	"evermos-mini/models"
	"evermos-mini/routes"
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// === Inisialisasi Database ===
	config.InitDB()

	// === Auto Migrate semua tabel ===
	if err := config.DB.AutoMigrate(
		&models.User{},
		&models.Toko{},
		&models.Alamat{},
		&models.Category{},
		&models.Produk{},
		&models.LogProduk{},
		&models.Transaksi{},
		&models.DetailTransaksi{},
	); err != nil {
		panic(err)
	}

	// === Inisiasi seed database ===
	config.SeedData()

	// === Setup Fiber app ===
	app := fiber.New()

	// Ini agar folder uploads bisa diakses dari browser
	app.Static("/uploads", "./uploads")

	// === Middleware Keamanan ===

	// 🔹 CORS agar API bisa diakses dari frontend
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// 🔹 Rate Limiter untuk mencegah spam request
	app.Use(limiter.New(limiter.Config{
		Max:        60,              // maksimal 60 request
		Expiration: 1 * time.Minute, // dalam 1 menit
	}))

	// 🔹 Logger untuk mencatat semua request ke terminal
	app.Use(logger.New())

	// === Setup Routes ===
	routes.SetupRoutes(app)

	// === Jalankan server di port dari .env ===
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = ":8080"
	}

	fmt.Println("Server running on http://localhost" + port)
	if err := app.Listen(port); err != nil {
		panic(err)
	}
}
