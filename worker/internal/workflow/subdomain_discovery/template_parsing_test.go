package subdomain_discovery

import (
	"fmt"
	"regexp"
	"reflect"
	"strings"
	"testing"
	"text/template"
	"text/template/parse"

	"github.com/yyhuni/lunafox/worker/internal/activity"
	"github.com/yyhuni/lunafox/worker/internal/pkg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func assertTemplateUsesVars(t *testing.T, s string, vars ...string) {
	t.Helper()
	for _, v := range vars {
		re := regexp.MustCompile(fmt.Sprintf(`(?s){{[^}]*\.%s[^}]*}}`, regexp.QuoteMeta(v)))
		assert.Regexp(t, re, s, "template should reference .%s", v)
	}
}

func dummyValueForType(t *testing.T, typ string) any {
	t.Helper()
	switch typ {
	case "string":
		return "value"
	case "integer":
		return 1
	case "boolean":
		return true
	default:
		t.Fatalf("unsupported config_schema.type %q", typ)
		return nil
	}
}

func templateVars(t *testing.T, s string) map[string]struct{} {
	t.Helper()

	funcMap := template.FuncMap{
		"quote": func(s string) string { return fmt.Sprintf("%q", s) },
		"lower": func(s string) string { return s },
		"upper": func(s string) string { return s },
		"join":  func(sep string, elems []string) string { return "" },
	}

	tmpl, err := template.New("t").Funcs(funcMap).Parse(s)
	require.NoError(t, err)

	vars := make(map[string]struct{})
	var walk func(n parse.Node)
	walk = func(n parse.Node) {
		if n == nil {
			return
		}
		// parse.Node is an interface; it can hold a typed-nil pointer.
		v := reflect.ValueOf(n)
		if v.Kind() == reflect.Ptr && v.IsNil() {
			return
		}
		switch nn := n.(type) {
		case *parse.ListNode:
			for _, child := range nn.Nodes {
				walk(child)
			}
		case *parse.ActionNode:
			walk(nn.Pipe)
		case *parse.PipeNode:
			for _, cmd := range nn.Cmds {
				walk(cmd)
			}
		case *parse.CommandNode:
			for _, arg := range nn.Args {
				walk(arg)
			}
		case *parse.FieldNode:
			if len(nn.Ident) > 0 {
				vars[nn.Ident[0]] = struct{}{}
			}
		case *parse.IfNode:
			walk(nn.Pipe)
			walk(nn.List)
			walk(nn.ElseList)
		case *parse.RangeNode:
			walk(nn.Pipe)
			walk(nn.List)
			walk(nn.ElseList)
		case *parse.WithNode:
			walk(nn.Pipe)
			walk(nn.List)
			walk(nn.ElseList)
		case *parse.TemplateNode:
			walk(nn.Pipe)
		case *parse.ChainNode:
			walk(nn.Node)
			if len(nn.Field) > 0 {
				vars[nn.Field[0]] = struct{}{}
			}
		default:
			// Ignore TextNode, IdentifierNode, VariableNode, etc.
		}
	}

	walk(tmpl.Tree.Root)
	return vars
}

func buildConfigRequired(t *testing.T, tmpl activity.CommandTemplate) map[string]any {
	t.Helper()

	cfg := make(map[string]any)
	for _, p := range append(tmpl.RuntimeParams, tmpl.CLIParams...) {
		if p.ConfigSchema.Required {
			cfg[p.SemanticID] = dummyValueForType(t, p.ConfigSchema.Type)
		}
	}
	return cfg
}

func buildConfigAll(t *testing.T, tmpl activity.CommandTemplate) map[string]any {
	t.Helper()

	cfg := make(map[string]any)
	for _, p := range append(tmpl.RuntimeParams, tmpl.CLIParams...) {
		cfg[p.SemanticID] = dummyValueForType(t, p.ConfigSchema.Type)
	}
	return cfg
}

