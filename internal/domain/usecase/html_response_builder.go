package usecase

import (
	"fmt"

	"goapp/internal/domain/entity"
)

// HtmlResponseBuilder builds HTML responses with user data
type HtmlResponseBuilder struct{}

// NewHtmlResponseBuilder creates a new HtmlResponseBuilder
func NewHtmlResponseBuilder() *HtmlResponseBuilder {
	return &HtmlResponseBuilder{}
}

// BuildSearchPage builds an HTML search results page
// VULNERABLE: Reflected XSS - no HTML escaping of query
func (b *HtmlResponseBuilder) BuildSearchPage(query string, products []entity.Product) string {
	html := fmt.Sprintf("<h2>Products matching: %s</h2>", query) // VULNERABLE: XSS
	html += "<table border='1'><tr><th>Name</th><th>Price</th><th>Category</th></tr>"
	for _, p := range products {
		html += fmt.Sprintf("<tr><td>%s</td><td>$%.2f</td><td>%s</td></tr>", p.Name, p.Price, p.Category)
	}
	html += "</table>"
	return html
}

// BuildUserSearchPage builds an HTML user search results page
// VULNERABLE: Reflected XSS - no HTML escaping of query
func (b *HtmlResponseBuilder) BuildUserSearchPage(query string, users []entity.User) string {
	html := fmt.Sprintf("<h2>Search results for: %s</h2><ul>", query) // VULNERABLE: XSS
	for _, u := range users {
		html += fmt.Sprintf("<li>%s - %s</li>", u.Username, u.Email)
	}
	html += "</ul>"
	return html
}

// BuildErrorPage builds an HTML error page
// VULNERABLE: Reflected XSS - no HTML escaping
func (b *HtmlResponseBuilder) BuildErrorPage(query string, errMsg string) string {
	return fmt.Sprintf("<h2>Products matching: %s</h2><p>Error: %s</p>", query, errMsg)
}

// BuildPingPage builds an HTML ping results page
// VULNERABLE: Reflected XSS - no HTML escaping
func (b *HtmlResponseBuilder) BuildPingPage(host string, output string) string {
	return fmt.Sprintf("<h2>Ping results for: %s</h2><pre>%s</pre>", host, output)
}
