package models

type Product struct {
	ID         int    `gorm:"primaryKey"`
	NamaProduk string `gorm:"not null"`
	Deskripsi  string
	Harga      float64
	Stok       int
}
