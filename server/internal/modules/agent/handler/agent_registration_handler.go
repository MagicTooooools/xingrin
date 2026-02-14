package handler

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"text/template"

	"github.com/gin-gonic/gin"
	agentapp "github.com/yyhuni/lunafox/server/internal/modules/agent/application"
	agentdomain "github.com/yyhuni/lunafox/server/internal/modules/agent/domain"
	"github.com/yyhuni/lunafox/server/internal/modules/agent/dto"
	"github.com/yyhuni/lunafox/server/internal/pkg/timeutil"
)

// CreateRegistrationToken creates a new registration token.
// POST /api/admin/agents/registration-tokens
func (h *AgentHandler) CreateRegistrationToken(c *gin.Context) {
	token, err := h.facade.CreateRegistrationToken(c.Request.Context())
	if err != nil {
		dto.InternalError(c, "Failed to create registration token")
		return
	}

	dto.Created(c, dto.RegistrationTokenResponse{
		Token:     token.Token,
		ExpiresAt: timeutil.ToUTC(token.ExpiresAt),
	})
}

// Register registers an agent using a registration token.
// POST /api/agent/register
func (h *AgentHandler) Register(c *gin.Context) {
	var req dto.AgentRegistrationRequest
	if !dto.BindJSON(c, &req) {
		return
	}

	agent, err := h.facade.RegisterAgent(
		c.Request.Context(),
		req.Token,
		req.Hostname,
		req.Version,
		getForwardedIP(c),
		agentdomain.AgentRegistrationOptions{
			MaxTasks:      req.MaxTasks,
			CPUThreshold:  req.CPUThreshold,
			MemThreshold:  req.MemThreshold,
			DiskThreshold: req.DiskThreshold,
		},
	)
	if err != nil {
		if errors.Is(err, agentapp.ErrRegistrationTokenInvalid) {
			dto.BadRequest(c, "Invalid or expired registration token")
			return
		}
		dto.InternalError(c, "Failed to register agent")
		return
	}

	dto.Created(c, dto.AgentRegistrationResponse{AgentID: agent.ID, APIKey: agent.APIKey})
}

// InstallScript returns an agent installation script.
// GET /api/agent/install-script?token=...
func (h *AgentHandler) InstallScript(c *gin.Context) {
	token := strings.TrimSpace(c.Query("token"))
	if token == "" {
		dto.BadRequest(c, "Missing registration token")
		return
	}
	if h.facade != nil {
		if err := h.facade.ValidateRegistrationToken(c.Request.Context(), token); err != nil {
			if errors.Is(err, agentapp.ErrRegistrationTokenInvalid) {
				dto.BadRequest(c, "Invalid or expired registration token")
				return
			}
			dto.InternalError(c, "Failed to validate registration token")
			return
		}
	}

	serverURL := h.publicURL
	if serverURL == "" {
		serverURL = inferServerURL(c)
	}
	if serverURL == "" {
		dto.InternalError(c, "Failed to infer server URL")
		return
	}

	version := h.serverVersion
	if version == "" {
		version = "unknown"
	}
	image := h.agentImage
	if image == "" {
		image = "yyhuni/lunafox-agent"
	}

	script, err := renderInstallScript(agentInstallSHTemplate, installTemplateData{
		Token:        token,
		ServerURL:    serverURL,
		AgentImage:   image,
		AgentVersion: version,
		WorkerToken:  h.workerToken,
	})
	if err != nil {
		dto.InternalError(c, "Failed to build install script")
		return
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%q", "install.sh"))
	c.Data(http.StatusOK, "text/plain; charset=utf-8", []byte(script))
}

func renderInstallScript(tpl *template.Template, data installTemplateData) (string, error) {
	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
