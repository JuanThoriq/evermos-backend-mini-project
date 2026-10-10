package controllers

import (
	"errors"
	"evermos-mini/config"
	"evermos-mini/models"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// POST /trx
func CreateTrx(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	var input struct {
		MetodeBayar string `json:"method_bayar"`
		AlamatKirim uint   `json:"alamat_kirim"`
		DetailTrx   []struct {
			ProdukID  uint `json:"product_id"`
			Kuantitas int  `json:"kuantitas"`
		} `json:"detail_trx"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"status": false, "message": "Input tidak valid"})
	}
	if input.MetodeBayar == "" || input.AlamatKirim == 0 || len(input.DetailTrx) == 0 {
		return c.Status(400).JSON(fiber.Map{"status": false, "message": "Metode bayar, alamat, dan produk wajib diisi"})
	}

	var transaksi models.Transaksi
	err := config.DB.Transaction(func(tx *gorm.DB) error {
		var alamat models.Alamat
		if err := tx.Where("id = ? AND user_id = ?", input.AlamatKirim, userID).First(&alamat).Error; err != nil {
			return errors.New("alamat tidak ditemukan")
		}

		transaksi = models.Transaksi{
			UserID: userID, AlamatPengiriman: input.AlamatKirim,
			MetodeBayar: input.MetodeBayar,
			Invoice:     fmt.Sprintf("INV-%d-%d", userID, time.Now().UnixNano()),
			Status:      "Menunggu Pembayaran",
		}
		if err := tx.Create(&transaksi).Error; err != nil {
			return err
		}

		for _, item := range input.DetailTrx {
			if item.ProdukID == 0 || item.Kuantitas < 1 {
				return errors.New("produk atau kuantitas tidak valid")
			}
			var produk models.Produk
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&produk, item.ProdukID).Error; err != nil {
				return errors.New("produk tidak ditemukan")
			}
			if produk.Stok < item.Kuantitas {
				return fmt.Errorf("stok %s tidak mencukupi", produk.NamaProduk)
			}

			logProduk := models.LogProduk{
				ProdukID: produk.ID, NamaProduk: produk.NamaProduk, Slug: produk.Slug,
				HargaReseller: produk.HargaReseller, HargaKonsumen: produk.HargaKonsumen,
				Deskripsi: produk.Deskripsi, TokoID: produk.TokoID, CategoryID: produk.CategoryID,
			}
			if err := tx.Create(&logProduk).Error; err != nil {
				return err
			}
			subtotal := produk.HargaKonsumen * float64(item.Kuantitas)
			if err := tx.Create(&models.DetailTransaksi{
				TransaksiID: transaksi.ID, LogProdukID: logProduk.ID,
				Kuantitas: item.Kuantitas, HargaTotal: subtotal,
			}).Error; err != nil {
				return err
			}
			if err := tx.Model(&produk).Update("stok", gorm.Expr("stok - ?", item.Kuantitas)).Error; err != nil {
				return err
			}
			transaksi.HargaTotal += subtotal
		}
		return tx.Model(&transaksi).Update("harga_total", transaksi.HargaTotal).Error
	})
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": false, "message": err.Error()})
	}

	config.DB.Preload("DetailTransaksi.LogProduk").First(&transaksi, transaksi.ID)
	return c.Status(201).JSON(fiber.Map{"status": true, "message": "Transaksi berhasil dibuat", "data": transaksi})
}

// GET /trx
func GetAllTrx(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	var transaksi []models.Transaksi
	if err := config.DB.Where("user_id = ?", userID).Preload("DetailTransaksi.LogProduk").Find(&transaksi).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"status": false, "message": "Gagal mengambil transaksi"})
	}
	return c.JSON(fiber.Map{"status": true, "data": transaksi})
}

// GET /trx/:id
func GetTrxByID(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	var transaksi models.Transaksi
	if err := config.DB.Preload("DetailTransaksi.LogProduk").Where("id = ? AND user_id = ?", c.Params("id"), userID).First(&transaksi).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"status": false, "message": "Transaksi tidak ditemukan"})
	}
	return c.JSON(fiber.Map{"status": true, "data": transaksi})
}
