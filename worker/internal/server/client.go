package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/orbit/worker/internal/pkg"
	"go.uber.org/zap"
)

// Client handles all HTTP communication with Server
// Implements Provider, ResultSaver, and StatusUpdater interfaces
type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
	maxRetries int
}

// NewClient creates a new server client
func NewClient(baseURL, token string) *Client {
	return &Client{
		baseURL: baseURL,
		token:   token,
		httpClient: &http.Client{
			Timeout: 5 * time.Minute,
		},
		maxRetries: 3,
	}
}

// --- HTTP helpers ---

func (c *Client) get(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("X-Worker-Token", c.token)
	req.Header.Set("Accept", "application/json")
	return c.httpClient.Do(req)
}

func (c *Client) postWithRetry(ctx context.Context, url string, body any) error {
	return c.doWithRetry(ctx, "POST", url, body)
}

func (c *Client) doWithRetry(ctx context.Context, method, url string, body any) error {
	var lastErr error
	for i := 0; i < c.maxRetries; i++ {
		// Check context before retry
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if i > 0 {
			// Use select to allow cancellation during sleep
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(time.Duration(1<<i) * time.Second):
			}
		}

		if err := c.doRequest(ctx, method, url, body); err == nil {
			return nil
		} else {
			lastErr = err
			pkg.Logger.Warn("API call failed, retrying",
				zap.String("url", url),
				zap.Int("attempt", i+1),
				zap.Error(err))
		}
	}

	pkg.Logger.Error("All retries failed", zap.String("url", url), zap.Error(lastErr))
	return lastErr
}

func (c *Client) doRequest(ctx context.Context, method, url string, body any) error {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Worker-Token", c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error: status=%d, body=%s", resp.StatusCode, string(respBody))
	}

	return nil
}

func fetchJSON[T any](ctx context.Context, c *Client, url string) (T, error) {
	var result T

	pkg.Logger.Debug("Fetching JSON", zap.String("url", url))

	resp, err := c.get(ctx, url)
	if err != nil {
		return result, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return result, fmt.Errorf("server error: status=%d, body=%s", resp.StatusCode, string(body))
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return result, fmt.Errorf("failed to decode JSON: %w", err)
	}

	return result, nil
}

// PostBatch sends a batch of items to the server
func (c *Client) PostBatch(ctx context.Context, scanID, targetID int, dataType string, items []any) error {
	var url string
	var body map[string]any

	switch dataType {
	case "subdomain":
		url = fmt.Sprintf("%s/api/worker/scans/%d/subdomains/bulk-upsert", c.baseURL, scanID)
		body = map[string]any{
			"targetId":   targetID,
			"subdomains": items,
		}
	case "website":
		url = fmt.Sprintf("%s/api/worker/scans/%d/websites/bulk-upsert", c.baseURL, scanID)
		body = map[string]any{
			"targetId": targetID,
			"websites": items,
		}
	case "endpoint":
		url = fmt.Sprintf("%s/api/worker/scans/%d/endpoints/bulk-upsert", c.baseURL, scanID)
		body = map[string]any{
			"targetId":  targetID,
			"endpoints": items,
		}
	default:
		url = fmt.Sprintf("%s/api/worker/scans/%d/%ss/bulk-upsert", c.baseURL, scanID, dataType)
		body = map[string]any{
			"targetId": targetID,
			"items":    items,
		}
	}

	return c.postWithRetry(ctx, url, body)
}
