package repository

import "goapp/internal/domain/entity"

// OrderRepository defines the interface for order data operations
// VULNERABLE: Implementations intentionally use unsafe patterns (SQLi, IDOR)
type OrderRepository interface {
	// FindByUserID finds orders by user ID (SQLi vulnerable)
	FindByUserID(userID string) ([]entity.Order, error)

	// FindByUserIDAndStatus finds orders by user ID and status (SQLi vulnerable)
	FindByUserIDAndStatus(userID, status string) ([]entity.Order, error)

	// FindByID finds order by ID (IDOR - no ownership check)
	FindByID(orderID string) (*entity.Order, error)

	// FindOrdersWithProductDetails finds orders with JOIN (SQLi in JOIN)
	FindOrdersWithProductDetails(username string) ([]map[string]interface{}, error)

	// Create creates a new order
	Create(order *entity.Order) error
}
