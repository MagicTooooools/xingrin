package task

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/yyhuni/lunafox/agent/internal/domain"
)

// Client handles HTTP API requests to the server.
type Client struct {
	baseURL string
	apiKey  string
	http    *http.Client
}

// NewClient creates a new task client.
func NewClient(serverURL, apiKey string) *Client {
	return &Client{
		baseURL: strings.TrimRight(serverURL, "/"),
		apiKey:  apiKey,
		http: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// PullTask requests a task from the server. Returns nil when no task available.
func (c *Client) PullTask(ctx context.Context) (*domain.Task, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/api/agent/tasks/pull", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Agent-Key", c.apiKey)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("pull task failed: status %d", resp.StatusCode)
	}

	var task domain.Task
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		return nil, err
	}
	return &task, nil
}

// UpdateStatus reports task status to the server with retry.
func (c *Client) UpdateStatus(ctx context.Context, taskID int, status, errorMessage string) error {
	payload := map[string]string{
		"status": status,
	}
	if errorMessage != "" {
		payload["errorMessage"] = errorMessage
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	var lastErr error
	for attempt := 0; attempt < 3; attempt++ {
		if attempt > 0 {
			backoff := time.Duration(5<<attempt) * time.Second // 5s, 10s, 20s
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(backoff):
			}
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPatch, fmt.Sprintf("%s/api/agent/tasks/%d/status", c.baseURL, taskID), bytes.NewReader(body))
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Agent-Key", c.apiKey)

		resp, err := c.http.Do(req)
		if err != nil {
			lastErr = err
			continue
		}
		resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			return nil
		}
		lastErr = fmt.Errorf("update status failed: status %d", resp.StatusCode)

		// Don't retry 4xx client errors (except 429)
		if resp.StatusCode >= 400 && resp.StatusCode < 500 && resp.StatusCode != 429 {
			return lastErr
		}
	}
	return lastErr
}
