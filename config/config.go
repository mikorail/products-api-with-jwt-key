package config

import (
	"products-api-with-jwt/models"

	"github.com/hashicorp/go-memdb"
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
	db.AutoMigrate(&models.User{}, &models.Product{})

	// Populate initial data
	populateInitialData(db)

	return db, nil
}

// SetupMemDB initializes the MemDB instance and creates the necessary schema
func SetupMemDB() (*memdb.MemDB, error) {
	// Define the schema for the tokens table

	var tokensSchema = &memdb.TableSchema{
		Name: "token_schemas",
		Indexes: map[string]*memdb.IndexSchema{
			"id": {
				Name:    "id",
				Unique:  true,
				Indexer: &memdb.StringFieldIndex{Field: "ID"}, // Ensure this matches the ID field in Token
			},
			"user_id": {
				Name:    "user_id",
				Unique:  false,
				Indexer: &memdb.StringFieldIndex{Field: "UserID"}, // Index for UserID
			},
			"username": {
				Name:    "username",
				Unique:  false,
				Indexer: &memdb.StringFieldIndex{Field: "Username"}, // Index for Username
			},
		},
	}
	var Schema = &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"token_schemas": tokensSchema,
		},
	}

	// Create the MemDB instance with the defined schema
	mdb, err := memdb.NewMemDB(Schema)
	if err != nil {
		return nil, err
	}

	return mdb, nil
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
			{Username: "admin", Password: string(passwordHash), Role: "admin", Department: "IT"},
			{Username: "user1", Password: string(passwordHash), Role: "user", Department: "Sales"},
			{Username: "user2", Password: string(passwordHash), Role: "user", Department: "Marketing"},
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
}
