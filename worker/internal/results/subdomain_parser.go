package results

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/yyhuni/lunafox/worker/internal/pkg"
	"go.uber.org/zap"
)

// Subdomain represents a single subdomain result.
type Subdomain struct {
	Name string `json:"name"`
}

// ParseSubdomains streams and deduplicates subdomains from multiple files.
// Deduplication is case-insensitive; first occurrence is preserved.
func ParseSubdomains(filePaths []string) (<-chan Subdomain, <-chan error) {
	out := make(chan Subdomain, 1000)
	errCh := make(chan error, 1)
	seen := make(map[string]struct{}, 500000)

	go func() {
		defer close(out)
		defer close(errCh)

		defer func() {
			if r := recover(); r != nil {
				pkg.Logger.Error("Panic in ParseSubdomains", zap.Any("panic", r))
				errCh <- fmt.Errorf("panic in subdomain parsing: %v", r)
			}
		}()

		for _, path := range filePaths {
			if err := streamSubdomainFile(path, seen, out); err != nil {
				pkg.Logger.Error("Error streaming subdomain file", zap.String("path", path), zap.Error(err))
				errCh <- err
				return
			}
		}
	}()

	return out, errCh
}

// streamSubdomainFile reads a single file and sends unique subdomains to the channel.
func streamSubdomainFile(filePath string, seen map[string]struct{}, out chan<- Subdomain) error {
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
			out <- Subdomain{Name: line}
		}
	}

	return scanner.Err()
}
