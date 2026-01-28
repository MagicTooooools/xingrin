package handler

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestInstallScriptUsesPublicURL(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/agents/install.sh?token=abc123", nil)

	h := NewAgentHandler(nil, "https://example.com", "v1.2.3", "yyhuni/orbit-agent", "worker-secret", nil, nil)
	h.InstallScript(c)

	if w.Code != 200 {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, `SERVER_URL="https://example.com"`) {
		t.Fatalf("expected server url in script")
	}
	if !strings.Contains(body, `AGENT_VERSION="v1.2.3"`) {
		t.Fatalf("expected version in script")
	}
	if !strings.Contains(body, `DEFAULT_WORKER_TOKEN="worker-secret"`) {
		t.Fatalf("expected worker token in script")
	}
}

func TestInstallScriptInfersURL(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest("GET", "/api/agents/install.sh?token=abc123", nil)
	req.Header.Set("X-Forwarded-Proto", "https")
	req.Header.Set("X-Forwarded-Host", "orbit.example.com")
	c.Request = req

	h := NewAgentHandler(nil, "", "v1.0.0", "yyhuni/orbit-agent", "worker-secret", nil, nil)
	h.InstallScript(c)

	if w.Code != 200 {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), `SERVER_URL="https://orbit.example.com"`) {
		t.Fatalf("expected inferred server url")
	}
}

func TestInstallScriptMissingToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/agents/install.sh", nil)

	h := NewAgentHandler(nil, "https://example.com", "v1.2.3", "yyhuni/orbit-agent", "worker-secret", nil, nil)
	h.InstallScript(c)

	if w.Code != 400 {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}
