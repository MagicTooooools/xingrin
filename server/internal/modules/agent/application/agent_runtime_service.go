package application

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/yyhuni/lunafox/server/internal/agentproto"
	"github.com/yyhuni/lunafox/server/internal/cache"
	agentdomain "github.com/yyhuni/lunafox/server/internal/modules/agent/domain"
	"github.com/yyhuni/lunafox/server/internal/pkg"
	"go.uber.org/zap"
)

// AgentRuntimeService orchestrates WebSocket runtime events.
type AgentRuntimeService struct {
	agentRepo       agentdomain.AgentRepository
	heartbeatCache  HeartbeatCachePort
	messageBus      AgentMessagePublisher
	clock           Clock
	serverVersion   string
	agentImage      string
	notifyMu        sync.Mutex
	notifiedVersion map[int]bool
}

func NewAgentRuntimeService(
	agentRepo agentdomain.AgentRepository,
	heartbeatCache HeartbeatCachePort,
	messageBus AgentMessagePublisher,
	clock Clock,
	serverVersion, agentImage string,
) *AgentRuntimeService {
	if clock == nil {
		panic("clock is required")
	}

	return &AgentRuntimeService{
		agentRepo:       agentRepo,
		heartbeatCache:  heartbeatCache,
		messageBus:      messageBus,
		clock:           clock,
		serverVersion:   serverVersion,
		agentImage:      agentImage,
		notifiedVersion: map[int]bool{},
	}
}

func (service *AgentRuntimeService) OnConnected(ctx context.Context, agent *agentdomain.Agent, ipAddress string) error {
	now := service.clock.NowUTC()
	agent.Status = "online"
	agent.ConnectedAt = &now
	agent.LastHeartbeat = &now
	agent.IPAddress = ipAddress

	if err := service.agentRepo.Update(ctx, agent); err != nil {
		return err
	}

	service.SendConfigUpdate(agent)
	return nil
}

func (service *AgentRuntimeService) OnDisconnected(ctx context.Context, agentID int) error {
	if err := service.agentRepo.UpdateStatus(ctx, agentID, "offline"); err != nil {
		return err
	}
	if service.heartbeatCache != nil {
		if err := service.heartbeatCache.Delete(ctx, agentID); err != nil {
			pkg.Warn("Failed to clear heartbeat cache on disconnect", zap.Int("agent_id", agentID), zap.Error(err))
		}
	}
	return nil
}

func (service *AgentRuntimeService) SendConfigUpdate(agent *agentdomain.Agent) {
	if service == nil || service.messageBus == nil || agent == nil {
		return
	}
	service.messageBus.SendConfigUpdate(agent.ID, BuildConfigUpdatePayload(agent))
}

func (service *AgentRuntimeService) HandleMessage(ctx context.Context, agent *agentdomain.Agent, raw []byte) error {
	var msg agentproto.Message
	if err := json.Unmarshal(raw, &msg); err != nil {
		return err
	}

	switch msg.Type {
	case agentproto.MessageTypeHeartbeat:
		var payload agentproto.HeartbeatPayload
		if err := json.Unmarshal(msg.Payload, &payload); err != nil {
			return err
		}
		return service.handleHeartbeat(ctx, agent.ID, payload)
	default:
		return nil
	}
}

func (service *AgentRuntimeService) handleHeartbeat(ctx context.Context, agentID int, payload agentproto.HeartbeatPayload) error {
	now := service.clock.NowUTC()
	update := agentdomain.AgentHeartbeatUpdate{
		LastHeartbeat: now,
		Version:       payload.Version,
		Hostname:      payload.Hostname,
	}

	if payload.Health != nil {
		update.HasHealth = true
		update.HealthState = payload.Health.State
		update.HealthReason = payload.Health.Reason
		update.HealthMessage = payload.Health.Message
		if payload.Health.Since != nil {
			since := payload.Health.Since.UTC()
			update.HealthSince = &since
		}
	}

	if err := service.agentRepo.UpdateHeartbeat(ctx, agentID, update); err != nil {
		return err
	}

	if service.heartbeatCache != nil {
		cachePayload := &cache.HeartbeatData{
			CPU:      payload.CPU,
			Mem:      payload.Mem,
			Disk:     payload.Disk,
			Tasks:    payload.Tasks,
			Version:  payload.Version,
			Hostname: payload.Hostname,
			Uptime:   payload.Uptime,
		}
		if payload.Health != nil {
			var since *time.Time
			if payload.Health.Since != nil {
				value := payload.Health.Since.UTC()
				since = &value
			}
			cachePayload.Health = &cache.HealthStatus{
				State:   payload.Health.State,
				Reason:  payload.Health.Reason,
				Message: payload.Health.Message,
				Since:   since,
			}
		}
		if err := service.heartbeatCache.Set(ctx, agentID, cachePayload); err != nil {
			pkg.Warn("Failed to cache heartbeat", zap.Error(err))
		}
	}

	service.maybeSendUpdateRequired(agentID, payload.Version)
	return nil
}

func (service *AgentRuntimeService) maybeSendUpdateRequired(agentID int, agentVersion string) {
	if service.messageBus == nil || service.serverVersion == "" || agentVersion == "" {
		return
	}
	if agentVersion == service.serverVersion {
		service.setNotified(agentID, false)
		return
	}
	if service.isNotified(agentID) {
		return
	}

	payload := agentproto.UpdateRequiredPayload{Version: service.serverVersion, Image: service.agentImage}
	if service.messageBus.SendUpdateRequired(agentID, payload) {
		service.setNotified(agentID, true)
	}
}

func (service *AgentRuntimeService) isNotified(agentID int) bool {
	service.notifyMu.Lock()
	defer service.notifyMu.Unlock()
	return service.notifiedVersion[agentID]
}

func (service *AgentRuntimeService) setNotified(agentID int, value bool) {
	service.notifyMu.Lock()
	defer service.notifyMu.Unlock()
	if !value {
		delete(service.notifiedVersion, agentID)
		return
	}
	service.notifiedVersion[agentID] = true
}
