package dto

import (
	"time"

	"github.com/yyhuni/lunafox/server/internal/agentproto"
)

// RegistrationTokenResponse represents a created registration token.
type RegistrationTokenResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
}

// AgentRegistrationRequest represents an agent self-registration request.
type AgentRegistrationRequest struct {
	Token    string `json:"token" binding:"required,len=8"`
	Hostname string `json:"hostname" binding:"required"`
	Version  string `json:"version" binding:"required"`
}

// HealthStatus represents agent health status.
type HealthStatus = agentproto.HealthStatus

// AgentHeartbeatResponse represents cached heartbeat data for an agent.
type AgentHeartbeatResponse struct {
	CPU       float64       `json:"cpu"`
	Mem       float64       `json:"mem"`
	Disk      float64       `json:"disk"`
	Tasks     int           `json:"tasks"`
	Uptime    int64         `json:"uptime"`
	UpdatedAt time.Time     `json:"updatedAt"`
	Health    *HealthStatus `json:"health,omitempty"`
}

// AgentResponse represents an agent in management APIs.
type AgentResponse struct {
	ID            int                     `json:"id"`
	Name          string                  `json:"name"`
	Status        string                  `json:"status"`
	Hostname      string                  `json:"hostname,omitempty"`
	IPAddress     string                  `json:"ipAddress,omitempty"`
	Version       string                  `json:"version,omitempty"`
	MaxTasks      int                     `json:"maxTasks"`
	CPUThreshold  int                     `json:"cpuThreshold"`
	MemThreshold  int                     `json:"memThreshold"`
	DiskThreshold int                     `json:"diskThreshold"`
	ConnectedAt   *time.Time              `json:"connectedAt,omitempty"`
	LastHeartbeat *time.Time              `json:"lastHeartbeat,omitempty"`
	Health        HealthStatus            `json:"health"`
	Heartbeat     *AgentHeartbeatResponse `json:"heartbeat,omitempty"`
	CreatedAt     time.Time               `json:"createdAt"`
}

// AgentListResponse represents a paginated list of agents.
type AgentListResponse struct {
	Results  []AgentResponse `json:"results"`
	Total    int64           `json:"total"`
	Page     int             `json:"page"`
	PageSize int             `json:"pageSize"`
}

// UpdateAgentConfigRequest represents a config update request.
type UpdateAgentConfigRequest struct {
	MaxTasks      *int `json:"maxTasks" binding:"omitempty,min=1,max=20"`
	CPUThreshold  *int `json:"cpuThreshold" binding:"omitempty,min=1,max=100"`
	MemThreshold  *int `json:"memThreshold" binding:"omitempty,min=1,max=100"`
	DiskThreshold *int `json:"diskThreshold" binding:"omitempty,min=1,max=100"`
}

