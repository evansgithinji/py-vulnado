package usecase

import (
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"
)

// NetworkUseCase handles network operations
type NetworkUseCase struct {
	// URLCache stores fetched URLs and their content (VULNERABLE: info leakage)
	URLCache map[string]string
}

// NewNetworkUseCase creates a new NetworkUseCase
func NewNetworkUseCase() *NetworkUseCase {
	return &NetworkUseCase{
		URLCache: make(map[string]string),
	}
}

// PingHost pings a host using shell
// VULNERABLE: Command injection via host parameter with sh -c
func (uc *NetworkUseCase) PingHost(host string) (string, error) {
	// VULNERABLE: User input in shell command via fmt.Sprintf
	cmd := fmt.Sprintf("ping -c 1 %s", host)
	output, err := exec.Command("sh", "-c", cmd).CombinedOutput()
	return string(output), err
}

// PingHostSafe pings a host using array form (less exploitable)
// Note: Still potentially vulnerable but harder to exploit
func (uc *NetworkUseCase) PingHostSafe(host string) (string, error) {
	cmd := exec.Command("ping", "-c", "1", host)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

// FetchURL fetches content from a URL
// VULNERABLE: SSRF - no URL validation
func (uc *NetworkUseCase) FetchURL(url string) (string, error) {
	// VULNERABLE: Direct HTTP request
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

// TestWebhook tests a webhook URL
// VULNERABLE: SSRF via POST request
func (uc *NetworkUseCase) TestWebhook(webhookURL, payload string) (string, error) {
	// VULNERABLE: Arbitrary POST request
	resp, err := http.Post(webhookURL, "application/json", strings.NewReader(payload))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// ProxyRequest proxies a request to another URL
// VULNERABLE: SSRF + request smuggling potential
func (uc *NetworkUseCase) ProxyRequest(targetURL string) (*http.Response, error) {
	// VULNERABLE: Direct proxy without validation
	return http.Get(targetURL)
}

// FetchAndCacheURL fetches a URL and caches the content
// VULNERABLE: SSRF + information leakage via cache
func (uc *NetworkUseCase) FetchAndCacheURL(url string) (string, error) {
	// Check cache first
	if cached, ok := uc.URLCache[url]; ok {
		return cached, nil
	}

	// VULNERABLE: SSRF - no URL validation
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// VULNERABLE: Cache stores all fetched content
	uc.URLCache[url] = string(content)
	return string(content), nil
}

// GetCachedURLs returns all cached URLs
// VULNERABLE: Information disclosure
func (uc *NetworkUseCase) GetCachedURLs() map[string]string {
	return uc.URLCache
}
