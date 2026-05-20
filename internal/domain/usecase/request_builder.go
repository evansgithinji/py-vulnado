package usecase

import "goapp/internal/domain/entity"

// RequestBuilder builds HTTP requests from DTOs and policies
type RequestBuilder struct{}

// NewRequestBuilder creates a new RequestBuilder
func NewRequestBuilder() *RequestBuilder {
	return &RequestBuilder{}
}

// BuildRequest builds a prepared HTTP request
// VULNERABLE: No URL validation - allows internal IPs, metadata endpoints
func (b *RequestBuilder) BuildRequest(dto *entity.HttpRequestDto, policy *entity.UrlPolicy) *entity.PreparedRequest {
	return &entity.PreparedRequest{
		URL:     dto.URL,
		Method:  dto.Method,
		Headers: dto.Headers,
		Body:    dto.Body,
		Timeout: dto.Timeout,
		Policy:  policy,
	}
}