func buildAutoParamsForTemplate(t *testing.T, tmpl activity.CommandTemplate, cfg map[string]any) map[string]any {
	t.Helper()

	// Determine which vars are provided via config (param.Var) and which vars belong to template params.
	paramVars := make(map[string]struct{})
	varsFromConfig := make(map[string]struct{})
	for _, p := range append(tmpl.RuntimeParams, tmpl.CLIParams...) {
		paramVars[p.Var] = struct{}{}
		if _, ok := cfg[p.SemanticID]; ok {
			varsFromConfig[p.Var] = struct{}{}
		}
	}

	used := templateVars(t, tmpl.BaseCommand)
	// CLI arg templates are appended only when the corresponding semantic_id exists in cfg.
	for _, p := range tmpl.CLIParams {
		if p.Arg == "" {
			continue
		}
		if _, ok := cfg[p.SemanticID]; !ok {
			continue
		}
		for v := range templateVars(t, p.Arg) {
			used[v] = struct{}{}
		}
	}

	params := make(map[string]any)
	for v := range used {
		if _, ok := varsFromConfig[v]; ok {
			continue
		}
		// Avoid populating template param variables via params, otherwise optional CLI args could be
		// accidentally enabled without going through config type validation.
		if _, ok := paramVars[v]; ok {
			continue
		}
		params[v] = fmt.Sprintf("%s-value", strings.ToLower(v))
	}

	return params
}

// TestLoadAllTemplates 验证所有工具模板都能成功加载
func TestLoadAllTemplates(t *testing.T) {
	require.NoError(t, pkg.InitLogger("error"))
	defer pkg.Sync()

	templates, err := loader.Load()
	require.NoError(t, err, "Failed to load templates from templates.yaml")
	require.NotEmpty(t, templates, "Templates should not be empty")

	// 验证所有预期的工具都存在
	expectedTools := []string{
		toolSubfinder,
		toolSublist3r,
		toolAssetfinder,
		toolSubdomainBruteforce,
		toolSubdomainResolve,
		toolSubdomainPermutationResolve,
	}

	for _, toolName := range expectedTools {
		tmpl, exists := templates[toolName]
		assert.True(t, exists, "Tool %s should exist in templates", toolName)
		assert.NotEmpty(t, tmpl.BaseCommand, "Tool %s should have a base command", toolName)
		assert.NotEmpty(t, tmpl.Metadata.Stage, "Tool %s should have a stage", toolName)
		assert.NotEmpty(t, tmpl.Metadata.DisplayName, "Tool %s should have a display name", toolName)
	}
}

// TestGetTemplate 验证 getTemplate 函数能正确获取单个模板
func TestGetTemplate(t *testing.T) {
	require.NoError(t, pkg.InitLogger("error"))
	defer pkg.Sync()
	tests := []struct {
		name     string
		toolName string
		wantErr  bool
	}{
		{
			name:     "Get subfinder template",
			toolName: toolSubfinder,
			wantErr:  false,
		},
		{
			name:     "Get non-existent tool",
			toolName: "non-existent-tool",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpl, err := getTemplate(tt.toolName)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.NotEmpty(t, tmpl.BaseCommand)
			}
		})
	}
}

// TestSubfinderTemplateStructure 验证 subfinder 模板的结构
func TestSubfinderTemplateStructure(t *testing.T) {
	require.NoError(t, pkg.InitLogger("error"))
	defer pkg.Sync()
	tmpl, err := getTemplate(toolSubfinder)
	require.NoError(t, err)

	// 验证 metadata
	assert.Equal(t, "Subfinder", tmpl.Metadata.DisplayName)
	assert.Equal(t, stageRecon, tmpl.Metadata.Stage)
	assert.NotEmpty(t, tmpl.Metadata.Description)

	// 验证 base command 包含必要的占位符（允许被 quote 等函数包裹）
	assertTemplateUsesVars(t, tmpl.BaseCommand, "Domain", "OutputFile")

	// 验证 runtime_params
	assert.NotEmpty(t, tmpl.RuntimeParams, "Subfinder should have runtime params")
	hasTimeout := false
	for _, param := range tmpl.RuntimeParams {
		if param.SemanticID == "timeout-runtime" {
			hasTimeout = true
			assert.Equal(t, "Timeout", param.Var)
			assert.Equal(t, "integer", param.ConfigSchema.Type)
			assert.True(t, param.ConfigSchema.Required)
		}
	}
	assert.True(t, hasTimeout, "Subfinder should have timeout-runtime parameter")

	// 验证 cli_params
	assert.NotEmpty(t, tmpl.CLIParams, "Subfinder should have CLI params")
	hasThreads := false
	for _, param := range tmpl.CLIParams {
		if param.SemanticID == "threads-cli" {
			hasThreads = true
			assert.Equal(t, "Threads", param.Var)
			assert.Contains(t, param.Arg, "{{.Threads}}")
		}
	}
	assert.True(t, hasThreads, "Subfinder should have threads-cli parameter")
}

