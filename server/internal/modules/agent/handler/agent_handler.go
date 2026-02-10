package handler

import (
	"strings"
	"text/template"

	"github.com/yyhuni/lunafox/server/internal/cache"
	agentapp "github.com/yyhuni/lunafox/server/internal/modules/agent/application"
	agentinstall "github.com/yyhuni/lunafox/server/internal/modules/agent/install"
)

// AgentHandler handles registration and admin APIs for agents.
type AgentHandler struct {
	facade         *agentapp.AgentFacade
	runtimeService *agentapp.AgentRuntimeService
	publicURL      string
	serverVersion  string
	agentImage     string
	workerToken    string
	heartbeatCache cache.HeartbeatCache
}

type installTemplateData struct {
	Token        string
	ServerURL    string
	AgentImage   string
	AgentVersion string
	WorkerToken  string
}

var agentInstallSHTemplate = template.Must(template.New("agent_install.sh").Parse(agentinstall.AgentInstallScript))

// NewAgentHandler creates a new AgentHandler.
func NewAgentHandler(
	facade *agentapp.AgentFacade,
	runtimeService *agentapp.AgentRuntimeService,
	publicURL, serverVersion, agentImage, workerToken string,
	heartbeatCache cache.HeartbeatCache,
) *AgentHandler {
	return &AgentHandler{
		facade:         facade,
		runtimeService: runtimeService,
		publicURL:      strings.TrimSpace(publicURL),
		serverVersion:  strings.TrimSpace(serverVersion),
		agentImage:     strings.TrimSpace(agentImage),
		workerToken:    strings.TrimSpace(workerToken),
		heartbeatCache: heartbeatCache,
	}
}
