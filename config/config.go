package config

import (
	"products-api-with-jwt/models"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// SetupDatabase initializes the SQLite database
func SetupDatabase() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		return db, err
	}

	// Migrate tables for User and Product models
	db.AutoMigrate(&models.User{}, &models.Product{}, &models.LoggingHistory{})

	// Populate initial data
	populateInitialData(db)

	return db, nil
}
func populateInitialData(db *gorm.DB) {
	// Check if initial data exists
	var count int64
	db.Model(&models.User{}).Count(&count)
	if count == 0 {
		// Hash password for example users
		passwordHash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

		// Add example users
		users := []models.User{
			{Username: "admin", Password: string(passwordHash), Role: "admin", Department: "IT", Active: false},
			{Username: "user1", Password: string(passwordHash), Role: "user", Department: "Sales", Active: false},
			{Username: "user2", Password: string(passwordHash), Role: "user", Department: "Marketing", Active: false},
		}
		db.Create(&users)
	}

	db.Model(&models.Product{}).Count(&count)
	if count == 0 {
		// Add example products
		products := []models.Product{
			{NamaProduk: "Produk A", Deskripsi: "Deskripsi Produk A", Harga: 1000, Stok: 10},
			{NamaProduk: "Produk B", Deskripsi: "Deskripsi Produk B", Harga: 2000, Stok: 15},
			{NamaProduk: "Produk C", Deskripsi: "Deskripsi Produk C", Harga: 3000, Stok: 20},
		}
		db.Create(&products)
	}

	db.Model(&models.LoggingHistory{}).Count(&count)
	if count == 0 {
		// Example data: Populate with a couple of sample records
		loggingHistoryEntries := []models.LoggingHistory{
			{UserID: 1, JWT: "aaa", ExpiredDate: time.Now().Add(time.Hour * 24), CreatedDate: time.Now().Add(time.Hour * 24)},
		}
		db.Create(&loggingHistoryEntries)
	}
}
