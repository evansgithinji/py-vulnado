package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"goapp/internal/domain/entity"
	"goapp/internal/domain/service"
)

// ReviewHandler handles review-related HTTP requests
type ReviewHandler struct {
	contentService *service.ContentService
}

// NewReviewHandler creates a new ReviewHandler
func NewReviewHandler(contentService *service.ContentService) *ReviewHandler {
	return &ReviewHandler{
		contentService: contentService,
	}
}

// RegisterRoutes registers review routes
func (h *ReviewHandler) RegisterRoutes(r *gin.Engine) {
	r.GET("/api/products/:id/reviews", h.GetReviews)
	r.POST("/api/products/:id/reviews", h.CreateReview)

	// HTML page
	r.GET("/reviews/:product_id", h.ReviewsPage)
	r.POST("/reviews/:product_id", h.PostReviewPage)
}

// GetReviews retrieves reviews for a product
func (h *ReviewHandler) GetReviews(c *gin.Context) {
	productIDStr := c.Param("id")
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	reviews, err := h.contentService.GetProductReviews(productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reviews)
}

// CreateReview creates a new review
// VULNERABLE: SQL Injection via comment + Stored XSS via ContentService deep call graph
func (h *ReviewHandler) CreateReview(c *gin.Context) {
	productIDStr := c.Param("id")
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	userIDStr := c.PostForm("user_id")
	userID, _ := strconv.ParseInt(userIDStr, 10, 64)
	if userID == 0 {
		userID = 1 // Default user
	}

	ratingStr := c.PostForm("rating")
	rating, _ := strconv.Atoi(ratingStr)
	if rating < 1 || rating > 5 {
		rating = 5
	}

	comment := c.PostForm("comment")

	// VULNERABLE: Stored XSS via ContentService deep call graph
	submission := &entity.ContentSubmission{
		Content:     comment,
		Author:      fmt.Sprintf("user_%d", userID),
		ContentType: "html",
		Context:     "review",
	}
	err = h.contentService.SubmitReview(productID, userID, rating, submission)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Review created"})
}

// ReviewsPage serves the reviews page for a product
// VULNERABLE: Stored XSS in reviews display
func (h *ReviewHandler) ReviewsPage(c *gin.Context) {
	productIDStr := c.Param("product_id")
	productID, _ := strconv.ParseInt(productIDStr, 10, 64)

	reviews, _ := h.contentService.GetProductReviews(productID)

	productName := fmt.Sprintf("Product #%d", productID)

	html := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head><title>Reviews for %s - Vulnerable App</title></head>
<body>
	<h1>Reviews for: %s</h1>
	<form method="POST" action="/reviews/%d">
		<input type="hidden" name="user_id" value="1">
		<label>Rating (1-5): <input type="number" name="rating" min="1" max="5" value="5"></label><br><br>
		<label>Comment: <textarea name="comment" rows="3" cols="50"></textarea></label><br><br>
		<button type="submit">Submit Review</button>
	</form>
	<h2>Existing Reviews</h2>
`, productName, productName, productID)

	for _, review := range reviews {
		// VULNERABLE: Stored XSS - comment displayed without encoding
		html += fmt.Sprintf(`
		<div style="border: 1px solid #ccc; padding: 10px; margin: 5px;">
			<b>Rating: %d/5</b>
			<p>%s</p>
		</div>
		`, review.Rating, review.Comment)
	}

	html += `
	<p>Hints:</p>
	<ul>
		<li>Stored XSS: Try &lt;script&gt;alert('XSS')&lt;/script&gt; in comment</li>
		<li>SQLi: Try '); DROP TABLE reviews; -- in comment</li>
	</ul>
</body>
</html>
`
	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, html)
}

// PostReviewPage handles form submission for reviews
func (h *ReviewHandler) PostReviewPage(c *gin.Context) {
	productIDStr := c.Param("product_id")
	productID, _ := strconv.ParseInt(productIDStr, 10, 64)

	userIDStr := c.PostForm("user_id")
	userID, _ := strconv.ParseInt(userIDStr, 10, 64)
	if userID == 0 {
		userID = 1
	}

	ratingStr := c.PostForm("rating")
	rating, _ := strconv.Atoi(ratingStr)
	if rating < 1 || rating > 5 {
		rating = 5
	}

	comment := c.PostForm("comment")

	submission := &entity.ContentSubmission{
		Content:     comment,
		Author:      fmt.Sprintf("user_%d", userID),
		ContentType: "html",
		Context:     "review",
	}
	h.contentService.SubmitReview(productID, userID, rating, submission)
	c.Redirect(http.StatusFound, fmt.Sprintf("/reviews/%d", productID))
}
