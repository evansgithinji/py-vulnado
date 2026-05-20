package entity

// Product represents a product in the store
type Product struct {
	ID       int64   `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Stock    int     `json:"stock"`
	Category string  `json:"category"`
}
