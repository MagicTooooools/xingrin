package subdomain_discovery

import (
	"fmt"
	"path/filepath"

	"github.com/orbit/worker/internal/activity"
	"github.com/orbit/worker/internal/pkg"
	"go.uber.org/zap"
)

// runPassiveStage executes all enabled passive discovery tools
func (w *Workflow) runPassiveStage(ctx *workflowContext) stageResult {
	stageConfig, ok := ctx.config[stagePassive].(map[string]any)
	if !ok {
		pkg.Logger.Debug("Passive stage not configured")
		return stageResult{}
	}

	var commands []activity.Command

	for _, domain := range ctx.domains {
		for _, toolName := range passiveTools {
			if !isToolEnabled(stageConfig, toolName) {
				continue
			}

			toolConfig, _ := stageConfig[toolName].(map[string]any)
			cmd := w.createPassiveCommand(ctx, domain, toolName, toolConfig)
			if cmd != nil {
				commands = append(commands, *cmd)
			}
		}
	}

	if len(commands) == 0 {
		pkg.Logger.Debug("No passive tools enabled")
		return stageResult{}
	}

	pkg.Logger.Info("Running passive stage", zap.Int("tools", len(commands)))
	results := w.runner.RunParallel(ctx.ctx, commands)
	return processResults(results)
}

// createPassiveCommand creates a command for a passive discovery tool
func (w *Workflow) createPassiveCommand(ctx *workflowContext, domain, toolName string, toolConfig map[string]any) *activity.Command {
	outputFile := filepath.Join(ctx.workDir, fmt.Sprintf("%s_%s.txt", toolName, sanitizeFilename(domain)))
	logFile := filepath.Join(ctx.workDir, fmt.Sprintf("%s_%s.log", toolName, sanitizeFilename(domain)))

	params := map[string]string{
		"domain":      domain,
		"output-file": outputFile,
	}

	// Add provider config for subfinder
	if toolName == toolSubfinder && ctx.providerConfigPath != "" {
		params["provider-config"] = ctx.providerConfigPath
	}

	cmdStr, err := buildCommand(toolName, params, toolConfig)
	if err != nil {
		pkg.Logger.Error("Failed to build command",
			zap.String("tool", toolName),
			zap.Error(err))
		return nil
	}

	return &activity.Command{
		Name:       fmt.Sprintf("%s_%s", toolName, sanitizeFilename(domain)),
		Command:    cmdStr,
		OutputFile: outputFile,
		LogFile:    logFile,
		Timeout:    getTimeout(toolConfig),
	}
}
