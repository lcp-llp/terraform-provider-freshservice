package provider

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Config holds the provider configuration
type Config struct {
	APIKey  string
	Domain  string
	BaseURL string
	Client  *http.Client
}

// NewConfig creates a new configuration instance
func NewConfig(apiKey, domain string) (*Config, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("api_key is required")
	}
	if domain == "" {
		return nil, fmt.Errorf("domain is required")
	}

	// Ensure domain has the correct format
	if !strings.HasSuffix(domain, ".freshservice.com") {
		if !strings.Contains(domain, ".") {
			domain = domain + ".freshservice.com"
		}
	}

	baseURL := fmt.Sprintf("https://%s/api/v2", domain)

	config := &Config{
		APIKey:  apiKey,
		Domain:  domain,
		BaseURL: baseURL,
		Client:  &http.Client{},
	}

	return config, nil
}

// NewRequest creates a new HTTP request with proper authentication
func (c *Config) NewRequest(ctx context.Context, method, endpoint string, body io.Reader) (*http.Request, error) {
	url := fmt.Sprintf("%s%s", c.BaseURL, endpoint)

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set authentication using Basic Auth with API key as username and "X" as password
	req.SetBasicAuth(c.APIKey, "X")

	// Set required headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	return req, nil
}

// DoRequest executes an HTTP request and returns the response
func (c *Config) DoRequest(req *http.Request) (*http.Response, error) {
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	// For 404 errors, return the response so caller can handle it
	if resp.StatusCode == 404 {
		return resp, nil
	}

	// Check for other API errors (but not 404)
	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, resp.Status)
	}

	return resp, nil
}
