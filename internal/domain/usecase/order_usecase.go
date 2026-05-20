package usecase

import (
	"goapp/internal/domain/entity"
	"goapp/internal/domain/repository"
)

// OrderUseCase handles order operations
type OrderUseCase struct {
	orderRepo repository.OrderRepository
}

// NewOrderUseCase creates a new OrderUseCase
func NewOrderUseCase(orderRepo repository.OrderRepository) *OrderUseCase {
	return &OrderUseCase{orderRepo: orderRepo}
}

// GetUserOrders retrieves orders for a user
// VULNERABLE: Passes unsanitized userID to repository (SQLi)
func (uc *OrderUseCase) GetUserOrders(userID string) ([]entity.Order, error) {
	return uc.orderRepo.FindByUserID(userID)
}

// GetUserOrdersByStatus retrieves orders by user and status
// VULNERABLE: Passes unsanitized parameters to repository (SQLi)
func (uc *OrderUseCase) GetUserOrdersByStatus(userID, status string) ([]entity.Order, error) {
	return uc.orderRepo.FindByUserIDAndStatus(userID, status)
}

// GetOrder retrieves an order by ID
// VULNERABLE: No authorization check (IDOR)
func (uc *OrderUseCase) GetOrder(orderID string) (*entity.Order, error) {
	return uc.orderRepo.FindByID(orderID)
}

// CreateOrder creates a new order
func (uc *OrderUseCase) CreateOrder(order *entity.Order) error {
	return uc.orderRepo.Create(order)
}

// GetOrderHistory retrieves order history with product details
// VULNERABLE: SQL Injection in JOIN query via username
func (uc *OrderUseCase) GetOrderHistory(username string) ([]map[string]interface{}, error) {
	return uc.orderRepo.FindOrdersWithProductDetails(username)
}
