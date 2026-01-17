package activity

import (
	"fmt"
	"strings"
)

// CommandBuilder builds commands from templates
type CommandBuilder struct{}

// NewCommandBuilder creates a new command builder
func NewCommandBuilder() *CommandBuilder {
	return &CommandBuilder{}
}

// Build constructs a command from a template with the given parameters
func (b *CommandBuilder) Build(tmpl CommandTemplate, params map[string]string, config map[string]any) (string, error) {
	if tmpl.Base == "" {
		return "", fmt.Errorf("template base command is empty")
	}

	// Start with base command
	cmd := tmpl.Base

	// Replace required placeholders
	for key, value := range params {
		placeholder := "{" + key + "}"
		cmd = strings.ReplaceAll(cmd, placeholder, value)
	}

	// Append optional parameters if present in config
	for configKey, flagTemplate := range tmpl.Optional {
		if value, ok := getConfigValue(config, configKey); ok {
			flag := strings.ReplaceAll(flagTemplate, "{"+configKey+"}", fmt.Sprintf("%v", value))
			cmd = cmd + " " + flag
		}
	}

	// Check for unreplaced placeholders (indicates missing required params)
	if strings.Contains(cmd, "{") && strings.Contains(cmd, "}") {
		return "", fmt.Errorf("command contains unreplaced placeholders: %s", cmd)
	}

	return cmd, nil
}

func getConfigValue(config map[string]any, key string) (any, bool) {
	if config == nil {
		return nil, false
	}
	value, ok := config[key]
	return value, ok
}
