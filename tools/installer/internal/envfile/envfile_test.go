package envfile

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRenderContainsPublicURL(t *testing.T) {
	content := Render(Data{
		ImageTag:    "dev",
		JWTSecret:   "jwt",
		WorkerToken: "worker",
		DBHost:      "postgres",
		DBPassword:  "postgres",
		RedisHost:   "redis",
		DBUser:      "postgres",
		DBName:      "lunafox",
		DBPort:      "5432",
		RedisPort:   "6379",
		Go111Module: "on",
		GoProxy:     "https://proxy.golang.org,direct",
		PublicURL:   "https://example.com:8083",
	})

	if !strings.Contains(content, "PUBLIC_URL=https://example.com:8083") {
		t.Fatalf("PUBLIC_URL not found in env content: %s", content)
	}
}

func TestReadWorkerToken(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte("WORKER_TOKEN=abc123\n"), 0o644); err != nil {
		t.Fatalf("write env: %v", err)
	}

	token, err := ReadWorkerToken(path)
	if err != nil {
		t.Fatalf("read worker token: %v", err)
	}
	if token != "abc123" {
		t.Fatalf("unexpected token: %s", token)
	}
}
