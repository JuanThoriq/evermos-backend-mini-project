package controllers

import (
	"evermos-mini/config"
	"evermos-mini/models"
	"evermos-mini/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// POST /auth/register
func Register(c *fiber.Ctx) error {
	var input struct {
		NamaUser     string `json:"nama" validate:"required"`
		Email        string `json:"email" validate:"required,email"`
		NoTelp       string `json:"no_telp" validate:"required"`
		KataSandi    string `json:"kata_sandi" validate:"required,min=6"`
		TanggalLahir string `json:"tanggal_Lahir"`
		Pekerjaan    string `json:"pekerjaan"`
		IDProvinsi   string `json:"id_provinsi"`
		IDKota       string `json:"id_kota"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  false,
			"message": "Input tidak valid",
		})
	}
	if !utils.ValidateStruct(c, input) {
		return nil
	}

	hash, err := utils.HashPassword(input.KataSandi)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": false, "message": "Gagal memproses kata sandi"})
	}

	var tgl time.Time
	if input.TanggalLahir != "" {
		tgl, err = time.Parse("02/01/2006", input.TanggalLahir)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"status": false, "message": "Format tanggal lahir harus DD/MM/YYYY"})
		}
	}

	user := models.User{
		NamaUser:     input.NamaUser,
		Email:        input.Email,
		NoTelp:       input.NoTelp,
		KataSandi:    hash,
		TanggalLahir: tgl,
		Pekerjaan:    input.Pekerjaan,
		IDProvinsi:   input.IDProvinsi,
		IDKota:       input.IDKota,
		Role:         "user",
	}

	err = config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&user).Error; err != nil {
			return err
		}
		return tx.Create(&models.Toko{NamaToko: "Toko " + input.NamaUser, UserID: user.ID}).Error
	})
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  false,
			"message": "Email atau nomor telepon sudah digunakan",
		})
	}

	return c.JSON(fiber.Map{
		"status":  true,
		"message": "Registrasi berhasil",
		"data":    user,
	})
}

// POST /auth/login
func Login(c *fiber.Ctx) error {
	var input struct {
		NoTelp    string `json:"no_telp" validate:"required"`
		KataSandi string `json:"kata_sandi" validate:"required"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"status": false, "message": "Input tidak valid"})
	}

	var user models.User
	if err := config.DB.Where("no_telp = ?", input.NoTelp).First(&user).Error; err != nil {
		return c.Status(401).JSON(fiber.Map{"status": false, "message": "Nomor telepon tidak ditemukan"})
	}

	if !utils.CheckPasswordHash(input.KataSandi, user.KataSandi) {
		return c.Status(401).JSON(fiber.Map{"status": false, "message": "Password salah"})
	}

	token, err := utils.GenerateJWT(user.ID, user.Role)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": false, "message": "Gagal membuat token"})
	}

	return c.JSON(fiber.Map{
		"status":  true,
		"message": "Login berhasil",
		"data": fiber.Map{
			"token": token,
			"user":  user,
		},
	})
}

// GET /user
func GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	var user models.User
	if err := config.DB.Preload("Toko").Preload("Alamat").Preload("Transaksi").
		First(&user, userID).Error; err != nil {
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
func UpdateProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	var input struct {
		Nama         string `json:"nama"`
		Email        string `json:"email"`
		NoTelp       string `json:"no_telp"`
		TanggalLahir string `json:"tanggal_Lahir"`
		JenisKelamin string `json:"jenis_kelamin"`
		Pekerjaan    string `json:"pekerjaan"`
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

	// isi semua field sesuai JSON body
	if input.Nama != "" {
		user.NamaUser = input.Nama
	}
	if input.Email != "" {
		user.Email = input.Email
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
	if input.IDProvinsi != "" {
		user.IDProvinsi = input.IDProvinsi
	}
	if input.IDKota != "" {
		user.IDKota = input.IDKota
	}
	if input.TanggalLahir != "" {
		tgl, err := time.Parse("02/01/2006", input.TanggalLahir)
		if err == nil {
			user.TanggalLahir = tgl
		}
	}

	config.DB.Save(&user)

	return c.JSON(fiber.Map{
		"status":  true,
		"message": "Profil berhasil diperbarui",
		"data":    user,
	})
}
