package client

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

// InternalClient demonstrates how to call the internal API endpoints.
type InternalClient struct {
	HTTPClient *http.Client
	BaseURL    string
}

func NewInternalClient(baseURL string) *InternalClient {
	return &InternalClient{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// GetInternalData calls the /api/v1/internal endpoint.
func (c *InternalClient) GetInternalData(ctx context.Context, token string) (string, error) {
	url := fmt.Sprintf("%s/api/v1/internal", c.BaseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
