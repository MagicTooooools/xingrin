package pkg

import (
	"os"
	"strings"
)

// ReadVersion reads version from a file path, returning "unknown" on failure.
func ReadVersion(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return "unknown"
	}
	version := strings.TrimSpace(string(data))
	if version == "" {
		return "unknown"
	}
	return version
}
