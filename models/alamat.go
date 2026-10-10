package models

import "time"

type Alamat struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	JudulAlamat  string    `json:"judul_alamat"`
	NamaPenerima string    `json:"nama_penerima"`
	NoTelp       string    `json:"no_telp"`
	Provinsi     string    `json:"provinsi"`
	Kota         string    `json:"kota"`
	Kecamatan    string    `json:"kecamatan"`
	Kelurahan    string    `json:"kelurahan"`
	DetailAlamat string    `json:"detail_alamat"`
	UserID       uint      `json:"user_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
