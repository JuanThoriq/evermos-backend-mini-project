package models

import "time"

type DetailTransaksi struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	LogProdukID uint      `json:"log_produk_id"`
	TransaksiID uint      `json:"transaksi_id"`
	Kuantitas   int       `json:"kuantitas"`
	HargaTotal  float64   `json:"harga_total"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	LogProduk LogProduk  `json:"log_produk" gorm:"foreignKey:LogProdukID"`
	Transaksi *Transaksi `json:"-" gorm:"foreignKey:TransaksiID"`
}
