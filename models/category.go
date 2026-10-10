package models

import "time"

type Category struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	NamaCategory string    `json:"nama_category"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	Produk []Produk `json:"produk"`
}
