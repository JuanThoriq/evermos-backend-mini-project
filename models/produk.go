package models

import "time"

type Produk struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	NamaProduk    string    `json:"nama_produk"`
	Slug          string    `json:"slug"`
	HargaReseller float64   `json:"harga_reseller"`
	HargaKonsumen float64   `json:"harga_konsumen"`
	Stok          int       `json:"stok"`
	Deskripsi     string    `json:"deskripsi"`
	Gambar        string    `json:"gambar"`
	TokoID        uint      `json:"toko_id"`
	CategoryID    uint      `json:"category_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	Toko     Toko     `json:"toko"`
	Category Category `json:"category"`
}

func (Produk) TableName() string { return "produk" }
