package config

import (
	"os"
	"testing"
)

func TestLoadConfigFromEnvAndFlags(t *testing.T) {
	t.Setenv("SERVER_URL", "https://example.com")
	t.Setenv("API_KEY", "abc12345")
	t.Setenv("MAX_TASKS", "5")
	t.Setenv("CPU_THRESHOLD", "80")
	t.Setenv("MEM_THRESHOLD", "81")
	t.Setenv("DISK_THRESHOLD", "82")

	cfg, err := Load([]string{})
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}
	if cfg.ServerURL != "https://example.com" {
		t.Fatalf("expected server url from env")
	}
	if cfg.MaxTasks != 5 {
		t.Fatalf("expected max tasks from env")
	}

	args := []string{
		"--server-url=https://override.example.com",
		"--api-key=deadbeef",
		"--max-tasks=9",
		"--cpu-threshold=70",
		"--mem-threshold=71",
		"--disk-threshold=72",
	}
	cfg, err = Load(args)
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}
	if cfg.ServerURL != "https://override.example.com" {
		t.Fatalf("expected server url from args")
	}
	if cfg.APIKey != "deadbeef" {
		t.Fatalf("expected api key from args")
	}
	if cfg.MaxTasks != 9 {
		t.Fatalf("expected max tasks from args")
	}
	if cfg.CPUThreshold != 70 || cfg.MemThreshold != 71 || cfg.DiskThreshold != 72 {
		t.Fatalf("expected thresholds from args")
	}
}

func TestLoadConfigMissingRequired(t *testing.T) {
	t.Setenv("SERVER_URL", "")
	t.Setenv("API_KEY", "")

	_, err := Load([]string{})
	if err == nil {
		t.Fatalf("expected error when required values missing")
	}
}

func TestReadVersionPrefersEnv(t *testing.T) {
	t.Setenv("AGENT_VERSION", "v9.9.9")
	if got := ReadVersion(); got != "v9.9.9" {
		t.Fatalf("expected version from env, got %s", got)
	}
}

func TestReadVersionFileFallback(t *testing.T) {
	tmp := t.TempDir()
	path := tmp + "/VERSION"
	if err := os.WriteFile(path, []byte("v1.2.3\n"), 0644); err != nil {
		t.Fatalf("write version file: %v", err)
	}
	if got := readVersionFile(path); got != "v1.2.3" {
		t.Fatalf("expected version from file, got %s", got)
	}
}
