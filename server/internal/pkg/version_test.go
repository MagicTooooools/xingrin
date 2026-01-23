package pkg

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadVersion(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "VERSION")
	if err := os.WriteFile(path, []byte("v1.2.3\n"), 0o600); err != nil {
		t.Fatalf("write version failed: %v", err)
	}

	if got := ReadVersion(path); got != "v1.2.3" {
		t.Fatalf("expected v1.2.3, got %q", got)
	}
}

func TestReadVersionUnknown(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "VERSION")
	if err := os.WriteFile(path, []byte("\n"), 0o600); err != nil {
		t.Fatalf("write version failed: %v", err)
	}
	if got := ReadVersion(path); got != "unknown" {
		t.Fatalf("expected unknown for empty file, got %q", got)
	}

	if got := ReadVersion(filepath.Join(dir, "missing")); got != "unknown" {
		t.Fatalf("expected unknown for missing file, got %q", got)
	}
}
