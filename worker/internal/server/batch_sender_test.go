package server

import (
	"context"
	"testing"

	"github.com/orbit/worker/internal/pkg"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type failingClient struct{}

func (f failingClient) GetProviderConfig(ctx context.Context, scanID int, toolName string) (*ProviderConfig, error) {
	return nil, nil
}

func (f failingClient) EnsureWordlistLocal(ctx context.Context, wordlistName, basePath string) (string, error) {
	return "", nil
}

func (f failingClient) PostBatch(ctx context.Context, scanID, targetID int, dataType string, items []any) error {
	return &HTTPError{StatusCode: 500, Body: "server error"}
}

func TestBatchSender_RequeuesOnFailure(t *testing.T) {
	prevLogger := pkg.Logger
	pkg.Logger = zap.NewNop()
	t.Cleanup(func() {
		pkg.Logger = prevLogger
	})
	ctx := context.Background()
	sender := NewBatchSender(ctx, failingClient{}, 1, 2, "subdomain", 2)

	require.NoError(t, sender.Add(map[string]string{"name": "a.example.com"}))
	err := sender.Add(map[string]string{"name": "b.example.com"})
	require.Error(t, err)

	sender.mu.Lock()
	defer sender.mu.Unlock()
	require.Len(t, sender.batch, 2)
}
