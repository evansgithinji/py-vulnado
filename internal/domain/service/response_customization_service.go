package service

import (
	"goapp/internal/domain/entity"
	"goapp/internal/domain/repository"
	"goapp/internal/domain/usecase"
)

// ResponseCustomizationService orchestrates custom header setting operations.
type ResponseCustomizationService struct {
	policyRepo   repository.HeaderPolicyRepository
	processor    *usecase.HeaderValueProcessor
	headerWriter *usecase.ResponseHeaderWriter
}

// NewResponseCustomizationService creates a new ResponseCustomizationService.
func NewResponseCustomizationService(
	policyRepo repository.HeaderPolicyRepository,
	processor *usecase.HeaderValueProcessor,
	headerWriter *usecase.ResponseHeaderWriter,
) *ResponseCustomizationService {
	return &ResponseCustomizationService{
		policyRepo:   policyRepo,
		processor:    processor,
		headerWriter: headerWriter,
	}
}

// SetCustomHeader processes and returns a header name-value pair.
// Call graph: Service.SetCustomHeader → policyRepo.FindByName → processor.ProcessValue → headerWriter.BuildHeaderPair
func (s *ResponseCustomizationService) SetCustomHeader(req *entity.HeaderRequest) (string, string) {
	// Layer 3: Repository lookup for policy
	policy, err := s.policyRepo.FindByName(req.PolicyName)
	if err != nil || policy == nil {
		policy = &entity.HeaderPolicy{Name: "default", AllowRaw: true}
	}

	// Layer 4: Process header value according to policy
	processedName, processedValue := s.processor.ProcessValue(policy, req.HeaderName, req.HeaderValue)

	// Layer 5: Build final header pair
	return s.headerWriter.BuildHeaderPair(processedName, processedValue)
}

// LocalizationService orchestrates language redirect operations.
type LocalizationService struct {
	localeRepo      repository.LocaleRepository
	redirectBuilder *usecase.RedirectBuilder
	cookieManager   *usecase.CookieManager
}

// NewLocalizationService creates a new LocalizationService.
func NewLocalizationService(
	localeRepo repository.LocaleRepository,
	redirectBuilder *usecase.RedirectBuilder,
	cookieManager *usecase.CookieManager,
) *LocalizationService {
	return &LocalizationService{
		localeRepo:      localeRepo,
		redirectBuilder: redirectBuilder,
		cookieManager:   cookieManager,
	}
}

// GetRedirectURL builds a redirect URL for the given locale request.
// Call graph: Service.GetRedirectURL → localeRepo.FindByCode → redirectBuilder.BuildRedirectURL
func (s *LocalizationService) GetRedirectURL(req *entity.LocaleRequest) string {
	// Layer 3: Repository lookup for locale
	locale, err := s.localeRepo.FindByCode(req.Lang)
	if err != nil || locale == nil {
		locale = s.localeRepo.GetDefaultLocale()
	}

	// Layer 4: Build redirect URL (vulnerable - no sanitization of lang)
	return s.redirectBuilder.BuildRedirectURL(locale, req.Lang)
}
