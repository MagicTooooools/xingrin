package subdomain_discovery

import (
	"bufio"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/orbit/worker/internal/activity"
)

const (
	defaultTimeout = 86400 // default max timeout: 24 hours
)

// buildCommand gets the template and builds the command string
func buildCommand(toolName string, params map[string]string, config map[string]any) (string, error) {
	tmpl, err := getTemplate(toolName)
	if err != nil {
		return "", err
	}
	builder := activity.NewCommandBuilder()
	return builder.Build(tmpl, params, config)
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
	toolConfig, ok := stageConfig[toolName].(map[string]any)
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
// If not configured, returns default (24 hours)
func getTimeout(toolConfig map[string]any) time.Duration {
	if toolConfig == nil {
		return time.Duration(defaultTimeout) * time.Second
	}

	if timeout, ok := toolConfig["timeout"].(int); ok && timeout > 0 {
		return time.Duration(timeout) * time.Second
	}
	if timeout, ok := toolConfig["timeout"].(float64); ok && timeout > 0 {
		return time.Duration(timeout) * time.Second
	}

	return time.Duration(defaultTimeout) * time.Second
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
