package config

import "testing"

func TestBuildWebSocketURL(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"https://example.com", "wss://example.com/api/agents/ws"},
		{"http://example.com", "ws://example.com/api/agents/ws"},
		{"https://example.com/api", "wss://example.com/api/agents/ws"},
		{"https://example.com/base", "wss://example.com/base/api/agents/ws"},
		{"wss://example.com", "wss://example.com/api/agents/ws"},
	}

	for _, tt := range tests {
		got, err := BuildWebSocketURL(tt.input)
		if err != nil {
			t.Fatalf("unexpected error for %s: %v", tt.input, err)
		}
		if got != tt.expected {
			t.Fatalf("input %s expected %s got %s", tt.input, tt.expected, got)
		}
	}
}

func TestBuildWebSocketURLInvalid(t *testing.T) {
	if _, err := BuildWebSocketURL("example.com"); err == nil {
		t.Fatalf("expected error for missing scheme")
	}
}
