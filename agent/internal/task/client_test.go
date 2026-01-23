package task

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestClientPullTaskNoContent(t *testing.T) {
	client := &Client{
		baseURL: "http://example",
		apiKey:  "key",
		http: &http.Client{
			Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
				if r.URL.Path != "/api/agent/tasks/pull" {
					t.Fatalf("unexpected path %s", r.URL.Path)
				}
				return &http.Response{
					StatusCode: http.StatusNoContent,
					Body:       io.NopCloser(strings.NewReader("")),
					Header:     http.Header{},
				}, nil
			}),
		},
	}
	task, err := client.PullTask(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if task != nil {
		t.Fatalf("expected nil task")
	}
}

func TestClientPullTaskOK(t *testing.T) {
	client := &Client{
		baseURL: "http://example",
		apiKey:  "key",
		http: &http.Client{
			Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
				if r.Header.Get("X-Agent-Key") == "" {
					t.Fatalf("missing api key header")
				}
				body, _ := json.Marshal(Task{ID: 1})
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader(body)),
					Header:     http.Header{},
				}, nil
			}),
		},
	}
	task, err := client.PullTask(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if task == nil || task.ID != 1 {
		t.Fatalf("unexpected task")
	}
}

func TestClientUpdateStatus(t *testing.T) {
	client := &Client{
		baseURL: "http://example",
		apiKey:  "key",
		http: &http.Client{
			Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
				if r.Method != http.MethodPatch {
					t.Fatalf("expected PATCH")
				}
				if r.Header.Get("X-Agent-Key") == "" {
					t.Fatalf("missing api key header")
				}
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader("")),
					Header:     http.Header{},
				}, nil
			}),
		},
	}
	if err := client.UpdateStatus(context.Background(), 1, "completed", ""); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestClientPullTaskErrorStatus(t *testing.T) {
	client := &Client{
		baseURL: "http://example",
		apiKey:  "key",
		http: &http.Client{
			Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusBadRequest,
					Body:       io.NopCloser(strings.NewReader("bad")),
					Header:     http.Header{},
				}, nil
			}),
		},
	}
	if _, err := client.PullTask(context.Background()); err == nil {
		t.Fatalf("expected error for non-200 status")
	}
}

func TestClientPullTaskBadJSON(t *testing.T) {
	client := &Client{
		baseURL: "http://example",
		apiKey:  "key",
		http: &http.Client{
			Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader("{bad json")),
					Header:     http.Header{},
				}, nil
			}),
		},
	}
	if _, err := client.PullTask(context.Background()); err == nil {
		t.Fatalf("expected error for invalid json")
	}
}

func TestClientUpdateStatusIncludesErrorMessage(t *testing.T) {
	client := &Client{
		baseURL: "http://example",
		apiKey:  "key",
		http: &http.Client{
			Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
				body, err := io.ReadAll(r.Body)
				if err != nil {
					t.Fatalf("read body: %v", err)
				}
				var payload map[string]string
				if err := json.Unmarshal(body, &payload); err != nil {
					t.Fatalf("unmarshal body: %v", err)
				}
				if payload["status"] != "failed" {
					t.Fatalf("expected status failed")
				}
				if payload["errorMessage"] != "boom" {
					t.Fatalf("expected error message")
				}
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader("")),
					Header:     http.Header{},
				}, nil
			}),
		},
	}
	if err := client.UpdateStatus(context.Background(), 1, "failed", "boom"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestClientUpdateStatusErrorStatus(t *testing.T) {
	client := &Client{
		baseURL: "http://example",
		apiKey:  "key",
		http: &http.Client{
			Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       io.NopCloser(strings.NewReader("")),
					Header:     http.Header{},
				}, nil
			}),
		},
	}
	if err := client.UpdateStatus(context.Background(), 1, "completed", ""); err == nil {
		t.Fatalf("expected error for non-200 status")
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}
