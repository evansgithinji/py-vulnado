package usecase

import (
	"io"
	"net/http"
	"strings"
	"time"

	"goapp/internal/domain/entity"
)

// HttpClientAdapter executes HTTP requests
type HttpClientAdapter struct{}

// NewHttpClientAdapter creates a new HttpClientAdapter
func NewHttpClientAdapter() *HttpClientAdapter {
	return &HttpClientAdapter{}
}

// Execute executes a prepared HTTP request and returns structured response
// VULNERABLE: Fetches any URL without restriction (SSRF)
func (a *HttpClientAdapter) Execute(prepared *entity.PreparedRequest) (map[string]interface{}, error) {
	client := &http.Client{
		Timeout: time.Duration(prepared.Timeout) * time.Second,
	}

	var resp *http.Response
	var err error

	if strings.ToUpper(prepared.Method) == "POST" {
		req, reqErr := http.NewRequest("POST", prepared.URL, strings.NewReader(prepared.Body))
		if reqErr != nil {
			return nil, reqErr
		}
		for k, v := range prepared.Headers {
			req.Header.Set(k, v)
		}
		if req.Header.Get("Content-Type") == "" {
			req.Header.Set("Content-Type", "application/json")
		}
		resp, err = client.Do(req)
	} else {
		resp, err = client.Get(prepared.URL)
	}

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"status_code": resp.StatusCode,
		"body":        string(body),
	}, nil
}

// ExecuteRaw executes a prepared HTTP request and returns raw body string
// VULNERABLE: Returns raw response body (SSRF)
func (a *HttpClientAdapter) ExecuteRaw(prepared *entity.PreparedRequest) (string, error) {
	client := &http.Client{
		Timeout: time.Duration(prepared.Timeout) * time.Second,
	}

	resp, err := client.Get(prepared.URL)
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

// ExecuteProxy executes a prepared HTTP request and returns the raw http.Response
// VULNERABLE: Direct proxy without validation
func (a *HttpClientAdapter) ExecuteProxy(prepared *entity.PreparedRequest) (*http.Response, error) {
	return http.Get(prepared.URL)
}
