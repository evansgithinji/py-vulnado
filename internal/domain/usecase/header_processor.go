package usecase

import (
	"fmt"

	"goapp/internal/domain/entity"
)

// HeaderValueProcessor processes header values according to policy.
type HeaderValueProcessor struct{}

// NewHeaderValueProcessor creates a new HeaderValueProcessor.
func NewHeaderValueProcessor() *HeaderValueProcessor {
	return &HeaderValueProcessor{}
}

// ProcessValue applies the header policy to the given value.
// VULNERABLE: HTTP Header Injection / CRLF (CWE-113) - no sanitization when AllowRaw is true
func (p *HeaderValueProcessor) ProcessValue(policy *entity.HeaderPolicy, name, value string) (string, string) {
	if policy.AllowRaw {
		// VULNERABLE: CRLF injection - user controls header value without sanitization
		return name, value
	}
	return name, policy.DefaultValue
}

// ResponseHeaderWriter writes processed headers to a response context.
type ResponseHeaderWriter struct{}

// NewResponseHeaderWriter creates a new ResponseHeaderWriter.
func NewResponseHeaderWriter() *ResponseHeaderWriter {
	return &ResponseHeaderWriter{}
}

// BuildHeaderPair constructs the final header name and value pair.
func (w *ResponseHeaderWriter) BuildHeaderPair(headerName, headerValue string) (string, string) {
	return headerName, headerValue
}

// CookieManager handles cookie-related header operations.
type CookieManager struct{}

// NewCookieManager creates a new CookieManager.
func NewCookieManager() *CookieManager {
	return &CookieManager{}
}

// BuildSetCookieValue constructs a Set-Cookie header value.
func (m *CookieManager) BuildSetCookieValue(name, value, path string) string {
	return fmt.Sprintf("%s=%s; Path=%s", name, value, path)
}

// RedirectBuilder constructs redirect URLs.
type RedirectBuilder struct{}

// NewRedirectBuilder creates a new RedirectBuilder.
func NewRedirectBuilder() *RedirectBuilder {
	return &RedirectBuilder{}
}

// BuildRedirectURL constructs a redirect URL with the language parameter.
// VULNERABLE: HTTP Header Injection / CRLF (CWE-113) - no sanitization of lang parameter
func (b *RedirectBuilder) BuildRedirectURL(locale *entity.LocaleConfig, lang string) string {
	// VULNERABLE: CRLF injection via lang parameter in Location header
	return fmt.Sprintf(locale.URLPattern, lang)
}
