package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadVersionEnv(t *testing.T) {
	t.Setenv("AGENT_VERSION", "v1.2.3")

	got := ReadVersion()
	if got != "v1.2.3" {
		t.Fatalf("expected version from env, got %q", got)
	}
}

func TestReadVersionFromFile(t *testing.T) {
	t.Setenv("AGENT_VERSION", "")
	dir := t.TempDir()
	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd failed: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(oldWD)
	})
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}

	versionPath := filepath.Join(dir, "VERSION")
	if err := os.WriteFile(versionPath, []byte("v2.0.0\n"), 0o600); err != nil {
		t.Fatalf("write version failed: %v", err)
	}

	got := ReadVersion()
	if got != "v2.0.0" {
		t.Fatalf("expected version from file, got %q", got)
	}
}

func TestReadVersionUnknown(t *testing.T) {
	t.Setenv("AGENT_VERSION", "")
	dir := t.TempDir()
	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd failed: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(oldWD)
	})
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}

	got := ReadVersion()
	if got != "unknown" {
		t.Fatalf("expected unknown, got %q", got)
	}
}
