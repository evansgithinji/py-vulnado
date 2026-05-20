package service

import (
	"goapp/internal/domain/entity"
	"goapp/internal/domain/repository"
	"goapp/internal/domain/usecase"
)

// NavigationService orchestrates URL redirect resolution with deep call graph
type NavigationService struct {
	policyRepo  *repository.RedirectPolicyRepository
	urlResolver *usecase.UrlResolver
}

// NewNavigationService creates a new NavigationService
func NewNavigationService(
	policyRepo *repository.RedirectPolicyRepository,
	urlResolver *usecase.UrlResolver,
) *NavigationService {
	return &NavigationService{
		policyRepo:  policyRepo,
		urlResolver: urlResolver,
	}
}

// ResolveRedirect resolves a redirect URL through the deep call graph
// Deep call graph for Open Redirect (depth 4):
// Handler -> NavigationService.ResolveRedirect()
//   -> RedirectPolicyRepository.GetPolicy() -> RedirectPolicy
//   -> UrlResolver.ResolveUrl() -> finalUrl (VULNERABLE: allows external)
func (s *NavigationService) ResolveRedirect(request *entity.RedirectRequest) string {
	policy := s.policyRepo.GetPolicy(request.Context)
	return s.urlResolver.ResolveUrl(request.TargetURL, policy)
}
