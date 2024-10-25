package services

import (
	"errors"
	"products-api-with-jwt/models"

	"gorm.io/gorm"
)

type ProductService struct {
	DB *gorm.DB
}

func NewProductService(db *gorm.DB) *ProductService {
	return &ProductService{DB: db}
}

// GetAllProducts mengambil semua produk dari database
func (s *ProductService) GetAllProducts() ([]models.Product, error) {
	var products []models.Product
	if err := s.DB.Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

// GetProductByID mengambil produk berdasarkan ID
func (s *ProductService) GetProductByID(id int) (*models.Product, error) {
	var product models.Product
	if err := s.DB.First(&product, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}
	return &product, nil
}

// CreateProduct menambah produk baru ke database
func (s *ProductService) CreateProduct(product *models.Product) (models.Product, error) {
	// Menyimpan produk baru ke database
	if err := s.DB.Create(product).Error; err != nil {
		return models.Product{}, err // Kembalikan error jika terjadi kesalahan
	}
	return *product, nil // Kembalikan produk yang baru dibuat
}

// UpdateProduct memperbarui hanya field yang disediakan dalam permintaan
func (s *ProductService) UpdateProduct(id int, updatedProduct *models.Product) (*models.Product, error) {
	var product models.Product

	// Cari produk berdasarkan ID
	if err := s.DB.First(&product, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	// Perbarui hanya field yang disediakan
	if updatedProduct.NamaProduk != "" {
		product.NamaProduk = updatedProduct.NamaProduk
	}
	if updatedProduct.Deskripsi != "" {
		product.Deskripsi = updatedProduct.Deskripsi
	}
	if updatedProduct.Harga != 0 {
		product.Harga = updatedProduct.Harga
	}
	if updatedProduct.Stok != 0 {
		product.Stok = updatedProduct.Stok
	}

	// Simpan perubahan ke database
	if err := s.DB.Save(&product).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

// DeleteProduct menghapus produk berdasarkan ID
func (s *ProductService) DeleteProduct(id int) error {
	if err := s.DB.Delete(&models.Product{}, id).Error; err != nil {
		return err
	}
	return nil
}
