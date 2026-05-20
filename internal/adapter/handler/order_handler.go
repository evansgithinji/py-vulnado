package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"goapp/internal/domain/usecase"
)

// OrderHandler handles order-related HTTP requests
type OrderHandler struct {
	orderUseCase *usecase.OrderUseCase
}

// NewOrderHandler creates a new OrderHandler
func NewOrderHandler(orderUseCase *usecase.OrderUseCase) *OrderHandler {
	return &OrderHandler{orderUseCase: orderUseCase}
}

// RegisterRoutes registers order routes
func (h *OrderHandler) RegisterRoutes(r *gin.Engine) {
	r.GET("/api/orders", h.GetOrders)
	r.GET("/api/orders/:id", h.GetOrder)
	r.GET("/api/orders/history", h.GetOrderHistory)

	// HTML pages
	r.GET("/orders", h.OrdersPage)
	r.GET("/orders/history", h.OrderHistoryPage)
}

// GetOrders retrieves orders for a user
// VULNERABLE: SQL Injection via user_id and status parameters
func (h *OrderHandler) GetOrders(c *gin.Context) {
	// VULNERABLE: user_id passed as string for SQLi demonstration
	userID := c.Query("user_id")
	status := c.Query("status")

	var orders interface{}
	var err error

	if status != "" {
		// VULNERABLE: Both parameters injectable
		orders, err = h.orderUseCase.GetUserOrdersByStatus(userID, status)
	} else {
		// VULNERABLE: user_id injectable
		orders, err = h.orderUseCase.GetUserOrders(userID)
	}

	if err != nil {
		// VULNERABLE: Information disclosure
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "Failed to fetch orders",
			"detail": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// GetOrder retrieves an order by ID
// VULNERABLE: IDOR - No authorization check
func (h *OrderHandler) GetOrder(c *gin.Context) {
	orderID := c.Param("id")

	// VULNERABLE: IDOR - Returns any order without ownership check
	// VULNERABLE: Potential SQLi via orderID
	order, err := h.orderUseCase.GetOrder(orderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

// GetOrderHistory retrieves order history with product details
// VULNERABLE: SQL Injection in JOIN query
func (h *OrderHandler) GetOrderHistory(c *gin.Context) {
	username := c.Query("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username required"})
		return
	}

	// VULNERABLE: SQLi in JOIN query via username
	orders, err := h.orderUseCase.GetOrderHistory(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "Failed to fetch order history",
			"detail": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// OrdersPage serves the orders page
// VULNERABLE: Reflected XSS
func (h *OrderHandler) OrdersPage(c *gin.Context) {
	userID := c.Query("user_id")
	status := c.Query("status")

	var orders interface{}
	if userID != "" {
		if status != "" {
			orders, _ = h.orderUseCase.GetUserOrdersByStatus(userID, status)
		} else {
			orders, _ = h.orderUseCase.GetUserOrders(userID)
		}
	}

	// VULNERABLE: Reflected XSS via userID and status
	html := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head><title>Orders - Vulnerable App</title></head>
<body>
	<h1>Orders</h1>
	<form method="GET" action="/orders">
		<label>User ID: <input type="text" name="user_id" value="%s"></label><br><br>
		<label>Status: <input type="text" name="status" value="%s"></label><br><br>
		<button type="submit">View Orders</button>
	</form>
	<h2>Orders for User: %s (Status: %s)</h2>
`, userID, status, userID, status) // VULNERABLE: Reflected XSS

	if orderList, ok := orders.([]interface{}); ok && len(orderList) > 0 {
		html += `<table border="1">
			<tr><th>ID</th><th>User ID</th><th>Product ID</th><th>Total</th><th>Status</th></tr>`
		// Display orders
		html += `</table>`
	}

	html += `
	<p>Hints:</p>
	<ul>
		<li>IDOR: Try /api/orders/1, /api/orders/2, etc.</li>
		<li>SQLi: Try user_id=1 OR 1=1</li>
	</ul>
</body>
</html>
`
	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, html)
}

// OrderHistoryPage serves the order history page with product details
// VULNERABLE: SQL Injection in JOIN query + Reflected XSS
func (h *OrderHandler) OrderHistoryPage(c *gin.Context) {
	username := c.Query("username")

	var orders []map[string]interface{}
	if username != "" {
		orders, _ = h.orderUseCase.GetOrderHistory(username)
	}

	// VULNERABLE: Reflected XSS via username
	html := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head><title>Order History - Vulnerable App</title></head>
<body>
	<h1>Order History (JOIN-based SQL Injection Demo)</h1>
	<form method="GET" action="/orders/history">
		<label>Username: <input type="text" name="username" value="%s"></label>
		<button type="submit">View History</button>
	</form>
	<h2>Order History for: %s</h2>
`, username, username)

	if len(orders) > 0 {
		html += `<table border="1">
			<tr><th>Product</th><th>Price</th><th>Quantity</th><th>Total</th><th>Status</th></tr>`
		for _, o := range orders {
			html += fmt.Sprintf("<tr><td>%v</td><td>$%.2f</td><td>%v</td><td>$%.2f</td><td>%v</td></tr>",
				o["product_name"], o["price"], o["quantity"], o["total"], o["status"])
		}
		html += `</table>`
	}

	html += `
	<p>Hints:</p>
	<ul>
		<li>JOIN SQLi: Try username=' OR '1'='1</li>
		<li>UNION: Try username=' UNION SELECT ...</li>
		<li>This query JOINs orders, products, and users tables</li>
	</ul>
</body>
</html>
`
	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, html)
}
