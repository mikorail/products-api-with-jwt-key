package services

import (
	"errors"
	"os"
	"strconv"
	"time"

	"products-api-with-jwt/models"

	global "products-api-with-jwt/global"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hashicorp/go-memdb"
	"gorm.io/gorm"
)

type AuthService struct {
	DB        *gorm.DB
	MemDB     *memdb.MemDB // Add MemDB instance to AuthService
	secretKey string
}

// NewAuthService menginisialisasi AuthService baru
func NewAuthService(db *gorm.DB, MemDB *memdb.MemDB) *AuthService {
	return &AuthService{
		DB:        db,
		MemDB:     MemDB,                          // Add MemDB instance to AuthServic
		secretKey: os.Getenv(global.ENVSecretKey), // Ganti dengan secret key yang aman
	}
}

// ValidateCredentials memvalidasi username dan password
func (s *AuthService) ValidateCredentials(username, password string) (models.User, error) {
	var user models.User
	// Cari pengguna berdasarkan username
	if err := s.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return user, errors.New("invalid username or password")
	}

	// Verifikasi password (gunakan bcrypt atau metode hashing yang sama yang digunakan saat menyimpan)
	if err := models.CheckPasswordHash(password, user.Password); err != nil {
		return user, errors.New("invalid username or password")
	}

	return user, nil
}

// GenerateToken menghasilkan token JWT untuk pengguna
func (s *AuthService) GenerateToken(username string, expiration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(expiration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secretKey))
}

// StoreTokenInMemDB stores the JWT token for the user in MemDB
func (as *AuthService) StoreTokenInMemDB(userID int, username string, token string) error {
	txn := as.MemDB.Txn(true)
	defer txn.Abort()

	// Create a new token entry with the user ID
	newToken := &models.Token{
		ID:       strconv.Itoa(userID), // Ensure this is a unique ID
		UserID:   strconv.Itoa(userID), // Store the user ID as a string
		Username: username,             // Store the username
		JWT:      token,                // Store the JWT
	}

	// Check if there is an existing token for the user
	if existingToken, err := txn.First("token_schemas", "user_id", newToken.UserID); err == nil && existingToken != nil {
		// Remove existing token
		if err := txn.Delete("token_schemas", existingToken); err != nil {
			return err
		}
	}

	// Insert the new token
	if err := txn.Insert("token_schemas", newToken); err != nil {
		return err
	}

	txn.Commit() // Commit the transaction
	return nil
}