// TestSubdomainBruteforceTemplateStructure 验证 subdomain-bruteforce 模板的结构
func TestSubdomainBruteforceTemplateStructure(t *testing.T) {
	require.NoError(t, pkg.InitLogger("error"))
	defer pkg.Sync()
	tmpl, err := getTemplate(toolSubdomainBruteforce)
	require.NoError(t, err)

	// 验证 metadata
	assert.Equal(t, "Subdomain Bruteforce", tmpl.Metadata.DisplayName)
	assert.Equal(t, stageBruteforce, tmpl.Metadata.Stage)

	// 验证 internal_params
	assert.NotEmpty(t, tmpl.InternalParams, "Subdomain bruteforce should have internal params")
	assert.Contains(t, tmpl.InternalParams, "subdomain-wordlist-base-path-runtime")
	assert.Contains(t, tmpl.InternalParams, "resolvers-path-cli")

	// 验证 base command（允许被 quote 等函数包裹）
	assert.Contains(t, tmpl.BaseCommand, "puredns bruteforce")
	assertTemplateUsesVars(t, tmpl.BaseCommand, "Wordlist", "Domain", "Resolvers", "OutputFile")

	// 验证特定的 CLI 参数
	paramNames := []string{"threads-cli", "rate-limit-cli", "wildcard-tests-cli", "wildcard-batch-cli"}
	for _, paramName := range paramNames {
		found := false
		for _, param := range tmpl.CLIParams {
			if param.SemanticID == paramName {
				found = true
				break
			}
		}
		assert.True(t, found, "Should have parameter: %s", paramName)
	}
}

// TestSubdomainPermutationResolveTemplateStructure 验证 subdomain-permutation-resolve 模板的结构
func TestSubdomainPermutationResolveTemplateStructure(t *testing.T) {
	require.NoError(t, pkg.InitLogger("error"))
	defer pkg.Sync()
	tmpl, err := getTemplate(toolSubdomainPermutationResolve)
	require.NoError(t, err)

	// 验证 metadata
	assert.Equal(t, "Subdomain Permutation Resolve", tmpl.Metadata.DisplayName)
	assert.Equal(t, stagePermutation, tmpl.Metadata.Stage)

	// 验证 base command 包含管道操作（允许被 quote 等函数包裹）
	assert.Contains(t, tmpl.BaseCommand, "cat")
	assert.Contains(t, tmpl.BaseCommand, "dnsgen")
	assert.Contains(t, tmpl.BaseCommand, "puredns resolve")
	assertTemplateUsesVars(t, tmpl.BaseCommand, "InputFile", "Resolvers", "OutputFile")

	// 验证特定的 runtime 参数
	wildcardParams := []string{
		"wildcard-sample-timeout-runtime",
		"wildcard-sample-multiplier-runtime",
		"wildcard-expansion-threshold-runtime",
	}
	for _, paramName := range wildcardParams {
		found := false
		for _, param := range tmpl.RuntimeParams {
			if param.SemanticID == paramName {
				found = true
				assert.Equal(t, "integer", param.ConfigSchema.Type)
				break
			}
		}
		assert.True(t, found, "Should have parameter: %s", paramName)
	}
}

