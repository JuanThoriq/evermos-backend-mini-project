package utils

import (
	"evermos-mini/config"
	"evermos-mini/models"
)

func CreateLogProduk(produk models.Produk) {
	log := models.LogProduk{
		ProdukID:      produk.ID,
		NamaProduk:    produk.NamaProduk,
		Slug:          produk.Slug,
		HargaReseller: produk.HargaReseller,
		HargaKonsumen: produk.HargaKonsumen,
		Deskripsi:     produk.Deskripsi,
		TokoID:        produk.TokoID,
		CategoryID:    produk.CategoryID,
	}

	config.DB.Create(&log)
}
