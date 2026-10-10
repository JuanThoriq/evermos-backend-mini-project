package controllers

import (
	"evermos-mini/config"
	"evermos-mini/models"
	"github.com/gofiber/fiber/v2"
)

// GET /category
func GetAllKategori(c *fiber.Ctx) error {
	var kategori []models.Category
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	search := c.Query("search", "")
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit
	query := config.DB.Model(&models.Category{})

	if search != "" {
		query = query.Where("nama_category LIKE ?", "%"+search+"%")
	}

	var total int64
	query.Count(&total)
	query.Offset(offset).Limit(limit).Find(&kategori)

	return c.JSON(fiber.Map{
		"status": true,
		"data":   kategori,
		"meta": fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// GET /category/:id
func GetKategoriByID(c *fiber.Ctx) error {
	id := c.Params("id")

	var category models.Category
	if err := config.DB.First(&category, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"status": false, "message": "Kategori tidak ditemukan"})
	}

	return c.JSON(fiber.Map{"status": true, "data": category})
}

// POST /category (Admin only)
func CreateKategori(c *fiber.Ctx) error {
	role := c.Locals("role")
	if role != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  false,
			"message": "Hanya admin yang dapat menambah kategori",
		})
	}

	var input struct {
		NamaCategory string `json:"nama_category"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"status": false, "message": "Input tidak valid"})
	}

	if input.NamaCategory == "" {
		return c.Status(400).JSON(fiber.Map{"status": false, "message": "Nama kategori wajib diisi"})
	}

	category := models.Category{NamaCategory: input.NamaCategory}
	config.DB.Create(&category)

	return c.JSON(fiber.Map{"status": true, "message": "Kategori berhasil ditambahkan", "data": category})
}

// PUT /category/:id (Admin only)
func UpdateKategori(c *fiber.Ctx) error {
	role := c.Locals("role")
	if role != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  false,
			"message": "Hanya admin yang dapat mengubah kategori",
		})
	}

	id := c.Params("id")

	var category models.Category
	if err := config.DB.First(&category, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"status": false, "message": "Kategori tidak ditemukan"})
	}

	var input struct {
		NamaCategory string `json:"nama_category"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"status": false, "message": "Input tidak valid"})
	}

	if input.NamaCategory == "" {
		return c.Status(400).JSON(fiber.Map{"status": false, "message": "Nama kategori wajib diisi"})
	}

	category.NamaCategory = input.NamaCategory
	config.DB.Save(&category)

	return c.JSON(fiber.Map{"status": true, "message": "Kategori berhasil diperbarui", "data": category})
}

// DELETE /category/:id (Admin only)
func DeleteKategori(c *fiber.Ctx) error {
	role := c.Locals("role")
	if role != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  false,
			"message": "Hanya admin yang dapat menghapus kategori",
		})
	}

	id := c.Params("id")

	var category models.Category
	if err := config.DB.First(&category, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"status": false, "message": "Kategori tidak ditemukan"})
	}

	config.DB.Delete(&category)
	return c.JSON(fiber.Map{"status": true, "message": "Kategori berhasil dihapus"})
}
