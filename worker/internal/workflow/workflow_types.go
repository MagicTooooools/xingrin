package workflow

import (
	"context"

	"github.com/yyhuni/lunafox/worker/internal/server"
)

// Workflow defines the workflow interface
type Workflow interface {
	// Execute executes the workflow
	Execute(params *Params) (*Output, error)
	// Name returns the workflow name
	Name() string
	// SaveResults saves results to the server
	SaveResults(ctx context.Context, client server.ServerClient, params *Params, output *Output) error
}

// Params defines workflow execution parameters
type Params struct {
	ScanID       int
	TargetID     int
	TargetName   string
	TargetType   string
	WorkDir      string
	ScanConfig   map[string]any
	Config       map[string]any
	ServerClient server.ServerClient
}

// Output defines workflow output
type Output struct {
	Data    interface{}
	Metrics *Metrics
}

// Metrics defines workflow execution metrics
type Metrics struct {
	ProcessedCount int
	FailedCount    int
	FailedTools    []string
}

// WorkflowNotes defines workflow-level notes by category.
type WorkflowNotes struct {
	Global []string `yaml:"global,omitempty" json:"global,omitempty"`
	Config []string `yaml:"config,omitempty" json:"config,omitempty"`
}

// DocSection defines a documentation section header.
type DocSection struct {
	ID    string `yaml:"id" json:"id"`
	Title string `yaml:"title" json:"title"`
}

// DocMetadata defines documentation structure and labels.
type DocMetadata struct {
	Sections      []DocSection      `yaml:"sections,omitempty" json:"sections,omitempty"`
	Labels        map[string]string `yaml:"labels,omitempty" json:"labels,omitempty"`
	ExampleHeader []string          `yaml:"example_header,omitempty" json:"exampleHeader,omitempty"`
}

// WorkflowMetadata defines workflow metadata structure
type WorkflowMetadata struct {
	Name        string          `yaml:"name" json:"name"`
	DisplayName string          `yaml:"display_name" json:"displayName"`
	Description string          `yaml:"description" json:"description"`
	Version     string          `yaml:"version" json:"version"`
	TargetTypes []string        `yaml:"target_types" json:"targetTypes"`
	Notes       WorkflowNotes   `yaml:"notes,omitempty" json:"notes,omitempty"`
	Doc         *DocMetadata    `yaml:"doc,omitempty" json:"doc,omitempty"`
	Stages      []StageMetadata `yaml:"stages" json:"stages"`
}

// StageMetadata defines stage metadata structure
type StageMetadata struct {
	ID          string   `yaml:"id" json:"id"`
	Name        string   `yaml:"name" json:"name"`
	Description string   `yaml:"description" json:"description"`
	Required    bool     `yaml:"required" json:"required"`
	Parallel    bool     `yaml:"parallel" json:"parallel"`
	Notes       []string `yaml:"notes,omitempty" json:"notes,omitempty"`
}
