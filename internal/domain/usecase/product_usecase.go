package usecase

import (
	"goapp/internal/domain/entity"
	"goapp/internal/domain/repository"
)

// ProductUseCase handles product operations
type ProductUseCase struct {
	productRepo repository.ProductRepository
}

// NewProductUseCase creates a new ProductUseCase
func NewProductUseCase(productRepo repository.ProductRepository) *ProductUseCase {
	return &ProductUseCase{productRepo: productRepo}
}

// SearchProducts searches for products
// VULNERABLE: Passes unsanitized query and category to repository (SQLi)
func (uc *ProductUseCase) SearchProducts(query, category string) ([]entity.Product, error) {
	return uc.productRepo.Search(query, category)
}

// SearchProductsWithSort searches for products with sorting
// VULNERABLE: ORDER BY injection via sortBy and order parameters
func (uc *ProductUseCase) SearchProductsWithSort(query, category, sortBy, order string) ([]entity.Product, error) {
	return uc.productRepo.SearchWithSort(query, category, sortBy, order)
}

// GetProduct retrieves a product by ID
func (uc *ProductUseCase) GetProduct(productID int64) (*entity.Product, error) {
	return uc.productRepo.FindByID(productID)
}

// GetAllProducts retrieves all products
func (uc *ProductUseCase) GetAllProducts() ([]entity.Product, error) {
	return uc.productRepo.GetAll()
}
