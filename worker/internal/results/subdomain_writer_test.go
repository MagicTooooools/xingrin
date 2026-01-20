package results

import (
	"context"
	"testing"

	"github.com/orbit/worker/internal/pkg"
	"github.com/orbit/worker/internal/server"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type captureClient struct {
	t       *testing.T
	batches [][]Subdomain
}

func (c *captureClient) GetProviderConfig(ctx context.Context, scanID int, toolName string) (*server.ProviderConfig, error) {
	return nil, nil
}

func (c *captureClient) EnsureWordlistLocal(ctx context.Context, wordlistName, basePath string) (string, error) {
	return "", nil
}

func (c *captureClient) PostBatch(ctx context.Context, scanID, targetID int, dataType string, items []any) error {
	require.Equal(c.t, "subdomain", dataType)
	batch := make([]Subdomain, 0, len(items))
	for _, item := range items {
		sd, ok := item.(Subdomain)
		require.True(c.t, ok)
		batch = append(batch, sd)
	}
	c.batches = append(c.batches, batch)
	return nil
}

func TestWriteSubdomains(t *testing.T) {
	prevLogger := pkg.Logger
	pkg.Logger = zap.NewNop()
	t.Cleanup(func() {
		pkg.Logger = prevLogger
	})

	ch := make(chan Subdomain, 2)
	ch <- Subdomain{Name: "a.example.com"}
	ch <- Subdomain{Name: "b.example.com"}
	close(ch)

	client := &captureClient{t: t}
	items, batches, err := WriteSubdomains(context.Background(), client, 1, 2, ch)
	require.NoError(t, err)
	require.Equal(t, 2, items)
	require.Equal(t, 1, batches)
	require.Len(t, client.batches, 1)
	require.Equal(t, []Subdomain{{Name: "a.example.com"}, {Name: "b.example.com"}}, client.batches[0])
}
