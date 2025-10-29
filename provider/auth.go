package provider

import (
	"context"
	"encoding/json"
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

// CreateTicket creates a new ticket via the Freshservice API
func (c *Config) CreateTicket(ctx context.Context, ticket Ticket) (*Ticket, error) {
	ticketData := map[string]interface{}{
		"subject":     ticket.Subject,
		"description": ticket.Description,
		"priority":    ticket.Priority,
		"status":      ticket.Status,
	}

	// Add optional fields if they exist
	if ticket.Email != "" {
		ticketData["email"] = ticket.Email
	}
	if ticket.WorkspaceID != 0 {
		ticketData["workspace_id"] = ticket.WorkspaceID
	}
	if len(ticket.Assets) > 0 {
		ticketData["assets"] = ticket.Assets
	}

	body, err := json.Marshal(ticketData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal ticket data: %w", err)
	}

	req, err := c.NewRequest(ctx, "POST", "/tickets", strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}

	resp, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var ticketResp TicketResponse
	if err := json.NewDecoder(resp.Body).Decode(&ticketResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &ticketResp.Ticket, nil
}

// GetTicket retrieves a ticket by ID
func (c *Config) GetTicket(ctx context.Context, ticketID string) (*Ticket, error) {
	endpoint := fmt.Sprintf("/tickets/%s", ticketID)

	req, err := c.NewRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return nil, fmt.Errorf("ticket with ID %s not found", ticketID)
	}

	var ticketResp TicketResponse
	if err := json.NewDecoder(resp.Body).Decode(&ticketResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &ticketResp.Ticket, nil
}

// UpdateTicket updates an existing ticket
func (c *Config) UpdateTicket(ctx context.Context, ticketID string, ticket Ticket) (*Ticket, error) {
	ticketData := map[string]interface{}{
		"subject":     ticket.Subject,
		"description": ticket.Description,
		"priority":    ticket.Priority,
		"status":      ticket.Status,
	}

	// Add optional fields if they exist
	if ticket.Email != "" {
		ticketData["email"] = ticket.Email
	}
	if ticket.WorkspaceID != 0 {
		ticketData["workspace_id"] = ticket.WorkspaceID
	}
	if len(ticket.Assets) > 0 {
		ticketData["assets"] = ticket.Assets
	}

	body, err := json.Marshal(ticketData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal ticket data: %w", err)
	}

	endpoint := fmt.Sprintf("/tickets/%s", ticketID)
	req, err := c.NewRequest(ctx, "PUT", endpoint, strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}

	resp, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var ticketResp TicketResponse
	if err := json.NewDecoder(resp.Body).Decode(&ticketResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &ticketResp.Ticket, nil
}

// DeleteTicket deletes a ticket by ID
func (c *Config) DeleteTicket(ctx context.Context, ticketID string) error {
	endpoint := fmt.Sprintf("/tickets/%s", ticketID)

	req, err := c.NewRequest(ctx, "DELETE", endpoint, nil)
	if err != nil {
		return err
	}

	resp, err := c.DoRequest(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		// Ticket already deleted or doesn't exist
		return nil
	}

	return nil
}