// TestCommandGeneration_Subfinder 测试 subfinder 命令生成
func TestCommandGeneration_Subfinder(t *testing.T) {
	require.NoError(t, pkg.InitLogger("error"))
	defer pkg.Sync()
	tmpl, err := getTemplate(toolSubfinder)
	require.NoError(t, err)

	builder := activity.NewCommandBuilder()

	// 测试用例：基本参数
	params := map[string]any{
		"Domain":         "example.com",
		"OutputFile":     "/tmp/output.txt",
		"ProviderConfig": "", // optional in template but must exist when using missingkey=error
	}
	config := map[string]any{
		"timeout-runtime": 3600,
		"threads-cli":     10,
	}

	cmd, err := builder.Build(tmpl, params, config)
	require.NoError(t, err)
	assert.NotEmpty(t, cmd)

	// 验证命令包含必要的部分
	assert.Contains(t, cmd, "subfinder")
	assert.Contains(t, cmd, "-d example.com")
	assert.Contains(t, cmd, "-all")
	assert.Contains(t, cmd, "-o \"/tmp/output.txt\"")
	assert.Contains(t, cmd, "-v")
	assert.Contains(t, cmd, "-t 10")

	t.Logf("Generated command: %s", cmd)
}

// TestCommandGeneration_SubfinderWithProviderConfig 测试带 provider config 的 subfinder 命令
func TestCommandGeneration_SubfinderWithProviderConfig(t *testing.T) {
	require.NoError(t, pkg.InitLogger("error"))
	defer pkg.Sync()
	tmpl, err := getTemplate(toolSubfinder)
	require.NoError(t, err)

	builder := activity.NewCommandBuilder()

	params := map[string]any{
		"Domain":         "example.com",
		"OutputFile":     "/tmp/output.txt",
		"ProviderConfig": "/etc/subfinder/config.yaml",
	}
	config := map[string]any{
		"timeout-runtime": 3600,
		"threads-cli":     10,
	}

	cmd, err := builder.Build(tmpl, params, config)
	require.NoError(t, err)

	// 验证包含 provider config
	assert.Contains(t, cmd, "-pc \"/etc/subfinder/config.yaml\"")

	t.Logf("Generated command: %s", cmd)
}

// TestCommandGeneration_Sublist3r 测试 sublist3r 命令生成
func TestCommandGeneration_Sublist3r(t *testing.T) {
	require.NoError(t, pkg.InitLogger("error"))
	defer pkg.Sync()
	tmpl, err := getTemplate(toolSublist3r)
	require.NoError(t, err)

	builder := activity.NewCommandBuilder()

	params := map[string]any{
		"Domain":     "example.com",
		"OutputFile": "/tmp/output.txt",
	}
	config := map[string]any{
		"timeout-runtime": 3600,
		"threads-cli":     10,
	}

	cmd, err := builder.Build(tmpl, params, config)
	require.NoError(t, err)

	assert.Contains(t, cmd, "python3")
	assert.Contains(t, cmd, "/opt/lunafox-tools/share/Sublist3r/sublist3r.py")
	assert.Contains(t, cmd, "-d example.com")
	assert.Contains(t, cmd, "-o \"/tmp/output.txt\"")
	assert.Contains(t, cmd, "-t 10")

	t.Logf("Generated command: %s", cmd)
}

// TestCommandGeneration_Assetfinder 测试 assetfinder 命令生成
func TestCommandGeneration_Assetfinder(t *testing.T) {
	require.NoError(t, pkg.InitLogger("error"))
	defer pkg.Sync()
	tmpl, err := getTemplate(toolAssetfinder)
	require.NoError(t, err)

	builder := activity.NewCommandBuilder()

	params := map[string]any{
		"Domain":     "example.com",
		"OutputFile": "/tmp/output.txt",
	}
	config := map[string]any{
		"timeout-runtime": 3600,
	}

	cmd, err := builder.Build(tmpl, params, config)
	require.NoError(t, err)

	assert.Contains(t, cmd, "assetfinder")
	assert.Contains(t, cmd, "--subs-only")
	assert.Contains(t, cmd, "example.com")
	assert.Contains(t, cmd, "> \"/tmp/output.txt\"")

	t.Logf("Generated command: %s", cmd)
}

