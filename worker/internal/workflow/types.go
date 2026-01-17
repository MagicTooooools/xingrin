package workflow

import (
	"context"
	"time"

	"github.com/orbit/worker/internal/server"
)

// Params contains parameters for a workflow execution
type Params struct {
	ScanID       int
	TargetID     int
	TargetName   string
	TargetType   string
	WorkDir      string
	ScanConfig   map[string]any
	ServerClient server.ServerClient
}

// Output contains the result of a workflow execution
type Output struct {
	Data    any      // Workflow-specific data (file paths for streaming, or parsed data)
	Metrics *Metrics // Execution statistics (optional)
}

// Metrics contains execution statistics
type Metrics struct {
	ProcessedCount int
	FailedCount    int
	FailedTools    []string
	Duration       time.Duration
}

// Workflow defines the interface for scan workflows
type Workflow interface {
	Name() string
	Execute(params *Params) (*Output, error)
	SaveResults(ctx context.Context, client server.ServerClient, params *Params, output *Output) error
}
