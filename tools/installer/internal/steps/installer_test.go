package steps

import (
	"crypto/tls"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yyhuni/lunafox/tools/installer/internal/cli"
	"github.com/yyhuni/lunafox/tools/installer/internal/ui"
)

func TestParseBuildxDriver(t *testing.T) {
	driver := parseBuildxDriver("Name: default\nDriver: docker-container\n")
	if driver != "docker-container" {
		t.Fatalf("unexpected driver: %s", driver)
	}
}

func TestBuildBakeContentIncludesCaches(t *testing.T) {
	content := buildBakeContent("/root/lunafox", "on", "https://proxy.golang.org,direct", true, "/tmp/a", "/tmp/b", "docker.io", "yyhuni")
	if !strings.Contains(content, "cache-from") {
		t.Fatalf("expected cache-from in bake content")
	}
}

func TestResolveVersionProdRequiresVersion(t *testing.T) {
	dir := t.TempDir()
	installer := NewInstaller(cli.Options{
		Mode:           cli.ModeProd,
		Version:        "",
		DockerDir:      filepath.Join(dir, "docker"),
		ComposeFile:    filepath.Join(dir, "docker", "docker-compose.yml"),
		ImageRegistry:  "docker.io",
		ImageNamespace: "yyhuni",
	}, nil, ui.NewPrinter(io.Discard, io.Discard))

	err := installer.resolveVersion()
	if err == nil {
		t.Fatalf("expected resolveVersion to fail without explicit version")
	}
}

func TestResolveVersionProdUsesExplicitVersion(t *testing.T) {
	installer := NewInstaller(cli.Options{
		Mode:           cli.ModeProd,
		Version:        "v1.2.3",
		ImageRegistry:  "docker.io",
		ImageNamespace: "yyhuni",
	}, nil, ui.NewPrinter(io.Discard, io.Discard))

	if err := installer.resolveVersion(); err != nil {
		t.Fatalf("resolveVersion failed: %v", err)
	}
	if installer.version != "v1.2.3" {
		t.Fatalf("unexpected version: %s", installer.version)
	}
}

func TestCheckURLReadyAndWarm(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		switch request.URL.Path {
		case "/ok":
			writer.WriteHeader(http.StatusOK)
		case "/not-found":
			writer.WriteHeader(http.StatusNotFound)
		default:
			writer.WriteHeader(http.StatusInternalServerError)
		}
	}))
	defer server.Close()

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, //nolint:gosec
		},
	}

	if !checkURLReady(client, server.URL+"/ok") {
		t.Fatalf("expected /ok to be ready")
	}
	if checkURLReady(client, server.URL+"/not-found") {
		t.Fatalf("expected /not-found not to be ready")
	}
	if !checkURLWarm(client, server.URL+"/not-found") {
		t.Fatalf("expected /not-found to be warm")
	}
}

func TestResolvePublicPort(t *testing.T) {
	if got := resolvePublicPort("https://example.com:8443", ""); got != "8443" {
		t.Fatalf("unexpected port from url: %s", got)
	}
	if got := resolvePublicPort("https://example.com", ""); got != "443" {
		t.Fatalf("unexpected https default port: %s", got)
	}
	if got := resolvePublicPort("", "18083"); got != "18083" {
		t.Fatalf("unexpected preferred port: %s", got)
	}
	if got := resolvePublicPort("", ""); got != cli.DefaultPublicPort {
		t.Fatalf("unexpected fallback port: %s", got)
	}
}
