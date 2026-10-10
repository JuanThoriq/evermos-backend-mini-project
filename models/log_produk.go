package models

import "time"

type LogProduk struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	ProdukID      uint      `json:"produk_id"` // foreign key ke produk
	NamaProduk    string    `json:"nama_produk"`
	Slug          string    `json:"slug"`
	HargaReseller float64   `json:"harga_reseller"`
	HargaKonsumen float64   `json:"harga_konsumen"`
	Deskripsi     string    `json:"deskripsi"`
	TokoID        uint      `json:"toko_id"`
	CategoryID    uint      `json:"category_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	// relasi antar tabel
	Produk   Produk   `gorm:"foreignKey:ProdukID" json:"produk,omitempty"`
	Toko     Toko     `gorm:"foreignKey:TokoID" json:"toko,omitempty"`
	Category Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
}
