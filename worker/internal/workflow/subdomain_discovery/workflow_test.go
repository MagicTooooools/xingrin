package subdomain_discovery

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yyhuni/lunafox/worker/internal/server"
	"github.com/yyhuni/lunafox/worker/internal/workflow"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type providerClient struct {
	config *server.ProviderConfig
	err    error
}

func (p providerClient) GetProviderConfig(ctx context.Context, scanID int, toolName string) (*server.ProviderConfig, error) {
	if p.err != nil {
		return nil, p.err
	}
	return p.config, nil
}

func (p providerClient) EnsureWordlistLocal(ctx context.Context, wordlistName, basePath string) (string, error) {
	return "", nil
}

func (p providerClient) PostBatch(ctx context.Context, scanID, targetID int, dataType string, items []any) error {
	return nil
}

type capturePostClient struct {
	calls int
	items []any
	err   error
}

func (c *capturePostClient) GetProviderConfig(ctx context.Context, scanID int, toolName string) (*server.ProviderConfig, error) {
	return nil, nil
}

func (c *capturePostClient) EnsureWordlistLocal(ctx context.Context, wordlistName, basePath string) (string, error) {
	return "", nil
}

func (c *capturePostClient) PostBatch(ctx context.Context, scanID, targetID int, dataType string, items []any) error {
	c.calls++
	c.items = append(c.items, items...)
	return c.err
}

func validScanConfig() map[string]any {
	return map[string]any{
		stageRecon: map[string]any{
			"enabled": true,
			"tools": map[string]any{
				toolSubfinder: map[string]any{
					"enabled":         true,
					"timeout-runtime": 3600,
					"threads-cli":     10,
				},
				toolSublist3r: map[string]any{
					"enabled":         true,
					"timeout-runtime": 3600,
					"threads-cli":     10,
				},
				toolAssetfinder: map[string]any{
					"enabled":         true,
					"timeout-runtime": 3600,
				},
			},
		},
		stageBruteforce: map[string]any{
			"enabled": false,
			"tools": map[string]any{
				toolSubdomainBruteforce: map[string]any{"enabled": false},
			},
		},
		stagePermutation: map[string]any{
			"enabled": false,
			"tools": map[string]any{
				toolSubdomainPermutationResolve: map[string]any{"enabled": false},
			},
		},
		stageResolve: map[string]any{
			"enabled": false,
			"tools": map[string]any{
				toolSubdomainResolve: map[string]any{"enabled": false},
			},
		},
	}
}

func TestInitializeMissingConfig(t *testing.T) {
	withNopLogger(t)
	w := New(t.TempDir())
	_, err := w.initialize(&workflow.Params{
		ScanConfig:   nil,
		TargetType:   "domain",
		TargetName:   "example.com",
		WorkDir:      t.TempDir(),
		ServerClient: providerClient{},
	})
	require.Error(t, err)
}

func TestInitializeInvalidTargetType(t *testing.T) {
	withNopLogger(t)
	w := New(t.TempDir())
	_, err := w.initialize(&workflow.Params{
		ScanConfig:   validScanConfig(),
		TargetType:   "ip",
		TargetName:   "1.1.1.1",
		WorkDir:      t.TempDir(),
		ServerClient: providerClient{},
	})
	require.Error(t, err)
}

func TestInitializeInvalidDomain(t *testing.T) {
	withNopLogger(t)
	w := New(t.TempDir())
	_, err := w.initialize(&workflow.Params{
		ScanConfig:   validScanConfig(),
		TargetType:   "domain",
		TargetName:   "bad domain",
		WorkDir:      t.TempDir(),
		ServerClient: providerClient{},
	})
	require.Error(t, err)
}

func TestInitializeNormalizesDomain(t *testing.T) {
	withNopLogger(t)
	w := New(t.TempDir())
	ctx, err := w.initialize(&workflow.Params{
		ScanConfig:   validScanConfig(),
		TargetType:   "domain",
		TargetName:   "Example.COM.",
		WorkDir:      t.TempDir(),
		ServerClient: providerClient{},
	})
	require.NoError(t, err)
	require.Equal(t, []string{"example.com"}, ctx.domains)
}

func TestInitializeNestedConfig(t *testing.T) {
	withNopLogger(t)
	w := New(t.TempDir())
	nested := map[string]any{
		Name: validScanConfig(),
	}

	ctx, err := w.initialize(&workflow.Params{
		ScanConfig:   nested,
		TargetType:   "domain",
		TargetName:   "example.com",
		WorkDir:      t.TempDir(),
		ServerClient: providerClient{},
	})
	require.NoError(t, err)
	require.NotNil(t, ctx)
}

