package main

import (
	"log"
	"products-api-with-jwt/config"
	"products-api-with-jwt/controllers"
	_ "products-api-with-jwt/docs" // Import docs for Swagger
	"products-api-with-jwt/middlewares"
	"products-api-with-jwt/services"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	// Setup database (SQLite)
	db, err := config.SetupDatabase()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Database connected successfully")

	// Initialize DB for services
	authService := services.NewAuthService(db)
	productService := services.NewProductService(db)

	// Initialize controllers
	authController := controllers.NewAuthController(authService)
	productController := controllers.NewProductController(productService)

	// Initialize router
	r := gin.Default()
	r.Use(middlewares.LoggingMiddleware())
	r.Use(middlewares.RateLimiterMiddleware())

	// Endpoint login (does not require JWT authentication)
	auth := r.Group("/auth")
	auth.POST("/login", authController.Login)
	auth.POST("/logout", authController.Logout)

	// Other endpoints require JWT authentication
	protected := r.Group("/")
	protected.Use(middlewares.JWTAuthMiddleware(authService))

	// Product endpoints
	product := protected.Group("/products")
	product.GET("/", productController.GetProducts)         // Get all products
	product.GET("/:id", productController.GetProductByID)   // Get product by ID
	product.POST("/", productController.CreateProduct)      // Add new product
	product.PUT("/:id", productController.UpdateProduct)    // Update product
	product.DELETE("/:id", productController.DeleteProduct) // Delete product

	// Swagger endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Run server on port 8080
	r.Run(":8080")
}
