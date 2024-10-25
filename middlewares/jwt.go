package middlewares

import (
	"fmt"
	"net/http"
	"os"
	"products-api-with-jwt/global"
	"products-api-with-jwt/models"
	"products-api-with-jwt/services"

	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte(os.Getenv(global.ENVSecretKey))

// GenerateToken membuat JWT token dengan waktu kedaluwarsa dinamis (1 hari atau 7 hari)
func GenerateToken(username string, expiration time.Duration) (string, error) {
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
		Issuer:    "your-app",
		Subject:   username,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

// JWTAuthMiddleware memvalidasi token JWT di header Authorization setiap request
func JWTAuthMiddleware(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, models.ApiResponse{
				Status:  "error",
				Code:    http.StatusBadRequest,
				Message: "Authorization header required"})
			c.Abort()
			return
		}

		// Ambil token dari header Authorization
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, models.ApiResponse{
				Status:  "error",
				Code:    http.StatusBadRequest,
				Message: "Bearer token required"})
			c.Abort()
			return
		}

		// Validasi token JWT
		token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
			// Verifikasi metode signing
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, models.ApiResponse{
				Status:  "error",
				Code:    http.StatusBadRequest,
				Message: "Invalid or expired token"})
			c.Abort()
			return
		}

		// Extract claims untuk mengambil informasi user jika diperlukan
		claims, ok := token.Claims.(*jwt.RegisteredClaims)
		if !ok || !token.Valid {
			c.JSON(http.StatusUnauthorized, models.ApiResponse{
				Status:  "error",
				Code:    http.StatusBadRequest,
				Message: "Invalid token claims"})
			c.Abort()
			return
		}

		// Check if user is active
		idCheck, err := authService.GetUserIDFromToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, models.ApiResponse{
				Status:  "error",
				Code:    http.StatusBadRequest,
				Message: "Userid not found"})
			c.Abort()
			return
		}
		user, error := authService.GetUserById(int(idCheck))
		if error != nil {
			c.JSON(http.StatusUnauthorized, models.ApiResponse{
				Status:  "error",
				Code:    http.StatusBadRequest,
				Message: "User not found"})
			c.Abort()
			return
		}

		// Check if user is active
		if !user.Active {
			c.JSON(http.StatusUnauthorized, models.ApiResponse{
				Status:  "error",
				Code:    http.StatusBadRequest,
				Message: "You are not logged in"})
			c.Abort()
			return
		}

		// Menyimpan username atau informasi lain dari token ke context
		c.Set("user_id", claims.Subject)
		c.Next()
	}
}
