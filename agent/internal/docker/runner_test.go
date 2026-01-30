package docker

import (
	"testing"

	"github.com/yyhuni/lunafox/agent/internal/domain"
)

func TestResolveWorkerImage(t *testing.T) {
	if _, err := resolveWorkerImage(""); err == nil {
		t.Fatalf("expected error for empty version")
	}
	if got, err := resolveWorkerImage("v1.2.3"); err != nil || got != workerImagePrefix+"v1.2.3" {
		t.Fatalf("expected version image, got %s, err: %v", got, err)
	}
}

func TestBuildWorkerEnv(t *testing.T) {
	spec := &domain.Task{
		ScanID:       1,
		TargetID:     2,
		TargetName:   "example.com",
		TargetType:   "domain",
		WorkflowName: "subdomain_discovery",
		WorkspaceDir: "/opt/lunafox/results",
		Config:       "config-yaml",
	}

	env := buildWorkerEnv(spec, "https://server", "token")
	expected := []string{
		"SERVER_URL=https://server",
		"SERVER_TOKEN=token",
		"SCAN_ID=1",
		"TARGET_ID=2",
		"TARGET_NAME=example.com",
		"TARGET_TYPE=domain",
		"WORKFLOW_NAME=subdomain_discovery",
		"WORKSPACE_DIR=/opt/lunafox/results",
		"CONFIG=config-yaml",
	}

	if len(env) != len(expected) {
		t.Fatalf("expected %d env entries, got %d", len(expected), len(env))
	}
	for i, item := range expected {
		if env[i] != item {
			t.Fatalf("expected env[%d]=%s got %s", i, item, env[i])
		}
	}
}
