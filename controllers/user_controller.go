package controllers

import (
	"evermos-mini/config"
	"evermos-mini/models"
	"time"

	"github.com/gofiber/fiber/v2"
)

// GET /user
func GetMyProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	var user models.User
	if err := config.DB.Preload("Toko").Preload("Alamat").First(&user, userID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  false,
			"message": "User tidak ditemukan",
		})
	}

	return c.JSON(fiber.Map{
		"status": true,
		"data":   user,
	})
}

// PUT /user
func UpdateMyProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	var input struct {
		Nama         string `json:"nama"`
		NoTelp       string `json:"no_telp"`
		TanggalLahir string `json:"tanggal_Lahir"`
		JenisKelamin string `json:"jenis_kelamin"`
		Pekerjaan    string `json:"pekerjaan"`
		Email        string `json:"email"`
		IDProvinsi   string `json:"id_provinsi"`
		IDKota       string `json:"id_kota"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  false,
			"message": "Input tidak valid",
		})
	}

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  false,
			"message": "User tidak ditemukan",
		})
	}

	if input.Nama != "" {
		user.NamaUser = input.Nama
	}
	if input.NoTelp != "" {
		user.NoTelp = input.NoTelp
	}
	if input.JenisKelamin != "" {
		user.JenisKelamin = input.JenisKelamin
	}
	if input.Pekerjaan != "" {
		user.Pekerjaan = input.Pekerjaan
	}
	if input.Email != "" {
		user.Email = input.Email
	}
	if input.IDProvinsi != "" {
		user.IDProvinsi = input.IDProvinsi
	}
	if input.IDKota != "" {
		user.IDKota = input.IDKota
	}

	// Parse tanggal lahir dari format Postman Rakamin
	if input.TanggalLahir != "" {
		parsed, err := time.Parse("02/01/2006", input.TanggalLahir)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"status": false, "message": "Format tanggal lahir harus DD/MM/YYYY"})
		}
		user.TanggalLahir = parsed
	}

	if err := config.DB.Save(&user).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{"status": false, "message": "Email atau nomor telepon sudah digunakan"})
	}

	return c.JSON(fiber.Map{
		"status":  true,
		"message": "Profil berhasil diperbarui",
		"data":    user,
	})
}
