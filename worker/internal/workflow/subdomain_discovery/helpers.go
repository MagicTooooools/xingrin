package subdomain_discovery

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/orbit/worker/internal/activity"
)

// normalizeToolConfig maps external config_schema.key values to internal parameter names.
func normalizeToolConfig(toolName string, config map[string]any) (map[string]any, error) {
	tmpl, err := getTemplate(toolName)
	if err != nil {
		return nil, err
	}
	return activity.MapConfigKeys(tmpl, config)
}

// buildCommand gets the template and builds the command string.
// config must be normalized (internal parameter names).
func buildCommand(toolName string, params map[string]any, config map[string]any) (string, error) {
	tmpl, err := getTemplate(toolName)
	if err != nil {
		return "", err
	}
	builder := activity.NewCommandBuilder()
	return builder.Build(tmpl, params, config)
}

// validateExplicitConfig ensures stage/tool enabled flags are explicitly set in config.
func validateExplicitConfig(config map[string]any) error {
	if config == nil {
		return fmt.Errorf("config is required")
	}

	metadata, err := loader.GetMetadata()
	if err != nil {
		return err
	}
	templates, err := loader.Load()
	if err != nil {
		return err
	}

	toolsByStage := make(map[string][]string)
	for toolName, tmpl := range templates {
		stage := tmpl.Metadata.Stage
		toolsByStage[stage] = append(toolsByStage[stage], toolName)
	}

	for _, stage := range metadata.Stages {
		stageConfigRaw, ok := config[stage.ID]
		if !ok {
			return fmt.Errorf("stage %s config is required", stage.ID)
		}
		stageConfig, ok := stageConfigRaw.(map[string]any)
		if !ok {
			return fmt.Errorf("stage %s config must be map", stage.ID)
		}
		enabledRaw, ok := stageConfig["enabled"]
		if !ok {
			return fmt.Errorf("stage %s.enabled is required", stage.ID)
		}
		if _, ok := enabledRaw.(bool); !ok {
			return fmt.Errorf("stage %s.enabled must be boolean", stage.ID)
		}

		toolsConfigRaw, ok := stageConfig["tools"]
		if !ok {
			return fmt.Errorf("stage %s.tools is required", stage.ID)
		}
		toolsConfig, ok := toolsConfigRaw.(map[string]any)
		if !ok {
			return fmt.Errorf("stage %s.tools must be map", stage.ID)
		}

		for _, toolName := range toolsByStage[stage.ID] {
			toolConfigRaw, ok := toolsConfig[toolName]
			if !ok {
				return fmt.Errorf("tool %s config is required", toolName)
			}
			toolConfig, ok := toolConfigRaw.(map[string]any)
			if !ok {
				return fmt.Errorf("tool %s config must be map", toolName)
			}
			enabledRaw, ok := toolConfig["enabled"]
			if !ok {
				return fmt.Errorf("tool %s.enabled is required", toolName)
			}
			if _, ok := enabledRaw.(bool); !ok {
				return fmt.Errorf("tool %s.enabled must be boolean", toolName)
			}
		}
	}

	return nil
}

// isStageEnabled checks if a stage is enabled in the config
func isStageEnabled(config map[string]any, stageName string) bool {
	stageConfig, ok := config[stageName].(map[string]any)
	if !ok {
		return false
	}
	enabled, ok := stageConfig["enabled"].(bool)
	return ok && enabled
}

// isToolEnabled checks if a specific tool is enabled within a stage
func isToolEnabled(stageConfig map[string]any, toolName string) bool {
	toolsConfig, ok := stageConfig["tools"].(map[string]any)
	if !ok {
		return false
	}
	toolConfig, ok := toolsConfig[toolName].(map[string]any)
	if !ok {
		return false
	}
	enabled, ok := toolConfig["enabled"].(bool)
	return ok && enabled
}

// getConfigPath retrieves a nested config section by path
func getConfigPath(config map[string]any, path string) map[string]any {
	if config == nil {
		return nil
	}
	parts := strings.Split(path, ".")
	current := config
	for _, part := range parts {
		next, ok := current[part].(map[string]any)
		if !ok {
			return nil
		}
		current = next
	}
	return current
}

// getTimeout extracts timeout from tool config
// Requires an explicit timeout value.
func getTimeout(toolConfig map[string]any) (time.Duration, error) {
	seconds, err := getIntValue(toolConfig, "timeout-runtime")
	if err != nil {
		return 0, fmt.Errorf("timeout: %w", err)
	}
	if seconds <= 0 {
		return 0, fmt.Errorf("timeout must be > 0")
	}
	return time.Duration(seconds) * time.Second, nil
}

// getIntValue extracts an integer value from config.
func getIntValue(config map[string]any, key string) (int, error) {
	if config == nil {
		return 0, fmt.Errorf("%s is required", key)
	}
	raw, ok := config[key]
	if !ok {
		return 0, fmt.Errorf("%s is required", key)
	}
	switch v := raw.(type) {
	case int:
		return v, nil
	case int8:
		return int(v), nil
	case int16:
		return int(v), nil
	case int32:
		return int(v), nil
	case int64:
		return int(v), nil
	case uint:
		return int(v), nil
	case uint8:
		return int(v), nil
	case uint16:
		return int(v), nil
	case uint32:
		return int(v), nil
	case uint64:
		return int(v), nil
	case float64:
		return int(v), nil
	default:
		return 0, fmt.Errorf("%s must be integer", key)
	}
}

// countFileLines counts non-empty lines in a file
func countFileLines(filePath string) int {
	file, err := os.Open(filePath)
	if err != nil {
		return 0
	}
	defer func() { _ = file.Close() }()

	count := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.TrimSpace(scanner.Text()) != "" {
			count++
		}
	}
	return count
}

// getStringValue extracts a string value from config with a default
func getStringValue(config map[string]any, key, defaultValue string) string {
	if config == nil {
		return defaultValue
	}
	if value, ok := config[key].(string); ok && value != "" {
		return value
	}
	return defaultValue
}

// sanitizeFilename removes or replaces characters that are invalid in filenames
func sanitizeFilename(name string) string {
	// Replace common problematic characters
	re := regexp.MustCompile(`[<>:"/\\|?*\s]`)
	return re.ReplaceAllString(name, "_")
}