// TestCommandGeneration_SubdomainBruteforce 测试 subdomain-bruteforce 命令生成
func TestCommandGeneration_SubdomainBruteforce(t *testing.T) {
	require.NoError(t, pkg.InitLogger("error"))
	defer pkg.Sync()
	tmpl, err := getTemplate(toolSubdomainBruteforce)
	require.NoError(t, err)

	builder := activity.NewCommandBuilder()

	params := map[string]any{
		"Domain":     "example.com",
		"OutputFile": "/tmp/output.txt",
		"Wordlist":   "/opt/lunafox/wordlists/subdomains-top1million-110000.txt",
		"Resolvers":  "/opt/lunafox/wordlists/resolvers.txt",
	}
	config := map[string]any{
		"timeout-runtime":                   3600,
		"subdomain-wordlist-name-runtime":   "subdomains-top1million-110000.txt",
		"threads-cli":                       100,
		"rate-limit-cli":                    150,
		"wildcard-tests-cli":                50,
		"wildcard-batch-cli":                1000000,
	}

	cmd, err := builder.Build(tmpl, params, config)
	require.NoError(t, err)

	assert.Contains(t, cmd, "puredns bruteforce")
	assert.Contains(t, cmd, "\"/opt/lunafox/wordlists/subdomains-top1million-110000.txt\"")
	assert.Contains(t, cmd, "example.com")
	assert.Contains(t, cmd, "-r \"/opt/lunafox/wordlists/resolvers.txt\"")
	assert.Contains(t, cmd, "--write \"/tmp/output.txt\"")
	assert.Contains(t, cmd, "-t 100")
	assert.Contains(t, cmd, "--rate-limit 150")
	assert.Contains(t, cmd, "--wildcard-tests 50")
	assert.Contains(t, cmd, "--wildcard-batch 1000000")

	t.Logf("Generated command: %s", cmd)
}

// TestCommandGeneration_SubdomainResolve 测试 subdomain-resolve 命令生成
func TestCommandGeneration_SubdomainResolve(t *testing.T) {
	require.NoError(t, pkg.InitLogger("error"))
	defer pkg.Sync()
	tmpl, err := getTemplate(toolSubdomainResolve)
	require.NoError(t, err)

	builder := activity.NewCommandBuilder()

	params := map[string]any{
		"InputFile":  "/tmp/input.txt",
		"OutputFile": "/tmp/output.txt",
		"Resolvers":  "/opt/lunafox/wordlists/resolvers.txt",
	}
	config := map[string]any{
		"timeout-runtime": 3600,
		"threads-cli":     100,
		"rate-limit-cli":  150,
	}

	cmd, err := builder.Build(tmpl, params, config)
	require.NoError(t, err)

	assert.Contains(t, cmd, "puredns resolve")
	assert.Contains(t, cmd, "\"/tmp/input.txt\"")
	assert.Contains(t, cmd, "-r \"/opt/lunafox/wordlists/resolvers.txt\"")
	assert.Contains(t, cmd, "--write \"/tmp/output.txt\"")
	assert.Contains(t, cmd, "-t 100")
	assert.Contains(t, cmd, "--rate-limit 150")

	t.Logf("Generated command: %s", cmd)
}

// TestCommandGeneration_SubdomainPermutationResolve 测试 subdomain-permutation-resolve 命令生成
func TestCommandGeneration_SubdomainPermutationResolve(t *testing.T) {
	require.NoError(t, pkg.InitLogger("error"))
	defer pkg.Sync()
	tmpl, err := getTemplate(toolSubdomainPermutationResolve)
	require.NoError(t, err)

	builder := activity.NewCommandBuilder()

	params := map[string]any{
		"InputFile":  "/tmp/input.txt",
		"OutputFile": "/tmp/output.txt",
		"Resolvers":  "/opt/lunafox/wordlists/resolvers.txt",
	}
	config := map[string]any{
		"timeout-runtime":                      3600,
		"wildcard-sample-timeout-runtime":      7200,
		"wildcard-sample-multiplier-runtime":   100,
		"wildcard-expansion-threshold-runtime": 50,
		"threads-cli":                          100,
		"rate-limit-cli":                       150,
		"wildcard-tests-cli":                   50,
		"wildcard-batch-cli":                   1000000,
	}

	cmd, err := builder.Build(tmpl, params, config)
	require.NoError(t, err)

	// 验证管道命令结构
	assert.Contains(t, cmd, "cat \"/tmp/input.txt\"")
	assert.Contains(t, cmd, "dnsgen -")
	assert.Contains(t, cmd, "puredns resolve")
	assert.Contains(t, cmd, "-r \"/opt/lunafox/wordlists/resolvers.txt\"")
	assert.Contains(t, cmd, "--write \"/tmp/output.txt\"")
	assert.Contains(t, cmd, "-t 100")
	assert.Contains(t, cmd, "--rate-limit 150")
	assert.Contains(t, cmd, "--wildcard-tests 50")
	assert.Contains(t, cmd, "--wildcard-batch 1000000")

	t.Logf("Generated command: %s", cmd)
}

