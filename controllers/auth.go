package controllers

import (
	"net/http"
	"time"

	"products-api-with-jwt/models"
	"products-api-with-jwt/services"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	AuthService *services.AuthService
}

// NewAuthController menginisialisasi AuthController baru
func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{AuthService: authService}
}

// Login menangani permintaan login dan menghasilkan token JWT
func (ac *AuthController) Login(c *gin.Context) {
	var input models.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, models.ApiResponse{
			Status:  "error",
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	// Validate credentials
	user, err := ac.AuthService.ValidateCredentials(input.Username, input.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.ApiResponse{
			Status:  "error",
			Code:    http.StatusUnauthorized,
			Message: "Invalid username or password",
			Data:    nil,
		})
		return
	}

	expiration := time.Hour * 24 // Default to 1 day
	if input.RememberMe {
		expiration = time.Hour * 24 * 7 // 7 days
	}

	// Generate new token
	token, err := ac.AuthService.GenerateToken(user.Username, expiration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ApiResponse{
			Status:  "error",
			Code:    http.StatusInternalServerError,
			Message: "Could not generate token",
			Data:    nil,
		})
		return
	}

	// Store the new token in MemDB, expiring the previous one if exists
	err = ac.AuthService.StoreTokenInMemDB(user.ID, user.Username, token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ApiResponse{
			Status:  "error",
			Code:    http.StatusInternalServerError,
			Message: "Could not store token in MemDB",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, models.ApiResponse{
		Status:  "success",
		Code:    http.StatusOK,
		Message: "Login successful",
		Data:    gin.H{"token": token},
	})
}
