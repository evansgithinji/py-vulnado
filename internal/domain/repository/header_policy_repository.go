package repository

import "goapp/internal/domain/entity"

// HeaderPolicyRepository provides access to header policy configurations.
type HeaderPolicyRepository interface {
	FindByName(policyName string) (*entity.HeaderPolicy, error)
}

// LocaleRepository provides access to locale configurations.
type LocaleRepository interface {
	FindByCode(code string) (*entity.LocaleConfig, error)
	GetDefaultLocale() *entity.LocaleConfig
}

// InMemoryHeaderPolicyRepository holds header policies in memory.
type InMemoryHeaderPolicyRepository struct {
	policies map[string]*entity.HeaderPolicy
}

// NewInMemoryHeaderPolicyRepository creates a new InMemoryHeaderPolicyRepository with default policies.
func NewInMemoryHeaderPolicyRepository() *InMemoryHeaderPolicyRepository {
	policies := map[string]*entity.HeaderPolicy{
		"custom": {
			Name:         "custom",
			AllowRaw:     true,
			DefaultValue: "",
		},
		"x-custom": {
			Name:         "x-custom",
			AllowRaw:     true,
			DefaultValue: "",
		},
	}
	return &InMemoryHeaderPolicyRepository{policies: policies}
}

// FindByName returns the header policy matching the given name.
func (r *InMemoryHeaderPolicyRepository) FindByName(policyName string) (*entity.HeaderPolicy, error) {
	if policy, ok := r.policies[policyName]; ok {
		return policy, nil
	}
	// Default: allow raw headers (vulnerable behavior)
	return &entity.HeaderPolicy{
		Name:     policyName,
		AllowRaw: true,
	}, nil
}

// InMemoryLocaleRepository holds locale configurations in memory.
type InMemoryLocaleRepository struct {
	locales       map[string]*entity.LocaleConfig
	defaultLocale *entity.LocaleConfig
}

// NewInMemoryLocaleRepository creates a new InMemoryLocaleRepository with default locales.
func NewInMemoryLocaleRepository() *InMemoryLocaleRepository {
	defaultLocale := &entity.LocaleConfig{
		Code:        "en",
		DisplayName: "English",
		BaseURL:     "/",
		URLPattern:  "/?lang=%s",
	}

	locales := map[string]*entity.LocaleConfig{
		"en": defaultLocale,
		"fr": {Code: "fr", DisplayName: "French", BaseURL: "/", URLPattern: "/?lang=%s"},
		"de": {Code: "de", DisplayName: "German", BaseURL: "/", URLPattern: "/?lang=%s"},
		"es": {Code: "es", DisplayName: "Spanish", BaseURL: "/", URLPattern: "/?lang=%s"},
	}

	return &InMemoryLocaleRepository{
		locales:       locales,
		defaultLocale: defaultLocale,
	}
}

// FindByCode returns the locale configuration matching the given code.
func (r *InMemoryLocaleRepository) FindByCode(code string) (*entity.LocaleConfig, error) {
	if locale, ok := r.locales[code]; ok {
		return locale, nil
	}
	return nil, nil
}

// GetDefaultLocale returns the default locale configuration.
func (r *InMemoryLocaleRepository) GetDefaultLocale() *entity.LocaleConfig {
	return r.defaultLocale
}