// TestParameterValidation_MissingRequiredParam 测试缺少必填参数的情况
func TestParameterValidation_MissingRequiredParam(t *testing.T) {
	require.NoError(t, pkg.InitLogger("error"))
	defer pkg.Sync()
	tmpl, err := getTemplate(toolSubfinder)
	require.NoError(t, err)

	builder := activity.NewCommandBuilder()

	params := map[string]any{
		"Domain":     "example.com",
		"OutputFile": "/tmp/output.txt",
	}
	// 缺少必填的 timeout-runtime 参数
	config := map[string]any{
		"threads-cli": 10,
	}

	_, err = builder.Build(tmpl, params, config)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "required parameter")
	assert.Contains(t, err.Error(), "timeout-runtime")
}

// TestParameterValidation_InvalidType 测试参数类型错误的情况
func TestParameterValidation_InvalidType(t *testing.T) {
	require.NoError(t, pkg.InitLogger("error"))
	defer pkg.Sync()
	tmpl, err := getTemplate(toolSubfinder)
	require.NoError(t, err)

	builder := activity.NewCommandBuilder()

	params := map[string]any{
		"Domain":     "example.com",
		"OutputFile": "/tmp/output.txt",
	}
	// threads-cli 应该是 integer，但提供了 string
	config := map[string]any{
		"timeout-runtime": 3600,
		"threads-cli":     "invalid",
	}

	_, err = builder.Build(tmpl, params, config)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "expected integer")
}

// TestParameterValidation_OptionalParamOmitted 测试可选参数缺失的情况（应该成功）
func TestParameterValidation_OptionalParamOmitted(t *testing.T) {
	require.NoError(t, pkg.InitLogger("error"))
	defer pkg.Sync()
	tmpl, err := getTemplate(toolAssetfinder)
	require.NoError(t, err)

	builder := activity.NewCommandBuilder()

	params := map[string]any{
		"Domain":     "example.com",
		"OutputFile": "/tmp/output.txt",
	}
	// assetfinder 只有 timeout-runtime 是必填的，没有其他必填的 CLI 参数
	config := map[string]any{
		"timeout-runtime": 3600,
	}

	cmd, err := builder.Build(tmpl, params, config)
	require.NoError(t, err)
	assert.NotEmpty(t, cmd)
}

// TestInternalParams_SubdomainBruteforce 测试 internal_params 的正确性
func TestInternalParams_SubdomainBruteforce(t *testing.T) {
	require.NoError(t, pkg.InitLogger("error"))
	defer pkg.Sync()
	tmpl, err := getTemplate(toolSubdomainBruteforce)
	require.NoError(t, err)

	// 验证 internal_params 存在且值正确
	assert.NotEmpty(t, tmpl.InternalParams)
	assert.Equal(t, "/opt/lunafox/wordlists", tmpl.InternalParams["subdomain-wordlist-base-path-runtime"])
	assert.Equal(t, "/opt/lunafox/wordlists/resolvers.txt", tmpl.InternalParams["resolvers-path-cli"])
}

// TestInternalParams_SubdomainResolve 测试 subdomain-resolve 的 internal_params
func TestInternalParams_SubdomainResolve(t *testing.T) {
	require.NoError(t, pkg.InitLogger("error"))
	defer pkg.Sync()
	tmpl, err := getTemplate(toolSubdomainResolve)
	require.NoError(t, err)

	assert.NotEmpty(t, tmpl.InternalParams)
	assert.Equal(t, "/opt/lunafox/wordlists/resolvers.txt", tmpl.InternalParams["resolvers-path-cli"])
}

