package application

import (
	"context"
	"time"

	"github.com/yyhuni/lunafox/server/internal/agentproto"
	"github.com/yyhuni/lunafox/server/internal/cache"
	agentdomain "github.com/yyhuni/lunafox/server/internal/modules/agent/domain"
	scandto "github.com/yyhuni/lunafox/server/internal/modules/scan/dto"
)

type AgentQueryStore interface {
	GetByID(ctx context.Context, id int) (*agentdomain.Agent, error)
	List(ctx context.Context, page, pageSize int, status string) ([]*agentdomain.Agent, int64, error)
}

type AgentCommandStore interface {
	Create(ctx context.Context, agent *agentdomain.Agent) error
	GetByID(ctx context.Context, id int) (*agentdomain.Agent, error)
	Update(ctx context.Context, agent *agentdomain.Agent) error
	Delete(ctx context.Context, id int) error
}

type AgentStore interface {
	AgentQueryStore
	AgentCommandStore
}

type RegistrationTokenStore interface {
	Create(ctx context.Context, token *agentdomain.RegistrationToken) error
	FindValid(ctx context.Context, token string, now time.Time) (*agentdomain.RegistrationToken, error)
	DeleteExpired(ctx context.Context, now time.Time) error
}

// HeartbeatCachePort defines heartbeat cache operations required by runtime service.
type HeartbeatCachePort interface {
	Set(ctx context.Context, agentID int, data *cache.HeartbeatData) error
	Get(ctx context.Context, agentID int) (*cache.HeartbeatData, error)
	Delete(ctx context.Context, agentID int) error
}

// AgentMessagePublisher emits typed messages to runtime connections.
type AgentMessagePublisher interface {
	SendConfigUpdate(agentID int, payload agentproto.ConfigUpdatePayload)
	SendUpdateRequired(agentID int, payload agentproto.UpdateRequiredPayload) bool
	SendTaskCancel(agentID, taskID int)
}

// Clock provides deterministic time source for services.
type Clock interface {
	NowUTC() time.Time
}

// TokenGenerator provides random token generation.
type TokenGenerator interface {
	GenerateHex(byteLen int) (string, error)
}

// ScanTaskRuntimePort describes the scan task runtime dependency for agent task service.
type ScanTaskRuntimePort interface {
	PullTask(ctx context.Context, agentID int) (*scandto.TaskAssignment, error)
	UpdateStatus(ctx context.Context, agentID, taskID int, status, errorMessage string) error
}
