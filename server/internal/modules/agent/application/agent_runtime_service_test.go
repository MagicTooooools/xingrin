package application

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/yyhuni/lunafox/server/internal/agentproto"
	"github.com/yyhuni/lunafox/server/internal/cache"
	agentdomain "github.com/yyhuni/lunafox/server/internal/modules/agent/domain"
)

type runtimeRepoStub struct {
	agentRepoStub
	heartbeats []agentdomain.AgentHeartbeatUpdate
	statuses   []string
}

func (repo *runtimeRepoStub) UpdateHeartbeat(_ context.Context, _ int, update agentdomain.AgentHeartbeatUpdate) error {
	repo.heartbeats = append(repo.heartbeats, update)
	return nil
}
func (repo *runtimeRepoStub) UpdateStatus(_ context.Context, _ int, status string) error {
	repo.statuses = append(repo.statuses, status)
	return nil
}

type cacheStub struct {
	setCalled bool
	setErr    error
	deleted   bool
}

func (cacheStore *cacheStub) Set(_ context.Context, _ int, _ *cache.HeartbeatData) error {
	cacheStore.setCalled = true
	return cacheStore.setErr
}
func (cacheStore *cacheStub) Get(_ context.Context, _ int) (*cache.HeartbeatData, error) {
	return nil, nil
}
func (cacheStore *cacheStub) Delete(_ context.Context, _ int) error {
	cacheStore.deleted = true
	return nil
}

type publisherStub struct {
	configSent        bool
	updateSent        bool
	updateSendSuccess bool
}

func (publisher *publisherStub) SendConfigUpdate(int, agentproto.ConfigUpdatePayload) {
	publisher.configSent = true
}
func (publisher *publisherStub) SendUpdateRequired(int, agentproto.UpdateRequiredPayload) bool {
	publisher.updateSent = true
	return publisher.updateSendSuccess
}
func (publisher *publisherStub) SendTaskCancel(int, int) {}

func TestAgentRuntimeServiceHeartbeatAndUpdateRequired(t *testing.T) {
	repo := &runtimeRepoStub{}
	cacheStore := &cacheStub{}
	publisher := &publisherStub{updateSendSuccess: true}
	service := NewAgentRuntimeService(
		repo,
		cacheStore,
		publisher,
		fixedClock{now: time.Date(2026, 1, 2, 10, 0, 0, 0, time.UTC)},
		"2.0.0",
		"img",
	)

	payload, _ := json.Marshal(agentproto.Message{
		Type: agentproto.MessageTypeHeartbeat,
		Payload: func() json.RawMessage {
			p, _ := json.Marshal(agentproto.HeartbeatPayload{Version: "1.0.0", Hostname: "node1", CPU: 1, Mem: 2, Disk: 3, Tasks: 1, Uptime: 10})
			return p
		}(),
		Timestamp: time.Now().UTC(),
	})

	err := service.HandleMessage(context.Background(), &agentdomain.Agent{ID: 1}, payload)
	if err != nil {
		t.Fatalf("HandleMessage error: %v", err)
	}
	if len(repo.heartbeats) != 1 {
		t.Fatalf("expected 1 heartbeat update")
	}
	if !cacheStore.setCalled {
		t.Fatalf("expected heartbeat cache set")
	}
	if !publisher.updateSent {
		t.Fatalf("expected update_required notification")
	}
}

func TestAgentRuntimeServiceOnDisconnected(t *testing.T) {
	repo := &runtimeRepoStub{}
	cacheStore := &cacheStub{}
	service := NewAgentRuntimeService(repo, cacheStore, &publisherStub{}, fixedClock{now: time.Now().UTC()}, "", "")
	if err := service.OnDisconnected(context.Background(), 1); err != nil {
		t.Fatalf("OnDisconnected error: %v", err)
	}
	if len(repo.statuses) != 1 || repo.statuses[0] != "offline" {
		t.Fatalf("expected offline status update")
	}
	if !cacheStore.deleted {
		t.Fatalf("expected cache deletion")
	}
}

func TestAgentRuntimeServiceCacheFailureNonBlocking(t *testing.T) {
	repo := &runtimeRepoStub{}
	cacheStore := &cacheStub{setErr: errors.New("boom")}
	service := NewAgentRuntimeService(repo, cacheStore, &publisherStub{}, fixedClock{now: time.Now().UTC()}, "", "")
	payload, _ := json.Marshal(agentproto.Message{
		Type: agentproto.MessageTypeHeartbeat,
		Payload: func() json.RawMessage {
			p, _ := json.Marshal(agentproto.HeartbeatPayload{Version: "1.0.0", Hostname: "node1"})
			return p
		}(),
	})

	if err := service.HandleMessage(context.Background(), &agentdomain.Agent{ID: 1}, payload); err != nil {
		t.Fatalf("expected non-blocking cache write, got %v", err)
	}
}
