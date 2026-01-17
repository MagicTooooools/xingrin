package subdomain_discovery

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/orbit/worker/internal/activity"
	"github.com/orbit/worker/internal/pkg"
	"github.com/orbit/worker/internal/pkg/validator"
	"go.uber.org/zap"
)

const (
	wildcardSampleTimeout = 2 * time.Hour // 2 hours for sampling
	wildcardTests         = 50             // puredns wildcard tests count
	wildcardBatch         = 1000000        // puredns wildcard batch size
	sampleMultiplier      = 100            // sample size = original count × 100
	expansionThreshold    = 50             // threshold = original count × 50
)

// wildcardCheckResult holds the result of wildcard detection
type wildcardCheckResult struct {
	isWildcard     bool
	originalCount  int
	sampleCount    int
	expansionRatio float64
	reason         string
}

// runMergeStage merges input files and runs a processing tool (resolve or permutation)
func (w *Workflow) runMergeStage(ctx *workflowContext, inputFiles []string, stageName, toolName string) stageResult {
	stageConfig, ok := ctx.config[stageName].(map[string]any)
	if !ok {
		pkg.Logger.Debug("Stage not configured", zap.String("stage", stageName))
		return stageResult{}
	}

	// Get tool-specific config
	toolConfig, _ := stageConfig[toolName].(map[string]any)

	// Merge all input files into one
	mergedFile := filepath.Join(ctx.workDir, fmt.Sprintf("%s_input.txt", stageName))
	if err := w.mergeFiles(inputFiles, mergedFile); err != nil {
		pkg.Logger.Error("Failed to merge files",
			zap.String("stage", stageName),
			zap.Error(err))
		return stageResult{failed: []string{stageName}}
	}

	// For permutation stage, check for wildcard domains first
	if stageName == stagePermutation {
		checkResult := w.checkWildcard(ctx.ctx, mergedFile, ctx.workDir)
		if checkResult.isWildcard {
			pkg.Logger.Warn("Skipping permutation stage due to wildcard detection",
				zap.String("reason", checkResult.reason),
				zap.Float64("expansionRatio", checkResult.expansionRatio))
			return stageResult{
				failed: []string{fmt.Sprintf("%s (wildcard: %s)", stageName, checkResult.reason)},
			}
		}
		pkg.Logger.Info("Wildcard check passed, proceeding with permutation",
			zap.Int("originalCount", checkResult.originalCount),
			zap.Int("sampleCount", checkResult.sampleCount))
	}

	outputFile := filepath.Join(ctx.workDir, fmt.Sprintf("%s_output.txt", stageName))
	logFile := filepath.Join(ctx.workDir, fmt.Sprintf("%s.log", stageName))

	params := map[string]string{
		"input-file":  mergedFile,
		"output-file": outputFile,
		"resolvers":   resolversPath,
	}

	cmdStr, err := buildCommand(toolName, params, toolConfig)
	if err != nil {
		pkg.Logger.Error("Failed to build command",
			zap.String("stage", stageName),
			zap.String("tool", toolName),
			zap.Error(err))
		return stageResult{failed: []string{stageName}}
	}

	timeout := getTimeout(toolConfig)

	cmd := activity.Command{
		Name:       stageName,
		Command:    cmdStr,
		OutputFile: outputFile,
		LogFile:    logFile,
		Timeout:    timeout,
	}

	pkg.Logger.Info("Running merge stage",
		zap.String("stage", stageName),
		zap.Int("inputFiles", len(inputFiles)),
		zap.Duration("timeout", timeout))

	results := w.runner.RunParallel(ctx.ctx, []activity.Command{cmd})
	return processResults(results)
}

// mergeFiles reads all input files, deduplicates entries, and writes to outputFile
// Uses streaming to minimize memory usage
func (w *Workflow) mergeFiles(inputFiles []string, outputFile string) error {
	out, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer func() { _ = out.Close() }()

	seen := make(map[string]struct{}, 100000) // pre-allocate for better performance
	writer := bufio.NewWriter(out)
	defer func() { _ = writer.Flush() }()

	for _, f := range inputFiles {
		if err := w.streamMergeFile(f, seen, writer); err != nil {
			pkg.Logger.Debug("Failed to process file during merge",
				zap.String("file", f),
				zap.Error(err))
			continue
		}
	}

	return writer.Flush()
}

