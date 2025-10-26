package reset

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// HTTPClient wraps http.Client with common configurations
type HTTPClient struct {
	client *http.Client
	token  string
}

// NewHTTPClient creates a new HTTP client with token
func NewHTTPClient(token string) *HTTPClient {
	return &HTTPClient{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		token: token,
	}
}

// DoRequest performs an HTTP request with automatic token injection
func (c *HTTPClient) DoRequest(method, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// Set common headers
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, text/plain, */*")

	return c.client.Do(req)
}
