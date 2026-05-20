package entity

// Review represents a product review
type Review struct {
	ID        int64  `json:"id"`
	ProductID int64  `json:"product_id"`
	UserID    int64  `json:"user_id"`
	Rating    int    `json:"rating"`
	Comment   string `json:"comment"`
}
