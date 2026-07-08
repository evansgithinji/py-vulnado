package persistence

import (
	"database/sql"
	"fmt"

	"goapp/internal/domain/entity"
)

// SQLiteReviewRepository implements ReviewRepository with SQLite
type SQLiteReviewRepository struct {
	db *sql.DB
}

// NewSQLiteReviewRepository creates a new SQLiteReviewRepository
func NewSQLiteReviewRepository(db *sql.DB) *SQLiteReviewRepository {
	return &SQLiteReviewRepository{db: db}
}

// FindByProductID finds reviews for a product
func (r *SQLiteReviewRepository) FindByProductID(productID int64) ([]entity.Review, error) {
	rows, err := r.db.Query(
		"SELECT id, product_id, user_id, rating, comment FROM reviews WHERE product_id = ?",
		productID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []entity.Review
	for rows.Next() {
		var review entity.Review
		if err := rows.Scan(&review.ID, &review.ProductID, &review.UserID, &review.Rating, &review.Comment); err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}

	return reviews, nil
}

// Create creates a new review
// VULNERABLE: SQL Injection via comment field + stored XSS
func (r *SQLiteReviewRepository) Create(review *entity.Review) error {
	// VULNERABLE: Direct string concatenation allows SQL injection
	// VULNERABLE: Comment stored without sanitization (stored XSS)
	query := fmt.Sprintf(
		"INSERT INTO reviews (product_id, user_id, rating, comment) VALUES (%d, %d, %d, '%s')",
		review.ProductID, review.UserID, review.Rating, review.Comment,
	)

	result, err := r.db.Exec(query)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	review.ID = id

	return nil
}
