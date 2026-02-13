package server

import "context"

// ServerClient defines the interface for Worker to communicate with Server.
// This interface allows for easier testing and decoupling.
type ServerClient interface {
	// GetProviderConfig fetches tool-specific configuration (e.g., API keys for subfinder)
	GetProviderConfig(ctx context.Context, scanID int, toolName string) (*ProviderConfig, error)

	// EnsureWordlistLocal ensures a wordlist file exists locally, downloading if needed
	EnsureWordlistLocal(ctx context.Context, wordlistName, basePath string) (string, error)

	// PostBatch sends a batch of data to the server (used by BatchSender)
	PostBatch(ctx context.Context, scanID, targetID int, dataType string, items []any) error
}

// ProviderConfig contains tool configuration file content
type ProviderConfig struct {
	Content string `json:"content"`
}
