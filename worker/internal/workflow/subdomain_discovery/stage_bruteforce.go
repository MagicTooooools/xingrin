package subdomain_discovery

import (
	"fmt"
	"path/filepath"

	"github.com/orbit/worker/internal/activity"
	"github.com/orbit/worker/internal/pkg"
	"go.uber.org/zap"
)

// runBruteforceStage executes subdomain bruteforce for all domains
func (w *Workflow) runBruteforceStage(ctx *workflowContext) stageResult {
	stageConfig, ok := ctx.config[stageBruteforce].(map[string]any)
	if !ok {
		pkg.Logger.Debug("Bruteforce stage not configured")
		return stageResult{}
	}

	// Get tool-specific config
	toolConfig, _ := stageConfig[toolSubdomainBruteforce].(map[string]any)

	// Get wordlist name from config (required, no default)
	wordlistName := getStringValue(toolConfig, "wordlist-name", "")
	if wordlistName == "" {
		pkg.Logger.Error("Bruteforce stage requires wordlist-name in config")
		return stageResult{failed: []string{stageBruteforce + " (missing wordlist-name in config)"}}
	}

	// Ensure wordlist exists locally (download from server if needed)
	wordlistPath, err := ctx.serverClient.EnsureWordlistLocal(ctx.ctx, wordlistName, wordlistBasePath)
	if err != nil {
		pkg.Logger.Error("Failed to get wordlist",
			zap.String("wordlist", wordlistName),
			zap.Error(err))
		return stageResult{failed: []string{stageBruteforce + " (wordlist: " + err.Error() + ")"}}
	}

	var commands []activity.Command

	for _, domain := range ctx.domains {
		cmd := w.createBruteforceCommand(ctx, domain, toolConfig, wordlistPath)
		if cmd != nil {
			commands = append(commands, *cmd)
		}
	}

	if len(commands) == 0 {
		pkg.Logger.Debug("No bruteforce commands created")
		return stageResult{}
	}

	pkg.Logger.Info("Running bruteforce stage",
		zap.Int("domains", len(commands)),
		zap.String("wordlist", wordlistPath))
	results := w.runner.RunParallel(ctx.ctx, commands)
	return processResults(results)
}

// createBruteforceCommand creates a bruteforce command for a domain
func (w *Workflow) createBruteforceCommand(ctx *workflowContext, domain string, toolConfig map[string]any, wordlistPath string) *activity.Command {
	outputFile := filepath.Join(ctx.workDir, fmt.Sprintf("bruteforce_%s.txt", sanitizeFilename(domain)))
	logFile := filepath.Join(ctx.workDir, fmt.Sprintf("bruteforce_%s.log", sanitizeFilename(domain)))

	params := map[string]string{
		"domain":      domain,
		"output-file": outputFile,
		"wordlist":    wordlistPath,
		"resolvers":   resolversPath,
	}

	cmdStr, err := buildCommand(toolSubdomainBruteforce, params, toolConfig)
	if err != nil {
		pkg.Logger.Error("Failed to build bruteforce command",
			zap.String("domain", domain),
			zap.Error(err))
		return nil
	}

	timeout := getTimeout(toolConfig)

	return &activity.Command{
		Name:       fmt.Sprintf("bruteforce_%s", sanitizeFilename(domain)),
		Command:    cmdStr,
		OutputFile: outputFile,
		LogFile:    logFile,
		Timeout:    timeout,
	}
}