// TestAllToolsHaveRequiredMetadata 验证所有工具都有必需的 metadata 字段
func TestAllToolsHaveRequiredMetadata(t *testing.T) {
	require.NoError(t, pkg.InitLogger("error"))
	defer pkg.Sync()
	templates, err := loader.Load()
	require.NoError(t, err)

	for toolName, tmpl := range templates {
		t.Run(toolName, func(t *testing.T) {
			assert.NotEmpty(t, tmpl.Metadata.DisplayName, "Tool %s should have display_name", toolName)
			assert.NotEmpty(t, tmpl.Metadata.Description, "Tool %s should have description", toolName)
			assert.NotEmpty(t, tmpl.Metadata.Stage, "Tool %s should have stage", toolName)
			assert.NotEmpty(t, tmpl.BaseCommand, "Tool %s should have base_command", toolName)
		})
	}
}

// TestAllToolsHaveTimeoutParam 验证所有工具都有 timeout-runtime 参数
func TestAllToolsHaveTimeoutParam(t *testing.T) {
	require.NoError(t, pkg.InitLogger("error"))
	defer pkg.Sync()
	templates, err := loader.Load()
	require.NoError(t, err)

	for toolName, tmpl := range templates {
		t.Run(toolName, func(t *testing.T) {
			hasTimeout := false
			for _, param := range tmpl.RuntimeParams {
				if param.SemanticID == "timeout-runtime" {
					hasTimeout = true
					assert.Equal(t, "Timeout", param.Var, "Tool %s timeout var should be 'Timeout'", toolName)
					assert.Equal(t, "integer", param.ConfigSchema.Type, "Tool %s timeout should be integer", toolName)
					assert.True(t, param.ConfigSchema.Required, "Tool %s timeout should be required", toolName)
					break
				}
			}
			assert.True(t, hasTimeout, "Tool %s should have timeout-runtime parameter", toolName)
		})
	}
}

// TestReconToolsHaveThreadsParam 验证所有 recon 工具都有 threads-cli 参数
func TestReconToolsHaveThreadsParam(t *testing.T) {
	require.NoError(t, pkg.InitLogger("error"))
	defer pkg.Sync()
	reconToolNames := []string{toolSubfinder, toolSublist3r}

	for _, toolName := range reconToolNames {
		t.Run(toolName, func(t *testing.T) {
			tmpl, err := getTemplate(toolName)
			require.NoError(t, err)

			hasThreads := false
			for _, param := range tmpl.CLIParams {
				if param.SemanticID == "threads-cli" {
					hasThreads = true
					assert.Equal(t, "Threads", param.Var)
					assert.Contains(t, param.Arg, "{{.Threads}}")
					break
				}
			}
			assert.True(t, hasThreads, "Recon tool %s should have threads-cli parameter", toolName)
		})
	}
}

// TestCommandTemplateQuoting 测试命令模板中的引号处理
func TestCommandTemplateQuoting(t *testing.T) {
	require.NoError(t, pkg.InitLogger("error"))
	defer pkg.Sync()
	tests := []struct {
		name     string
		toolName string
		params   map[string]any
		config   map[string]any
		checks   []string
	}{
		{
			name:     "Subfinder with spaces in path",
			toolName: toolSubfinder,
			params: map[string]any{
				"Domain":         "example.com",
				"OutputFile":     "/tmp/path with spaces/output.txt",
				"ProviderConfig": "", // optional in template but must exist when using missingkey=error
			},
			config: map[string]any{
				"timeout-runtime": 3600,
				"threads-cli":     10,
			},
			checks: []string{
				`-o "/tmp/path with spaces/output.txt"`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpl, err := getTemplate(tt.toolName)
			require.NoError(t, err)

			builder := activity.NewCommandBuilder()
			cmd, err := builder.Build(tmpl, tt.params, tt.config)
			require.NoError(t, err)

			for _, check := range tt.checks {
				assert.Contains(t, cmd, check, "Command should properly quote paths with spaces")
			}

			t.Logf("Generated command: %s", cmd)
		})
	}
}

