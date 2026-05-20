package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"goapp/internal/domain/entity"
	"goapp/internal/domain/service"
	"goapp/internal/domain/usecase"
)

// ProductHandler handles product-related HTTP requests
type ProductHandler struct {
	catalogService  *service.CatalogService
	htmlBuilder     *usecase.HtmlResponseBuilder
}

// NewProductHandler creates a new ProductHandler
func NewProductHandler(catalogService *service.CatalogService, htmlBuilder *usecase.HtmlResponseBuilder) *ProductHandler {
	return &ProductHandler{
		catalogService: catalogService,
		htmlBuilder:    htmlBuilder,
	}
}

// RegisterRoutes registers product routes
func (h *ProductHandler) RegisterRoutes(r *gin.Engine) {
	r.GET("/api/products", h.ListProducts)
	r.GET("/api/products/search", h.SearchProducts)
	r.GET("/api/products/search/sorted", h.SearchProductsSorted)
	r.GET("/api/products/:id", h.GetProduct)

	// HTML pages
	r.GET("/products", h.ProductsPage)
	r.GET("/products/search", h.SearchPage)
	r.GET("/products/search/sorted", h.SortedSearchPage)
}

// ListProducts returns all products
func (h *ProductHandler) ListProducts(c *gin.Context) {
	request := &entity.SearchRequest{Query: ""}
	products, err := h.catalogService.Search(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}

// SearchProducts searches for products
// VULNERABLE: SQL Injection via query and category parameters
func (h *ProductHandler) SearchProducts(c *gin.Context) {
	query := c.Query("q")
	category := c.Query("category")

	var products []entity.Product
	var err error

	if category != "" {
		products, err = h.catalogService.SearchByCategory(category)
	} else {
		request := &entity.SearchRequest{Query: query, Category: category}
		products, err = h.catalogService.Search(request)
	}

	if err != nil {
		// VULNERABLE: Information disclosure
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "Search failed",
			"detail": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, products)
}

// SearchProductsSorted searches for products with sorting
// VULNERABLE: SQL Injection via ORDER BY clause
func (h *ProductHandler) SearchProductsSorted(c *gin.Context) {
	query := c.Query("q")
	category := c.Query("category")
	sortBy := c.Query("sort")
	order := c.Query("order")

	request := &entity.SearchRequest{
		Query:    query,
		Category: category,
		SortBy:   sortBy,
		Order:    order,
	}

	products, err := h.catalogService.Search(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "Search failed",
			"detail": err.Error(),
			"query":  fmt.Sprintf("sort=%s, order=%s", sortBy, order),
		})
		return
	}

	c.JSON(http.StatusOK, products)
}

