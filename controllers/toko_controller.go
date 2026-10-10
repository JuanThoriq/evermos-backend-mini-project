package controllers

import (
	"evermos-mini/config"
	"evermos-mini/models"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// ==========================
// GET /toko?limit=10&page=1&nama=xxx
// ==========================
func GetAllToko(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	nama := c.Query("nama", "")

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	var toko []models.Toko
	query := config.DB.Preload("User")

	if nama != "" {
		query = query.Where("LOWER(nama_toko) LIKE ?", "%"+strings.ToLower(nama)+"%")
	}

	var total int64
	query.Model(&models.Toko{}).Count(&total)

	if err := query.Offset(offset).Limit(limit).Find(&toko).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  false,
			"message": "Gagal mengambil data toko",
			"error":   err.Error(),
		})
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return c.JSON(fiber.Map{
		"status": true,
		"pagination": fiber.Map{
			"page":        page,
			"limit":       limit,
			"total_data":  total,
			"total_pages": totalPages,
		},
		"data": toko,
	})
}

// ==========================
// GET /toko/:id_toko
// ==========================
func GetTokoByID(c *fiber.Ctx) error {
	id := c.Params("id_toko")

	var toko models.Toko
	if err := config.DB.Preload("User").First(&toko, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  false,
			"message": "Toko tidak ditemukan",
		})
	}

	return c.JSON(fiber.Map{"status": true, "data": toko})
}

// ==========================
// GET /toko/my
// ==========================
func GetMyToko(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	var toko models.Toko
	if err := config.DB.Where("user_id = ?", userID).Preload("User").First(&toko).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  false,
			"message": "Toko tidak ditemukan",
		})
	}

	return c.JSON(fiber.Map{"status": true, "data": toko})
}

// ==========================
// PUT /toko/:id_toko (form-data)
// ==========================
func UpdateToko(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	id := c.Params("id_toko")

	var toko models.Toko
	result := config.DB.First(&toko, id)
	if result.Error != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  false,
			"message": "Toko tidak ditemukan",
		})
	}

	if toko.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  false,
			"message": "Tidak punya akses ke toko ini",
		})
	}

	// Ambil form-data
	namaToko := c.FormValue("nama_toko")

	file, err := c.FormFile("photo")
	if err == nil {
		// Upload file baru
		filename := fmt.Sprintf("%d-%s", time.Now().Unix(), file.Filename)
		savePath := fmt.Sprintf("./uploads/%s", filename)
		if err := c.SaveFile(file, savePath); err != nil {
			return c.Status(500).JSON(fiber.Map{"status": false, "message": "Gagal menyimpan foto toko"})
		}
		toko.Photo = fmt.Sprintf("http://localhost:8080/uploads/%s", filename)
	}

	if namaToko != "" {
		toko.NamaToko = namaToko
	}

	config.DB.Save(&toko)

	return c.JSON(fiber.Map{
		"status":  true,
		"message": "Toko berhasil diperbarui",
		"data":    toko,
	})
}
