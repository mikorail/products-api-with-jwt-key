package controllers

import (
	"net/http"
	"strconv"

	"products-api-with-jwt/models"
	"products-api-with-jwt/services"

	"github.com/gin-gonic/gin"
)

type ProductController struct {
	ProductService *services.ProductService
}

// NewProductController menginisialisasi ProductController baru
func NewProductController(productService *services.ProductService) *ProductController {
	return &ProductController{ProductService: productService}
}

// GetProducts godoc
// @Summary Get all products
// @Description Get a list of all products
// @Tags products
// @Produce json
// @Success 200 {object} models.ApiResponse
// @Failure 500 {object} models.ApiResponse
// @Router /products [get]
func (pc *ProductController) GetProducts(c *gin.Context) {
	products, err := pc.ProductService.GetAllProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ApiResponse{
			Status:  "error",
			Code:    http.StatusInternalServerError,
			Message: "Could not retrieve products",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, models.ApiResponse{
		Status:  "success",
		Code:    http.StatusOK,
		Message: "Products retrieved successfully",
		Data:    products,
		Count:   len(products), // Optional count of items
	})
}

// GetProductByID godoc
// @Summary Get product by ID
// @Description Get details of a product by its ID
// @Tags products
// @Param id path int true "Product ID"
// @Produce json
// @Success 200 {object} models.ApiResponse
// @Failure 400 {object} models.ApiResponse
// @Failure 404 {object} models.ApiResponse
// @Router /products/{id} [get]
func (pc *ProductController) GetProductByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ApiResponse{
			Status:  "error",
			Code:    http.StatusBadRequest,
			Message: "Invalid product ID",
			Data:    nil,
		})
		return
	}

	product, err := pc.ProductService.GetProductByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ApiResponse{
			Status:  "error",
			Code:    http.StatusNotFound,
			Message: "Product not found",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, models.ApiResponse{
		Status:  "success",
		Code:    http.StatusOK,
		Message: "Product retrieved successfully",
		Data:    product,
	})
}

// CreateProduct godoc
// @Summary Create a new product
// @Description Create a new product with the given details
// @Tags products
// @Accept json
// @Produce json
// @Param product body models.Product true "Product"
// @Success 201 {object} models.ApiResponse
// @Failure 400 {object} models.ApiResponse
// @Failure 500 {object} models.ApiResponse
// @Router /products [post]
func (pc *ProductController) CreateProduct(c *gin.Context) {
	var input models.Product
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, models.ApiResponse{
			Status:  "error",
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	product, err := pc.ProductService.CreateProduct(&input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ApiResponse{
			Status:  "error",
			Code:    http.StatusInternalServerError,
			Message: "Could not create product",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusCreated, models.ApiResponse{
		Status:  "success",
		Code:    http.StatusCreated,
		Message: "Product created successfully",
		Data:    product,
	})
}

// UpdateProduct godoc
// @Summary Update a product by ID
// @Description Update a product's information by its ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param product body models.Product true "Product"
// @Success 200 {object} models.ApiResponse
// @Failure 400 {object} models.ApiResponse
// @Failure 404 {object} models.ApiResponse
// @Failure 500 {object} models.ApiResponse
// @Router /products/{id} [put]
func (pc *ProductController) UpdateProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ApiResponse{
			Status:  "error",
			Code:    http.StatusBadRequest,
			Message: "Invalid product ID",
			Data:    nil,
		})
		return
	}

	var input models.Product
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, models.ApiResponse{
			Status:  "error",
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	updatedProduct, err := pc.ProductService.UpdateProduct(id, &input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ApiResponse{
			Status:  "error",
			Code:    http.StatusInternalServerError,
			Message: "Could not update product",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, models.ApiResponse{
		Status:  "success",
		Code:    http.StatusOK,
		Message: "Product updated successfully",
		Data:    updatedProduct,
	})
}

// DeleteProduct godoc
// @Summary Delete a product by ID
// @Description Delete a product by its ID
// @Tags products
// @Param id path int true "Product ID"
// @Produce json
// @Success 200 {object} models.ApiResponse
// @Failure 400 {object} models.ApiResponse
// @Failure 500 {object} models.ApiResponse
// @Router /products/{id} [delete]
func (pc *ProductController) DeleteProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ApiResponse{
			Status:  "error",
			Code:    http.StatusBadRequest,
			Message: "Invalid product ID",
			Data:    nil,
		})
		return
	}

	if err := pc.ProductService.DeleteProduct(id); err != nil {
		c.JSON(http.StatusInternalServerError, models.ApiResponse{
			Status:  "error",
			Code:    http.StatusInternalServerError,
			Message: "Could not delete product",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, models.ApiResponse{
		Status:  "success",
		Code:    http.StatusOK,
		Message: "Product deleted successfully",
		Data:    nil,
	})
}
