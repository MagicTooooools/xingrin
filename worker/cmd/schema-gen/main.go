package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"

	"github.com/orbit/worker/internal/activity"
	"github.com/orbit/worker/internal/workflow"
	"gopkg.in/yaml.v3"
)

// JSONSchema represents a JSON Schema Draft 7 structure
type JSONSchema struct {
	Schema               string                     `json:"$schema"`
	Title                string                     `json:"title,omitempty"`
	Description          string                     `json:"description,omitempty"`
	Type                 string                     `json:"type"`
	Properties           map[string]*PropertySchema `json:"properties,omitempty"`
	Required             []string                   `json:"required,omitempty"`
	AdditionalProperties bool                       `json:"additionalProperties"`
	Metadata             map[string]interface{}     `json:"x-metadata,omitempty"`
}

// PropertySchema represents a property in JSON Schema
type PropertySchema struct {
	Type        string                     `json:"type,omitempty"`
	Description string                     `json:"description,omitempty"`
	Default     interface{}                `json:"default,omitempty"`
	Properties  map[string]*PropertySchema `json:"properties,omitempty"`
	Required    []string                   `json:"required,omitempty"`
	// Tool metadata extension fields
	Stage   string `json:"x-stage,omitempty"`
	Warning string `json:"x-warning,omitempty"`
}

// TemplateFile represents the structure of templates.yaml
type TemplateFile struct {
	Metadata workflow.WorkflowMetadata           `yaml:"metadata"`
	Tools    map[string]activity.CommandTemplate `yaml:"tools"`
}

func main() {
	inputFile := flag.String("input", "", "Input templates.yaml file")
	outputFile := flag.String("output", "", "Output JSON Schema file")
	flag.Parse()

	if *inputFile == "" || *outputFile == "" {
		fmt.Fprintf(os.Stderr, "Usage: schema-gen -input <templates.yaml> -output <schema.json>\n")
		os.Exit(1)
	}

	// Read template file
	data, err := os.ReadFile(*inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input file: %v\n", err)
		os.Exit(1)
	}

	// Parse YAML
	var templateFile TemplateFile
	if err := yaml.Unmarshal(data, &templateFile); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing YAML: %v\n", err)
		os.Exit(1)
	}

	// Generate JSON Schema
	schema := generateSchema(templateFile)

	// Output JSON
	output, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating JSON: %v\n", err)
		os.Exit(1)
	}

	// Write to file
	if err := os.WriteFile(*outputFile, output, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing output file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Schema generated successfully: %s\n", *outputFile)
}

func generateSchema(templateFile TemplateFile) *JSONSchema {
	schema := &JSONSchema{
		Schema:               "http://json-schema.org/draft-07/schema#",
		Title:                templateFile.Metadata.DisplayName,
		Description:          templateFile.Metadata.Description,
		Type:                 "object",
		Properties:           make(map[string]*PropertySchema),
		AdditionalProperties: false,
		Metadata: map[string]interface{}{
			"name":         templateFile.Metadata.Name,
			"version":      templateFile.Metadata.Version,
			"target_types": templateFile.Metadata.TargetTypes,
			"stages":       templateFile.Metadata.Stages,
		},
	}
	// Group tools by stage
	toolsByStage := make(map[string][]string)
	for toolName, tool := range templateFile.Tools {
		stage := tool.Metadata.Stage
		toolsByStage[stage] = append(toolsByStage[stage], toolName)
	}

	// Generate schema for each stage
	var requiredStages []string
	for _, stage := range templateFile.Metadata.Stages {
		stageSchema := &PropertySchema{
			Type:       "object",
			Properties: make(map[string]*PropertySchema),
		}

		stageSchema.Properties["enabled"] = &PropertySchema{
			Type:        "boolean",
			Description: "Whether to enable this stage",
		}

		toolsSchema := &PropertySchema{
			Type:       "object",
			Properties: make(map[string]*PropertySchema),
		}

		tools := toolsByStage[stage.ID]
		if len(tools) > 0 {
			sort.Strings(tools)
		}

		var requiredTools []string
		for _, toolName := range tools {
			tool := templateFile.Tools[toolName]
			toolSchema := &PropertySchema{
				Type:        "object",
				Description: tool.Metadata.Description,
				Properties:  make(map[string]*PropertySchema),
				Stage:       tool.Metadata.Stage,
				Warning:     tool.Metadata.Warning,
			}

			toolSchema.Properties["enabled"] = &PropertySchema{
				Type:        "boolean",
				Description: "Whether to enable this tool",
			}

			var requiredParams []string
			for _, param := range append(append([]activity.Parameter{}, tool.RuntimeParams...), tool.CLIParams...) {
				paramSchema := &PropertySchema{
					Description: param.Documentation.Description,
				}

				switch param.ConfigSchema.Type {
				case "string":
					paramSchema.Type = "string"
				case "integer":
					paramSchema.Type = "integer"
				case "boolean":
					paramSchema.Type = "boolean"
				default:
					paramSchema.Type = "string"
				}
				toolSchema.Properties[param.ConfigSchema.Key] = paramSchema
				if param.ConfigSchema.Required {
					requiredParams = append(requiredParams, param.ConfigSchema.Key)
				}
			}

			if len(requiredParams) > 0 {
				toolSchema.Required = requiredParams
			}

			toolsSchema.Properties[toolName] = toolSchema
			requiredTools = append(requiredTools, toolName)
		}

		if len(requiredTools) > 0 {
			toolsSchema.Required = requiredTools
		}

		stageSchema.Properties["tools"] = toolsSchema
		stageSchema.Required = []string{"enabled", "tools"}

		schema.Properties[stage.ID] = stageSchema
		requiredStages = append(requiredStages, stage.ID)
	}

	if len(requiredStages) > 0 {
		schema.Required = requiredStages
	}

	return schema
}
