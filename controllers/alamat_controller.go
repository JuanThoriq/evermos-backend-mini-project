package controllers

import (
	"evermos-mini/config"
	"evermos-mini/models"
	"time"

	"github.com/gofiber/fiber/v2"
)

// =================== GET MY ALAMAT ===================
func GetAllAlamat(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	var alamat []models.Alamat
	config.DB.Where("user_id = ?", userID).Find(&alamat)
	return c.JSON(fiber.Map{"status": true, "data": alamat})
}

// =================== GET ALAMAT BY ID ===================
func GetAlamatByID(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	id := c.Params("id")

	var alamat models.Alamat
	if err := config.DB.First(&alamat, "id = ? AND user_id = ?", id, userID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"status": false, "message": "Alamat tidak ditemukan"})
	}

	return c.JSON(fiber.Map{"status": true, "data": alamat})
}

// =================== CREATE ALAMAT ===================
func CreateAlamat(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	var input struct {
		JudulAlamat  string `json:"judul_alamat"`
		NamaPenerima string `json:"nama_penerima"`
		NoTelp       string `json:"no_telp"`
		Provinsi     string `json:"provinsi"`
		Kota         string `json:"kota"`
		Kecamatan    string `json:"kecamatan"`
		Kelurahan    string `json:"kelurahan"`
		DetailAlamat string `json:"detail_alamat"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": false, "message": "Input tidak valid"})
	}
	if input.JudulAlamat == "" || input.NamaPenerima == "" || input.NoTelp == "" || input.DetailAlamat == "" {
		return c.Status(400).JSON(fiber.Map{"status": false, "message": "Data alamat wajib diisi"})
	}

	alamat := models.Alamat{
		JudulAlamat:  input.JudulAlamat,
		NamaPenerima: input.NamaPenerima,
		NoTelp:       input.NoTelp,
		Provinsi:     input.Provinsi,
		Kota:         input.Kota,
		Kecamatan:    input.Kecamatan,
		Kelurahan:    input.Kelurahan,
		DetailAlamat: input.DetailAlamat,
		UserID:       userID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := config.DB.Create(&alamat).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"status": false, "message": "Gagal menambahkan alamat"})
	}
	return c.JSON(fiber.Map{"status": true, "message": "Alamat berhasil ditambahkan", "data": alamat})
}

// =================== UPDATE ALAMAT ===================
func UpdateAlamat(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	id := c.Params("id")

	var alamat models.Alamat
	if err := config.DB.First(&alamat, "id = ? AND user_id = ?", id, userID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"status": false, "message": "Alamat tidak ditemukan"})
	}

	var input struct {
		JudulAlamat  string `json:"judul_alamat"`
		NamaPenerima string `json:"nama_penerima"`
		NoTelp       string `json:"no_telp"`
		Provinsi     string `json:"provinsi"`
		Kota         string `json:"kota"`
		Kecamatan    string `json:"kecamatan"`
		Kelurahan    string `json:"kelurahan"`
		DetailAlamat string `json:"detail_alamat"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": false, "message": "Input tidak valid"})
	}

	alamat.JudulAlamat = input.JudulAlamat
	alamat.NamaPenerima = input.NamaPenerima
	alamat.NoTelp = input.NoTelp
	alamat.Provinsi = input.Provinsi
	alamat.Kota = input.Kota
	alamat.Kecamatan = input.Kecamatan
	alamat.Kelurahan = input.Kelurahan
	alamat.DetailAlamat = input.DetailAlamat
	alamat.UpdatedAt = time.Now()

	if err := config.DB.Save(&alamat).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"status": false, "message": "Gagal memperbarui alamat"})
	}
	return c.JSON(fiber.Map{"status": true, "message": "Alamat berhasil diperbarui", "data": alamat})
}

// =================== DELETE ALAMAT ===================
func DeleteAlamat(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	id := c.Params("id")

	var alamat models.Alamat
	if err := config.DB.First(&alamat, "id = ? AND user_id = ?", id, userID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"status": false, "message": "Alamat tidak ditemukan"})
	}

	if err := config.DB.Delete(&alamat).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"status": false, "message": "Gagal menghapus alamat"})
	}
	return c.JSON(fiber.Map{"status": true, "message": "Alamat berhasil dihapus"})
}
