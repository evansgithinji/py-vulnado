package persistence

import (
	"database/sql"
	"fmt"

	"goapp/internal/domain/entity"
)

// SQLiteOrderRepository implements OrderRepository with SQLite
type SQLiteOrderRepository struct {
	db *sql.DB
}

// NewSQLiteOrderRepository creates a new SQLiteOrderRepository
func NewSQLiteOrderRepository(db *sql.DB) *SQLiteOrderRepository {
	return &SQLiteOrderRepository{db: db}
}

// FindByUserID finds orders by user ID
// VULNERABLE: SQL Injection via userID parameter (string type allows injection)
func (r *SQLiteOrderRepository) FindByUserID(userID string) ([]entity.Order, error) {
	// VULNERABLE: Direct string concatenation with string userID
	sqlQuery := fmt.Sprintf(
		"SELECT id, user_id, product_id, quantity, total, status FROM orders WHERE user_id = %s",
		userID,
	)

	rows, err := r.db.Query(sqlQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []entity.Order
	for rows.Next() {
		var order entity.Order
		if err := rows.Scan(&order.ID, &order.UserID, &order.ProductID, &order.Quantity, &order.Total, &order.Status); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

// FindByUserIDAndStatus finds orders by user ID and status
// VULNERABLE: Multiple SQL injection points
func (r *SQLiteOrderRepository) FindByUserIDAndStatus(userID, status string) ([]entity.Order, error) {
	// VULNERABLE: Both parameters are injectable
	sqlQuery := fmt.Sprintf(
		"SELECT id, user_id, product_id, quantity, total, status FROM orders WHERE user_id = %s AND status = '%s'",
		userID, status,
	)

	rows, err := r.db.Query(sqlQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []entity.Order
	for rows.Next() {
		var order entity.Order
		if err := rows.Scan(&order.ID, &order.UserID, &order.ProductID, &order.Quantity, &order.Total, &order.Status); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

// FindByID finds order by ID
// VULNERABLE: IDOR - No authorization check, returns any order
func (r *SQLiteOrderRepository) FindByID(orderID string) (*entity.Order, error) {
	// VULNERABLE: IDOR + potential SQL injection
	sqlQuery := fmt.Sprintf(
		"SELECT id, user_id, product_id, quantity, total, status FROM orders WHERE id = %s",
		orderID,
	)

	var order entity.Order
	err := r.db.QueryRow(sqlQuery).Scan(&order.ID, &order.UserID, &order.ProductID, &order.Quantity, &order.Total, &order.Status)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

// Create creates a new order
func (r *SQLiteOrderRepository) Create(order *entity.Order) error {
	result, err := r.db.Exec(
		"INSERT INTO orders (user_id, product_id, quantity, total, status) VALUES (?, ?, ?, ?, ?)",
		order.UserID, order.ProductID, order.Quantity, order.Total, order.Status,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	order.ID = id

	return nil
}

// FindOrdersWithProductDetails finds orders with JOIN to products table
// VULNERABLE: SQL Injection in JOIN query
func (r *SQLiteOrderRepository) FindOrdersWithProductDetails(username string) ([]map[string]interface{}, error) {
	// VULNERABLE: JOIN-based SQL injection
	sqlQuery := fmt.Sprintf(`
		SELECT p.name, p.price, o.quantity, o.total, o.status
		FROM orders o
		JOIN products p ON o.product_id = p.id
		JOIN users u ON o.user_id = u.id
		WHERE u.username = '%s'
	`, username)

	rows, err := r.db.Query(sqlQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var productName string
		var price float64
		var quantity int
		var total float64
		var status string
		if err := rows.Scan(&productName, &price, &quantity, &total, &status); err != nil {
			return nil, err
		}
		results = append(results, map[string]interface{}{
			"product_name": productName,
			"price":        price,
			"quantity":     quantity,
			"total":        total,
			"status":       status,
		})
	}

	return results, nil
}
