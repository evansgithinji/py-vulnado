package repository

import "goapp/internal/domain/entity"

// ReviewRepository defines the interface for review data operations
// VULNERABLE: Implementations intentionally use unsafe patterns
type ReviewRepository interface {
	// FindByProductID finds reviews for a product
	FindByProductID(productID int64) ([]entity.Review, error)

	// Create creates a new review (SQLi vulnerable, stored XSS)
	Create(review *entity.Review) error
}
