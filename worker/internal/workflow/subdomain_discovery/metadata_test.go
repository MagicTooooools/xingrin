package subdomain_discovery

import (
	"testing"

	"github.com/orbit/worker/internal/pkg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// init 初始化测试环境
func init() {
	// 为测试环境初始化 logger（如果还没有初始化）
	if pkg.Logger == nil {
		logger, _ := zap.NewDevelopment()
		pkg.Logger = logger
	}
}

// TestStagesMatchMetadata 验证代码中的 stage 常量和 templates.yaml 中的定义一致
func TestStagesMatchMetadata(t *testing.T) {
	// 从 templates.yaml 加载 metadata
	metadata, err := loader.GetMetadata()
	require.NoError(t, err, "Failed to load metadata from templates.yaml")

	// 代码中定义的 stage 名称
	codeStages := map[string]bool{
		stageRecon:       true,
		stageBruteforce:  true,
		stagePermutation: true,
		stageResolve:     true,
	}

	// 从 metadata 中提取 stage IDs
	metadataStages := make(map[string]bool)
	for _, stage := range metadata.Stages {
		metadataStages[stage.ID] = true
	}

	// 检查代码中的每个 stage 是否在 metadata 中存在
	for stageName := range codeStages {
		assert.True(t, metadataStages[stageName],
			"Stage '%s' is defined in code but missing in templates.yaml metadata", stageName)
	}

	// 检查 metadata 中的每个 stage 是否在代码中存在
	for stageName := range metadataStages {
		assert.True(t, codeStages[stageName],
			"Stage '%s' is defined in templates.yaml but missing in code constants", stageName)
	}
}

// TestToolsMatchMetadata 验证代码中的 tool 常量和 templates.yaml 中的定义一致
func TestToolsMatchMetadata(t *testing.T) {
	// 从 templates.yaml 加载所有工具模板
	templates, err := loader.Load()
	require.NoError(t, err, "Failed to load templates from templates.yaml")

	// 代码中定义的 tool 名称
	codeTools := map[string]bool{
		toolSubfinder:                   true,
		toolSublist3r:                   true,
		toolAssetfinder:                 true,
		toolSubdomainBruteforce:         true,
		toolSubdomainPermutationResolve: true,
		toolSubdomainResolve:            true,
	}

	// 从 templates 中提取 tool 名称
	metadataTools := make(map[string]bool)
	for toolName := range templates {
		metadataTools[toolName] = true
	}

	// 检查代码中的每个 tool 是否在 templates 中存在
	for toolName := range codeTools {
		assert.True(t, metadataTools[toolName],
			"Tool '%s' is defined in code but missing in templates.yaml", toolName)
	}

	// 注意：不检查反向（templates 中有但代码中没有），因为可能有未使用的工具定义
}

// TestReconToolsMatchMetadata 验证 reconTools 列表和实际定义的工具一致
func TestReconToolsMatchMetadata(t *testing.T) {
	// reconTools 应该包含所有侦察工具
	expectedReconTools := []string{
		toolSubfinder,
		toolSublist3r,
		toolAssetfinder,
	}

	assert.ElementsMatch(t, expectedReconTools, reconTools,
		"reconTools list should match the expected reconnaissance tools")
}

// TestStageToolMapping 验证每个 stage 都有对应的工具
func TestStageToolMapping(t *testing.T) {
	// 加载所有工具模板
	templates, err := loader.Load()
	require.NoError(t, err, "Failed to load templates")

	// 统计每个 stage 的工具数量
	stageTools := make(map[string][]string)
	for toolName, tmpl := range templates {
		stageTools[tmpl.Metadata.Stage] = append(stageTools[tmpl.Metadata.Stage], toolName)
	}

	// 验证每个 stage 至少有一个工具
	metadata, err := loader.GetMetadata()
	require.NoError(t, err, "Failed to load metadata")

	for _, stage := range metadata.Stages {
		tools := stageTools[stage.ID]
		assert.NotEmpty(t, tools,
			"Stage '%s' should have at least one tool defined", stage.ID)

		t.Logf("Stage '%s' has %d tool(s): %v", stage.ID, len(tools), tools)
	}
}

// TestGeneratedConstantsNotEmpty 验证生成的常量不为空
func TestGeneratedConstantsNotEmpty(t *testing.T) {
	// 验证 stage 常量
	assert.NotEmpty(t, stageRecon, "stageRecon should not be empty")
	assert.NotEmpty(t, stageBruteforce, "stageBruteforce should not be empty")
	assert.NotEmpty(t, stagePermutation, "stagePermutation should not be empty")
	assert.NotEmpty(t, stageResolve, "stageResolve should not be empty")

	// 验证 tool 常量
	assert.NotEmpty(t, toolSubfinder, "toolSubfinder should not be empty")
	assert.NotEmpty(t, toolSublist3r, "toolSublist3r should not be empty")
	assert.NotEmpty(t, toolAssetfinder, "toolAssetfinder should not be empty")
	assert.NotEmpty(t, toolSubdomainBruteforce, "toolSubdomainBruteforce should not be empty")
	assert.NotEmpty(t, toolSubdomainPermutationResolve, "toolSubdomainPermutationResolve should not be empty")
	assert.NotEmpty(t, toolSubdomainResolve, "toolSubdomainResolve should not be empty")
}

// TestTemplateYAMLStructure 验证 templates.yaml 的基本结构
func TestTemplateYAMLStructure(t *testing.T) {
	metadata, err := loader.GetMetadata()
	require.NoError(t, err, "Failed to load metadata")

	// 验证 workflow metadata
	assert.NotEmpty(t, metadata.Name, "Workflow name should not be empty")
	assert.NotEmpty(t, metadata.DisplayName, "Workflow display name should not be empty")
	assert.NotEmpty(t, metadata.Version, "Workflow version should not be empty")
	assert.NotEmpty(t, metadata.Stages, "Workflow should have at least one stage")
}
