package persistence

import (
	"database/sql"
	"fmt"

	"goapp/internal/domain/entity"
)

// SQLiteProductRepository implements ProductRepository with SQLite
type SQLiteProductRepository struct {
	db *sql.DB
}

// NewSQLiteProductRepository creates a new SQLiteProductRepository
func NewSQLiteProductRepository(db *sql.DB) *SQLiteProductRepository {
	return &SQLiteProductRepository{db: db}
}

// Search searches products by name and category
// VULNERABLE: SQL Injection via string formatting
func (r *SQLiteProductRepository) Search(query, category string) ([]entity.Product, error) {
	// VULNERABLE: Direct string concatenation
	sqlQuery := fmt.Sprintf(
		"SELECT id, name, price, stock, category FROM products WHERE name LIKE '%%%s%%'",
		query,
	)

	// VULNERABLE: Additional SQL injection point
	if category != "" {
		sqlQuery += fmt.Sprintf(" AND category = '%s'", category)
	}

	rows, err := r.db.Query(sqlQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []entity.Product
	for rows.Next() {
		var product entity.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.Stock, &product.Category); err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

// SearchWithSort searches products with sorting
// VULNERABLE: ORDER BY injection
func (r *SQLiteProductRepository) SearchWithSort(query, category, sortBy, order string) ([]entity.Product, error) {
	// VULNERABLE: Direct string concatenation
	sqlQuery := fmt.Sprintf(
		"SELECT id, name, price, stock, category FROM products WHERE name LIKE '%%%s%%'",
		query,
	)

	if category != "" {
		sqlQuery += fmt.Sprintf(" AND category = '%s'", category)
	}

	// VULNERABLE: ORDER BY injection - sortBy and order directly concatenated
	if sortBy != "" {
		sqlQuery += fmt.Sprintf(" ORDER BY %s", sortBy)
		if order != "" {
			sqlQuery += fmt.Sprintf(" %s", order)
		}
	}

	rows, err := r.db.Query(sqlQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []entity.Product
	for rows.Next() {
		var product entity.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.Stock, &product.Category); err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

// FindByID finds product by ID
func (r *SQLiteProductRepository) FindByID(id int64) (*entity.Product, error) {
	var product entity.Product
	err := r.db.QueryRow(
		"SELECT id, name, price, stock, category FROM products WHERE id = ?",
		id,
	).Scan(&product.ID, &product.Name, &product.Price, &product.Stock, &product.Category)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// GetAll returns all products
func (r *SQLiteProductRepository) GetAll() ([]entity.Product, error) {
	rows, err := r.db.Query("SELECT id, name, price, stock, category FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []entity.Product
	for rows.Next() {
		var product entity.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.Stock, &product.Category); err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

// Create creates a new product
func (r *SQLiteProductRepository) Create(product *entity.Product) error {
	result, err := r.db.Exec(
		"INSERT INTO products (name, price, stock, category) VALUES (?, ?, ?, ?)",
		product.Name, product.Price, product.Stock, product.Category,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	product.ID = id

	return nil
}