func TestInitializeProviderConfigWritten(t *testing.T) {
	withNopLogger(t)
	w := New(t.TempDir())
	workDir := t.TempDir()

	ctx, err := w.initialize(&workflow.Params{
		ScanConfig:   validScanConfig(),
		TargetType:   "domain",
		TargetName:   "example.com",
		WorkDir:      workDir,
		ServerClient: providerClient{config: &server.ProviderConfig{Content: "api: key"}},
	})
	require.NoError(t, err)
	require.NotEmpty(t, ctx.providerConfigPath)

	data, err := os.ReadFile(ctx.providerConfigPath)
	require.NoError(t, err)
	assert.Equal(t, "api: key", string(data))
}

func TestSetupProviderConfigNoContent(t *testing.T) {
	withNopLogger(t)
	w := New(t.TempDir())
	workDir := t.TempDir()

	path, err := w.setupProviderConfig(context.Background(), &workflow.Params{
		ScanID:       1,
		ServerClient: providerClient{config: nil},
	}, workDir)
	require.NoError(t, err)
	assert.Equal(t, "", path)

	path, err = w.setupProviderConfig(context.Background(), &workflow.Params{
		ScanID:       1,
		ServerClient: providerClient{config: &server.ProviderConfig{Content: ""}},
	}, workDir)
	require.NoError(t, err)
	assert.Equal(t, "", path)
}

func TestSetupProviderConfigWriteError(t *testing.T) {
	withNopLogger(t)
	w := New(t.TempDir())
	workDir := t.TempDir()

	require.NoError(t, os.Chmod(workDir, 0500))
	t.Cleanup(func() {
		_ = os.Chmod(workDir, 0700)
	})

	_, err := w.setupProviderConfig(context.Background(), &workflow.Params{
		ScanID:       1,
		ServerClient: providerClient{config: &server.ProviderConfig{Content: "api: key"}},
	}, workDir)
	require.Error(t, err)
}

func TestInitializeProviderConfigErrorIgnored(t *testing.T) {
	withNopLogger(t)
	w := New(t.TempDir())
	_, err := w.initialize(&workflow.Params{
		ScanConfig:   validScanConfig(),
		TargetType:   "domain",
		TargetName:   "example.com",
		WorkDir:      t.TempDir(),
		ServerClient: providerClient{err: errors.New("boom")},
	})
	require.NoError(t, err)
}

func TestSaveResultsNoFiles(t *testing.T) {
	withNopLogger(t)
	w := New(t.TempDir())
	output := &workflow.Output{
		Data:    "not-files",
		Metrics: &workflow.Metrics{},
	}

	err := w.SaveResults(context.Background(), providerClient{}, &workflow.Params{}, output)
	require.NoError(t, err)
}

func TestSaveResultsSuccessUpdatesMetrics(t *testing.T) {
	withNopLogger(t)
	w := New(t.TempDir())
	dir := t.TempDir()
	file := filepath.Join(dir, "out.txt")
	require.NoError(t, os.WriteFile(file, []byte("a.example.com\nb.example.com\n"), 0644))

	client := &capturePostClient{}
	output := &workflow.Output{
		Data:    []string{file},
		Metrics: &workflow.Metrics{},
	}
	params := &workflow.Params{ScanID: 1, TargetID: 2}

	err := w.SaveResults(context.Background(), client, params, output)
	require.NoError(t, err)
	assert.Equal(t, 1, client.calls)
	assert.Equal(t, 2, output.Metrics.ProcessedCount)
}

func TestSaveResultsWriteSubdomainsError(t *testing.T) {
	withNopLogger(t)
	w := New(t.TempDir())
	dir := t.TempDir()
	file := filepath.Join(dir, "out.txt")
	require.NoError(t, os.WriteFile(file, []byte("a.example.com\n"), 0644))

	client := &capturePostClient{err: errors.New("post failed")}
	output := &workflow.Output{
		Data:    []string{file},
		Metrics: &workflow.Metrics{},
	}
	params := &workflow.Params{ScanID: 1, TargetID: 2}

	err := w.SaveResults(context.Background(), client, params, output)
	require.Error(t, err)
}

func TestSaveResultsParseError(t *testing.T) {
	withNopLogger(t)
	w := New(t.TempDir())
	dir := t.TempDir()
	file := filepath.Join(dir, "long.txt")
	longLine := strings.Repeat("a", 70*1024)
	require.NoError(t, os.WriteFile(file, []byte(longLine+"\n"), 0644))

	client := &capturePostClient{}
	output := &workflow.Output{
		Data:    []string{file},
		Metrics: &workflow.Metrics{},
	}
	params := &workflow.Params{ScanID: 1, TargetID: 2}

	err := w.SaveResults(context.Background(), client, params, output)
	require.Error(t, err)
}
