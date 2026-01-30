package job

import (
	"context"
	"time"

	"github.com/yyhuni/lunafox/server/internal/pkg"
	"github.com/yyhuni/lunafox/server/internal/repository"
	"go.uber.org/zap"
)

// AgentMonitor marks stale agents offline and recovers their tasks.
type AgentMonitor struct {
	agentRepo    repository.AgentRepository
	scanTaskRepo repository.ScanTaskRepository
	interval     time.Duration
	timeout      time.Duration
}

// NewAgentMonitor creates a new AgentMonitor.
func NewAgentMonitor(agentRepo repository.AgentRepository, scanTaskRepo repository.ScanTaskRepository, interval, timeout time.Duration) *AgentMonitor {
	return &AgentMonitor{
		agentRepo:    agentRepo,
		scanTaskRepo: scanTaskRepo,
		interval:     interval,
		timeout:      timeout,
	}
}

// Run starts the monitor loop.
func (m *AgentMonitor) Run(ctx context.Context) {
	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()

	m.check(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			m.check(ctx)
		}
	}
}

func (m *AgentMonitor) check(ctx context.Context) {
	cutoff := time.Now().Add(-m.timeout)

	agents, err := m.agentRepo.FindStaleOnline(ctx, cutoff)
	if err != nil {
		pkg.Warn("Failed to query stale agents", zap.Error(err))
		return
	}

	for _, agent := range agents {
		if err := m.agentRepo.UpdateStatus(ctx, agent.ID, "offline"); err != nil {
			pkg.Warn("Failed to mark agent offline", zap.Int("agent_id", agent.ID), zap.Error(err))
			continue
		}
		if err := m.scanTaskRepo.FailTasksForOfflineAgent(ctx, agent.ID); err != nil {
			pkg.Warn("Failed to fail tasks for offline agent", zap.Int("agent_id", agent.ID), zap.Error(err))
		}
	}
}
