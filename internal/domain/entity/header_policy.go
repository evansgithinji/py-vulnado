package entity

// HeaderPolicy defines a policy for setting response headers.
type HeaderPolicy struct {
	Name         string
	AllowRaw     bool
	DefaultValue string
}

// HeaderRequest is a DTO for header customization requests.
type HeaderRequest struct {
	HeaderName  string
	HeaderValue string
	PolicyName  string
}

// LocaleConfig holds locale configuration settings.
type LocaleConfig struct {
	Code        string
	DisplayName string
	BaseURL     string
	URLPattern  string
}

// LocaleRequest is a DTO for localization redirect requests.
type LocaleRequest struct {
	Lang string
}
