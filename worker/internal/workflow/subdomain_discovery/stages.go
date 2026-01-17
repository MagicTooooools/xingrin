package subdomain_discovery

import (
	"context"
	"os"

	"github.com/orbit/worker/internal/activity"
	"github.com/orbit/worker/internal/server"
)

const (
	// Stage names (kebab-case to match YAML config)
	stagePassive     = "passive-tools"
	stageBruteforce  = "bruteforce"
	stagePermutation = "permutation"
	stageResolve     = "resolve"

	// Tool names (kebab-case to match templates.yaml)
	toolSubfinder            = "subfinder"
	toolSublist3r            = "sublist3r"
	toolAssetfinder          = "assetfinder"
	toolSubdomainBruteforce  = "subdomain-bruteforce"
	toolSubdomainPermutation = "subdomain-permutation-resolve"
	toolSubdomainResolve     = "subdomain-resolve"
)

var (
	// Configurable paths with defaults
	resolversPath    = envOrDefault("RESOLVERS_PATH", "/opt/orbit/wordlists/resolvers.txt")
	wordlistBasePath = envOrDefault("WORDLIST_BASE_PATH", "/opt/orbit/wordlists")
)

var passiveTools = []string{toolSubfinder, toolSublist3r, toolAssetfinder}

// envOrDefault returns environment variable value or default (package-level init)
func envOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// workflowContext holds shared context for stage execution
type workflowContext struct {
	ctx                context.Context
	domains            []string
	config             map[string]any
	workDir            string
	providerConfigPath string
	serverClient       server.ServerClient // for downloading wordlists etc.
}

// stageResult holds the output of a stage execution
type stageResult struct {
	files   []string
	failed  []string
	success []string
}

// merge combines another stageResult into this one
func (sr *stageResult) merge(other stageResult) {
	sr.files = append(sr.files, other.files...)
	sr.failed = append(sr.failed, other.failed...)
	sr.success = append(sr.success, other.success...)
}

// runAllStages executes all discovery stages and collects results
func (w *Workflow) runAllStages(ctx *workflowContext) stageResult {
	var allResults stageResult

	// Stage 1: Passive collection (always runs if configured)
	allResults.merge(w.runPassiveStage(ctx))

	// Stage 2: Bruteforce (optional)
	if isStageEnabled(ctx.config, stageBruteforce) {
		allResults.merge(w.runBruteforceStage(ctx))
	}

	// Stage 3: Permutation (optional, requires previous output)
	if isStageEnabled(ctx.config, stagePermutation) && len(allResults.files) > 0 {
		allResults.merge(w.runMergeStage(ctx, allResults.files, stagePermutation, toolSubdomainPermutation))
	}

	// Stage 4: Resolve (optional, requires previous output)
	if isStageEnabled(ctx.config, stageResolve) && len(allResults.files) > 0 {
		allResults.merge(w.runMergeStage(ctx, allResults.files, stageResolve, toolSubdomainResolve))
	}

	return allResults
}



// processResults converts activity results to stageResult
func processResults(results []*activity.Result) stageResult {
	var sr stageResult
	for _, r := range results {
		if r.Error != nil {
			sr.failed = append(sr.failed, r.Name)
		} else {
			sr.success = append(sr.success, r.Name)
			if r.OutputFile != "" {
				sr.files = append(sr.files, r.OutputFile)
			}
		}
	}
	return sr
}
