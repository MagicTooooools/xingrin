package subdomain_discovery

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/orbit/worker/internal/activity"
	"github.com/orbit/worker/internal/pkg"
	"github.com/orbit/worker/internal/pkg/validator"
	"github.com/orbit/worker/internal/server"
	"github.com/orbit/worker/internal/workflow"
	"go.uber.org/zap"
)

const Name = "subdomain_discovery"

func init() {
	workflow.Register(Name, func(workDir string) workflow.Workflow {
		return New(workDir)
	})
}

// Workflow implements the subdomain discovery scan workflow
type Workflow struct {
	runner  *activity.Runner
	workDir string
}

// New creates a new subdomain discovery workflow
func New(workDir string) *Workflow {
	return &Workflow{
		runner:  activity.NewRunner(workDir),
		workDir: workDir,
	}
}

func (w *Workflow) Name() string {
	return Name
}

// Execute runs the subdomain discovery workflow
func (w *Workflow) Execute(params *workflow.Params) (*workflow.Output, error) {
	// Initialize and validate
	ctx, err := w.initialize(params)
	if err != nil {
		return nil, err
	}

	// Run all stages
	allResults := w.runAllStages(ctx)

	// Store result files for streaming in SaveResults
	output := &workflow.Output{
		Data: allResults.files, // Pass file paths instead of parsed data
		Metrics: &workflow.Metrics{
			ProcessedCount: 0, // Will be updated after streaming
			FailedCount:    len(allResults.failed),
			FailedTools:    allResults.failed,
		},
	}

	// Check for complete failure
	if len(allResults.failed) > 0 && len(allResults.success) == 0 {
		return output, fmt.Errorf("all tools failed")
	}

	return output, nil
}

// SaveResults streams subdomain results to the server in batches
func (w *Workflow) SaveResults(ctx context.Context, client server.ServerClient, params *workflow.Params, output *workflow.Output) error {
	files, ok := output.Data.([]string)
	if !ok || len(files) == 0 {
		return nil
	}

	// Create batch sender with context
	sender := server.NewBatchSender(ctx, client, params.ScanID, params.TargetID, "subdomain", 5000)

	// Stream and deduplicate from files
	subdomainCh, errCh := w.streamAndDeduplicate(files)

	// Send subdomains in batches
	for subdomain := range subdomainCh {
		if err := sender.Add(map[string]string{"name": subdomain}); err != nil {
			return err
		}
	}

	// Check for streaming errors
	select {
	case err := <-errCh:
		if err != nil {
			return fmt.Errorf("error streaming results: %w", err)
		}
	default:
	}

	// Flush remaining items
	if err := sender.Flush(); err != nil {
		return err
	}

	// Update metrics
	items, batches := sender.Stats()
	output.Metrics.ProcessedCount = items
	pkg.Logger.Info("Results saved",
		zap.Int("subdomains", items),
		zap.Int("batches", batches))

	return nil
}

// initialize validates params and prepares the workflow context
func (w *Workflow) initialize(params *workflow.Params) (*workflowContext, error) {
	// Config can be either nested under workflow name or flat
	// Try nested first: { "subdomain_discovery": { "passive-tools": ... } }
	// Then flat: { "passive-tools": ... }
	flowConfig := getConfigPath(params.ScanConfig, Name)
	if flowConfig == nil {
		// Use flat config directly
		flowConfig = params.ScanConfig
	}
	if flowConfig == nil {
		return nil, fmt.Errorf("missing %s config", Name)
	}

	workDir := filepath.Join(params.WorkDir, Name)
	if err := os.MkdirAll(workDir, 0755); err != nil {
		return nil, err
	}

	// Subdomain discovery only works for domain type targets
	if params.TargetType != "domain" {
		return nil, fmt.Errorf("subdomain discovery requires domain target, got %s", params.TargetType)
	}

	// Normalize domain first
	normalizedDomain, err := validator.NormalizeDomain(params.TargetName)
	if err != nil {
		return nil, fmt.Errorf("failed to normalize domain: %w", err)
	}

	// Validate normalized domain
	if err := validator.ValidateDomain(normalizedDomain); err != nil {
		return nil, fmt.Errorf("invalid target domain: %w", err)
	}

	// Wrap in slice for compatibility with multi-domain processing
	domains := []string{normalizedDomain}

	pkg.Logger.Info("Workflow initialized",
		zap.Int("scanId", params.ScanID),
		zap.String("targetName", params.TargetName),
		zap.String("targetType", params.TargetType))

	ctx := context.Background()
	providerConfigPath, err := w.setupProviderConfig(ctx, params, workDir)
	if err != nil {
		// Log warning but continue - provider config is optional (enhances results but not required)
		pkg.Logger.Warn("Failed to setup provider config, subfinder will run without API keys",
			zap.Error(err))
	}

	return &workflowContext{
		ctx:                ctx,
		domains:            domains,
		config:             flowConfig,
		workDir:            workDir,
		providerConfigPath: providerConfigPath,
		serverClient:       params.ServerClient,
	}, nil
}

// setupProviderConfig fetches and writes the subfinder provider config
// Returns empty string if no config available, error if fetch/write failed
func (w *Workflow) setupProviderConfig(ctx context.Context, params *workflow.Params, workDir string) (string, error) {
	providerConfig, err := params.ServerClient.GetProviderConfig(ctx, params.ScanID, toolSubfinder)
	if err != nil {
		return "", fmt.Errorf("failed to get provider config: %w", err)
	}
	if providerConfig == nil || providerConfig.Content == "" {
		return "", nil // No config available, not an error
	}

	configPath := filepath.Join(workDir, "provider-config.yaml")
	if err := os.WriteFile(configPath, []byte(providerConfig.Content), 0600); err != nil {
		return "", fmt.Errorf("failed to write provider config: %w", err)
	}
	pkg.Logger.Info("Provider config written", zap.String("path", configPath))
	return configPath, nil
}

// streamAndDeduplicate streams unique subdomains from multiple files.
// Returns a channel that yields deduplicated subdomains one by one.
// The channel is closed when all files are processed or an error occurs.
func (w *Workflow) streamAndDeduplicate(filePaths []string) (<-chan string, <-chan error) {
	out := make(chan string, 1000) // buffered for better throughput
	errCh := make(chan error, 1)
	seen := make(map[string]struct{}, 500000) // pre-allocate for large datasets

	go func() {
		defer close(out)
		defer close(errCh)

		// Recover from panic and send error to channel
		defer func() {
			if r := recover(); r != nil {
				pkg.Logger.Error("Panic in streamAndDeduplicate", zap.Any("panic", r))
				errCh <- fmt.Errorf("panic in stream processing: %v", r)
			}
		}()

		for _, path := range filePaths {
			if err := w.streamFile(path, seen, out); err != nil {
				pkg.Logger.Error("Error streaming file", zap.String("path", path), zap.Error(err))
				errCh <- err
				return
			}
		}
	}()

	return out, errCh
}

// streamFile reads a single file and sends unique subdomains to the channel
// Note: Files are already validated by mergeFiles, so we skip validation for performance
func (w *Workflow) streamFile(filePath string, seen map[string]struct{}, out chan<- string) error {
	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer func() { _ = file.Close() }()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		lower := strings.ToLower(line)
		if _, exists := seen[lower]; !exists {
			seen[lower] = struct{}{}
			out <- line
		}
	}

	return scanner.Err()
}