// TestCommandGeneration_NoExtraSpaces 验证生成的命令没有多余的空格
func TestCommandGeneration_NoExtraSpaces(t *testing.T) {
	require.NoError(t, pkg.InitLogger("error"))
	defer pkg.Sync()
	tmpl, err := getTemplate(toolSubfinder)
	require.NoError(t, err)

	builder := activity.NewCommandBuilder()

	params := map[string]any{
		"Domain":         "example.com",
		"OutputFile":     "/tmp/output.txt",
		"ProviderConfig": "", // optional in template but must exist when using missingkey=error
	}
	config := map[string]any{
		"timeout-runtime": 3600,
		"threads-cli":     10,
	}

	cmd, err := builder.Build(tmpl, params, config)
	require.NoError(t, err)

	// 验证没有连续的多个空格
	assert.NotContains(t, cmd, "  ", "Command should not contain double spaces")

	// 验证命令前后没有空格
	assert.Equal(t, strings.TrimSpace(cmd), cmd, "Command should be trimmed")
}

// TestStageAssignment 验证每个工具都分配到了正确的 stage
func TestStageAssignment(t *testing.T) {
	require.NoError(t, pkg.InitLogger("error"))
	defer pkg.Sync()
	expectedStages := map[string]string{
		toolSubfinder:                   stageRecon,
		toolSublist3r:                   stageRecon,
		toolAssetfinder:                 stageRecon,
		toolSubdomainBruteforce:         stageBruteforce,
		toolSubdomainResolve:            stageResolve,
		toolSubdomainPermutationResolve: stagePermutation,
	}

	for toolName, expectedStage := range expectedStages {
		t.Run(toolName, func(t *testing.T) {
			tmpl, err := getTemplate(toolName)
			require.NoError(t, err)
			assert.Equal(t, expectedStage, tmpl.Metadata.Stage,
				"Tool %s should be in stage %s", toolName, expectedStage)
		})
	}
}

func TestAllTemplates_Build_MinimalRequiredConfig(t *testing.T) {
	require.NoError(t, pkg.InitLogger("error"))
	defer pkg.Sync()

	templates, err := loader.Load()
	require.NoError(t, err)
	require.NotEmpty(t, templates)

	for toolName, tmpl := range templates {
		t.Run(toolName, func(t *testing.T) {
			cfg := buildConfigRequired(t, tmpl)
			params := buildAutoParamsForTemplate(t, tmpl, cfg)

			builder := activity.NewCommandBuilder()
			cmd, err := builder.Build(tmpl, params, cfg)
			require.NoError(t, err)
			assert.NotEmpty(t, cmd)
			assert.NotContains(t, cmd, "{{")
			assert.NotContains(t, cmd, "}}")
		})
	}
}

func TestAllTemplates_Build_WithAllParams(t *testing.T) {
	require.NoError(t, pkg.InitLogger("error"))
	defer pkg.Sync()

	templates, err := loader.Load()
	require.NoError(t, err)
	require.NotEmpty(t, templates)

	for toolName, tmpl := range templates {
		t.Run(toolName, func(t *testing.T) {
			cfg := buildConfigAll(t, tmpl)
			params := buildAutoParamsForTemplate(t, tmpl, cfg)

			builder := activity.NewCommandBuilder()
			cmd, err := builder.Build(tmpl, params, cfg)
			require.NoError(t, err)
			assert.NotEmpty(t, cmd)
			assert.NotContains(t, cmd, "{{")
			assert.NotContains(t, cmd, "}}")
		})
	}
}

func TestAllTemplates_BaseCommand_DoesNotReferenceOptionalConfigVars(t *testing.T) {
	require.NoError(t, pkg.InitLogger("error"))
	defer pkg.Sync()

	templates, err := loader.Load()
	require.NoError(t, err)
	require.NotEmpty(t, templates)

	for toolName, tmpl := range templates {
		t.Run(toolName, func(t *testing.T) {
			used := templateVars(t, tmpl.BaseCommand)
			for _, p := range append(tmpl.RuntimeParams, tmpl.CLIParams...) {
				if p.ConfigSchema.Required {
					continue
				}
				if _, ok := used[p.Var]; ok {
					t.Fatalf(
						"base_command references optional config param var %q (semantic_id=%q). "+
							"Optional params should not be referenced in base_command when using missingkey=error; "+
							"move it to cli_params.arg (gated by config), make it required, or ensure it is always provided via params.",
						p.Var,
						p.SemanticID,
					)
				}
			}
		})
	}
}