// GetProduct retrieves a product by ID
func (h *ProductHandler) GetProduct(c *gin.Context) {
	productIDStr := c.Param("id")
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	request := &entity.SearchRequest{Query: productIDStr}
	products, err := h.catalogService.Search(request)
	if err != nil || len(products) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	for _, p := range products {
		if p.ID == productID {
			c.JSON(http.StatusOK, p)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
}

// ProductsPage serves the products listing page
func (h *ProductHandler) ProductsPage(c *gin.Context) {
	request := &entity.SearchRequest{Query: ""}
	products, _ := h.catalogService.Search(request)

	html := `
<!DOCTYPE html>
<html>
<head><title>Products - Vulnerable App</title></head>
<body>
	<h1>Products</h1>
	<a href="/products/search">Search Products</a>
	<table border="1">
		<tr><th>ID</th><th>Name</th><th>Price</th><th>Stock</th><th>Category</th></tr>
`
	for _, p := range products {
		html += fmt.Sprintf("<tr><td>%d</td><td>%s</td><td>$%.2f</td><td>%d</td><td>%s</td></tr>",
			p.ID, p.Name, p.Price, p.Stock, p.Category)
	}

	html += `
	</table>
</body>
</html>
`
	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, html)
}

// SearchPage serves the product search page
// VULNERABLE: Reflected XSS via search parameters
func (h *ProductHandler) SearchPage(c *gin.Context) {
	query := c.Query("q")
	category := c.Query("category")

	var products []entity.Product
	if category != "" {
		products, _ = h.catalogService.SearchByCategory(category)
	} else {
		request := &entity.SearchRequest{Query: query, Category: category}
		products, _ = h.catalogService.Search(request)
	}

	// VULNERABLE: Reflected XSS via HtmlResponseBuilder
	if query != "" || category != "" {
		searchTerm := query
		if searchTerm == "" {
			searchTerm = category
		}
		htmlContent := h.htmlBuilder.BuildSearchPage(searchTerm, products)
		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, htmlContent)
		return
	}

	// VULNERABLE: Reflected XSS - query and category displayed without encoding
	html := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head><title>Product Search - Vulnerable App</title></head>
<body>
	<h1>Product Search</h1>
	<form method="GET" action="/products/search">
		<label>Search: <input type="text" name="q" value="%s"></label><br><br>
		<label>Category: <input type="text" name="category" value="%s"></label><br><br>
		<button type="submit">Search</button>
	</form>
	<h2>Search Results for: %s (Category: %s)</h2>
`, query, category, query, category) // VULNERABLE: Reflected XSS

	html += `<table border="1">
		<tr><th>ID</th><th>Name</th><th>Price</th><th>Stock</th><th>Category</th></tr>`

	for _, p := range products {
		html += fmt.Sprintf("<tr><td>%d</td><td>%s</td><td>$%.2f</td><td>%d</td><td>%s</td></tr>",
			p.ID, p.Name, p.Price, p.Stock, p.Category)
	}

	html += `
	</table>
	<p>Hint: Try ' OR '1'='1 in search, or ' UNION SELECT ... for data extraction</p>
</body>
</html>
`
	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, html)
}

// SortedSearchPage serves the sorted product search page
// VULNERABLE: ORDER BY injection
func (h *ProductHandler) SortedSearchPage(c *gin.Context) {
	query := c.Query("q")
	category := c.Query("category")
	sortBy := c.Query("sort")
	order := c.Query("order")

	request := &entity.SearchRequest{
		Query:    query,
		Category: category,
		SortBy:   sortBy,
		Order:    order,
	}
	products, _ := h.catalogService.Search(request)

	// VULNERABLE: Reflected XSS via all parameters
	html := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head><title>Sorted Product Search - Vulnerable App</title></head>
<body>
	<h1>Product Search with Sorting (ORDER BY Injection Demo)</h1>
	<form method="GET" action="/products/search/sorted">
		<label>Search: <input type="text" name="q" value="%s"></label><br><br>
		<label>Category: <input type="text" name="category" value="%s"></label><br><br>
		<label>Sort By: <input type="text" name="sort" value="%s" placeholder="name, price, stock"></label><br><br>
		<label>Order: <input type="text" name="order" value="%s" placeholder="ASC, DESC"></label><br><br>
		<button type="submit">Search</button>
	</form>
	<h2>Results (sorted by: %s %s)</h2>
`, query, category, sortBy, order, sortBy, order)

	html += `<table border="1">
		<tr><th>ID</th><th>Name</th><th>Price</th><th>Stock</th><th>Category</th></tr>`

	for _, p := range products {
		html += fmt.Sprintf("<tr><td>%d</td><td>%s</td><td>$%.2f</td><td>%d</td><td>%s</td></tr>",
			p.ID, p.Name, p.Price, p.Stock, p.Category)
	}

	html += `
	</table>
	<p>Hints:</p>
	<ul>
		<li>ORDER BY injection: Try sort=name;--</li>
		<li>Boolean-based: sort=(CASE WHEN 1=1 THEN name ELSE price END)</li>
		<li>Error-based: sort=name,(SELECT 1 FROM users)</li>
	</ul>
</body>
</html>
`
	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, html)
}