// streamMergeFile streams a single file and writes unique subdomains to the writer
func (w *Workflow) streamMergeFile(filePath string, seen map[string]struct{}, writer *bufio.Writer) error {
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
		if line == "" || !validator.IsValidSubdomainFormat(line) {
			continue
		}

		lower := strings.ToLower(line)
		if _, exists := seen[lower]; !exists {
			seen[lower] = struct{}{}
			if _, err := fmt.Fprintln(writer, line); err != nil {
				return err
			}
		}
	}

	return scanner.Err()
}

// checkWildcard performs wildcard detection before permutation stage
// Returns true if wildcard is detected (should skip permutation)
func (w *Workflow) checkWildcard(ctx context.Context, inputFile, workDir string) wildcardCheckResult {
	originalCount := countFileLines(inputFile)
	if originalCount == 0 {
		return wildcardCheckResult{
			isWildcard: false,
			reason:     "empty input file",
		}
	}

	sampleSize := originalCount * sampleMultiplier
	maxAllowed := originalCount * expansionThreshold
	sampleOutput := filepath.Join(workDir, "wildcard_sample.txt")
	logFile := filepath.Join(workDir, "wildcard_detection.log")

	// Build sampling command: dnsgen | head | puredns resolve
	sampleCmd := fmt.Sprintf(
		"cat '%s' | dnsgen - | head -n %d | puredns resolve -r '%s' --write '%s' --wildcard-tests %d --wildcard-batch %d --quiet",
		inputFile, sampleSize, resolversPath, sampleOutput, wildcardTests, wildcardBatch,
	)

	pkg.Logger.Info("Wildcard detection: sampling",
		zap.Int("originalCount", originalCount),
		zap.Int("sampleSize", sampleSize),
		zap.Int("threshold", maxAllowed))

	// Run sampling with runner
	result := w.runner.Run(ctx, activity.Command{
		Name:       "wildcard_detection",
		Command:    sampleCmd,
		OutputFile: sampleOutput,
		LogFile:    logFile,
		Timeout:    wildcardSampleTimeout,
	})

	// Handle execution errors
	if result.Error != nil {
		if result.ExitCode == activity.ExitCodeTimeout {
			pkg.Logger.Warn("Wildcard detection timeout")
			return wildcardCheckResult{
				isWildcard:    true,
				originalCount: originalCount,
				reason:        "sampling timeout",
			}
		}
		// Non-timeout error, continue anyway
		pkg.Logger.Debug("Wildcard sampling command error (continuing)", zap.Error(result.Error))
	}

	// Count sample results
	sampleCount := countFileLines(sampleOutput)

	pkg.Logger.Info("Wildcard detection: sample result",
		zap.Int("sampleCount", sampleCount),
		zap.Int("originalCount", originalCount),
		zap.Int("threshold", maxAllowed))

	// Check if expansion ratio exceeds threshold
	if sampleCount > maxAllowed {
		ratio := float64(sampleCount) / float64(originalCount)
		pkg.Logger.Warn("Wildcard detected: expansion ratio too high",
			zap.Int("sampleCount", sampleCount),
			zap.Int("threshold", maxAllowed),
			zap.Float64("ratio", ratio))

		return wildcardCheckResult{
			isWildcard:     true,
			originalCount:  originalCount,
			sampleCount:    sampleCount,
			expansionRatio: ratio,
			reason:         fmt.Sprintf("expansion ratio %.1fx exceeds threshold %dx", ratio, expansionThreshold),
		}
	}

	return wildcardCheckResult{
		isWildcard:     false,
		originalCount:  originalCount,
		sampleCount:    sampleCount,
		expansionRatio: float64(sampleCount) / float64(originalCount),
	}
}
