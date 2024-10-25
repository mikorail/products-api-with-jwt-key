package models

type User struct {
	ID         int    `gorm:"primaryKey"`
	Username   string `gorm:"unique;not null"`
	Password   string `gorm:"not null"`
	Role       string `gorm:"not null"`
	Department string `gorm:"not null"`
	Active     bool
}
