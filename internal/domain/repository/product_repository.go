package repository

import "goapp/internal/domain/entity"

// ProductRepository defines the interface for product data operations
// VULNERABLE: Implementations intentionally use unsafe SQL patterns
type ProductRepository interface {
	// Search searches products by name and category (SQLi vulnerable)
	Search(query, category string) ([]entity.Product, error)

	// SearchWithSort searches products with sorting (ORDER BY injection)
	SearchWithSort(query, category, sortBy, order string) ([]entity.Product, error)

	// FindByID finds product by ID
	FindByID(id int64) (*entity.Product, error)

	// GetAll returns all products
	GetAll() ([]entity.Product, error)

	// Create creates a new product
	Create(product *entity.Product) error
}
