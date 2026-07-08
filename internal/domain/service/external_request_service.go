package service

import (
	"net/http"

	"goapp/internal/domain/entity"
	"goapp/internal/domain/repository"
	"goapp/internal/domain/usecase"
)

// ExternalRequestService orchestrates external HTTP requests with deep call graph
type ExternalRequestService struct {
	policyRepo     *repository.UrlPolicyRepository
	requestBuilder *usecase.RequestBuilder
	httpClient     *usecase.HttpClientAdapter
}

// NewExternalRequestService creates a new ExternalRequestService
func NewExternalRequestService(
	policyRepo *repository.UrlPolicyRepository,
	requestBuilder *usecase.RequestBuilder,
	httpClient *usecase.HttpClientAdapter,
) *ExternalRequestService {
	return &ExternalRequestService{
		policyRepo:     policyRepo,
		requestBuilder: requestBuilder,
		httpClient:     httpClient,
	}
}

// ExecuteRequest executes an external HTTP request through the deep call graph
// Deep call graph for SSRF (depth 5):
// Handler -> ExternalRequestService.ExecuteRequest()
//   -> UrlPolicyRepository.GetPolicy() -> UrlPolicy
//   -> RequestBuilder.BuildRequest() -> PreparedRequest (VULNERABLE: no URL validation)
//   -> HttpClientAdapter.Execute() -> response (VULNERABLE: fetches any URL)
func (s *ExternalRequestService) ExecuteRequest(dto *entity.HttpRequestDto) (map[string]interface{}, error) {
	policy := s.policyRepo.GetPolicy("fetch")
	prepared := s.requestBuilder.BuildRequest(dto, policy)
	return s.httpClient.Execute(prepared)
}

// ProxyRequest proxies an HTTP request through the deep call graph
func (s *ExternalRequestService) ProxyRequest(url string) (*http.Response, error) {
	dto := &entity.HttpRequestDto{URL: url, Method: "GET", Timeout: 10}
	policy := s.policyRepo.GetPolicy("proxy")
	prepared := s.requestBuilder.BuildRequest(dto, policy)
	return s.httpClient.ExecuteProxy(prepared)
}

// TestWebhook tests a webhook URL through the deep call graph
func (s *ExternalRequestService) TestWebhook(url, payload string) (map[string]interface{}, error) {
	dto := &entity.HttpRequestDto{URL: url, Method: "POST", Body: payload, Timeout: 10}
	policy := s.policyRepo.GetPolicy("webhook")
	prepared := s.requestBuilder.BuildRequest(dto, policy)
	return s.httpClient.Execute(prepared)
}

// FetchRemoteFile fetches remote file content through the deep call graph
func (s *ExternalRequestService) FetchRemoteFile(url string) (string, error) {
	dto := &entity.HttpRequestDto{URL: url, Method: "GET", Timeout: 10}
	policy := s.policyRepo.GetPolicy("remote_file")
	prepared := s.requestBuilder.BuildRequest(dto, policy)
	return s.httpClient.ExecuteRaw(prepared)
}
