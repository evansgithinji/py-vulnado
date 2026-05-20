package usecase

import (
	"fmt"

	"goapp/internal/domain/entity"
)

// SearchResultMapper maps raw query results to product entities
type SearchResultMapper struct{}

// NewSearchResultMapper creates a new SearchResultMapper
func NewSearchResultMapper() *SearchResultMapper {
	return &SearchResultMapper{}
}

// MapToProducts maps raw database rows to Product entities
func (m *SearchResultMapper) MapToProducts(rows []map[string]interface{}) []entity.Product {
	var products []entity.Product
	for _, row := range rows {
		product := entity.Product{}

		if id, ok := row["id"]; ok {
			product.ID = toInt64(id)
		}
		if name, ok := row["name"]; ok {
			product.Name = fmt.Sprintf("%v", name)
		}
		if price, ok := row["price"]; ok {
			product.Price = toFloat64(price)
		}
		if stock, ok := row["stock"]; ok {
			product.Stock = int(toInt64(stock))
		}
		if category, ok := row["category"]; ok {
			product.Category = fmt.Sprintf("%v", category)
		}

		products = append(products, product)
	}
	return products
}

func toInt64(v interface{}) int64 {
	switch val := v.(type) {
	case int64:
		return val
	case int:
		return int64(val)
	case float64:
		return int64(val)
	default:
		return 0
	}
}

func toFloat64(v interface{}) float64 {
	switch val := v.(type) {
	case float64:
		return val
	case int64:
		return float64(val)
	case int:
		return float64(val)
	default:
		return 0
	}
}
