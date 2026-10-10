package config

import (
	"evermos-mini/models"
	"fmt"
)

func SeedData() {
	var count int64
	DB.Model(&models.Category{}).Count(&count)
	if count == 0 {
		categories := []models.Category{
			{NamaCategory: "Fashion"},
			{NamaCategory: "Aksesoris"},
			{NamaCategory: "Rumah Tangga"},
		}
		if err := DB.Create(&categories).Error; err != nil {
			fmt.Println("Gagal membuat kategori default:", err)
		} else {
			fmt.Println("Kategori default berhasil dibuat")
		}
	}
}
