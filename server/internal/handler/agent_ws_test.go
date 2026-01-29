package handler

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/yyhuni/orbit/server/internal/agentproto"
	"github.com/yyhuni/orbit/server/internal/cache"
	"github.com/yyhuni/orbit/server/internal/model"
	"github.com/yyhuni/orbit/server/internal/repository"
	ws "github.com/yyhuni/orbit/server/internal/websocket"
)

type fakeAgentRepo struct {
	lastUpdate repository.AgentHeartbeatUpdate
	updated    bool
}

func (f *fakeAgentRepo) Create(ctx context.Context, agent *model.Agent) error {
	return nil
}

func (f *fakeAgentRepo) FindByID(ctx context.Context, id int) (*model.Agent, error) {
	return nil, nil
}

func (f *fakeAgentRepo) FindByAPIKey(ctx context.Context, apiKey string) (*model.Agent, error) {
	return nil, nil
}

func (f *fakeAgentRepo) List(ctx context.Context, page, pageSize int, status string) ([]*model.Agent, int64, error) {
	return nil, 0, nil
}

func (f *fakeAgentRepo) FindStaleOnline(ctx context.Context, before time.Time) ([]*model.Agent, error) {
	return nil, nil
}

func (f *fakeAgentRepo) Update(ctx context.Context, agent *model.Agent) error {
	return nil
}

func (f *fakeAgentRepo) UpdateStatus(ctx context.Context, id int, status string) error {
	return nil
}

func (f *fakeAgentRepo) UpdateHeartbeat(ctx context.Context, id int, update repository.AgentHeartbeatUpdate) error {
	f.lastUpdate = update
	f.updated = true
	return nil
}

func (f *fakeAgentRepo) Delete(ctx context.Context, id int) error {
	return nil
}

type fakeHeartbeatCache struct {
	last *cache.HeartbeatData
}

func (f *fakeHeartbeatCache) Set(ctx context.Context, agentID int, data *cache.HeartbeatData) error {
	f.last = data
	return nil
}

func (f *fakeHeartbeatCache) Get(ctx context.Context, agentID int) (*cache.HeartbeatData, error) {
	return nil, nil
}

func (f *fakeHeartbeatCache) Delete(ctx context.Context, agentID int) error {
	return nil
}

func TestHandleHeartbeatUpdatesRepoAndCache(t *testing.T) {
	repo := &fakeAgentRepo{}
	cacheStore := &fakeHeartbeatCache{}
	handler := NewAgentWebSocketHandler(ws.NewHub(), repo, cacheStore, "", "")

	agent := &model.Agent{ID: 1}
	now := time.Now().UTC()
	payload := agentproto.HeartbeatPayload{
		CPU:      12.5,
		Mem:      34.5,
		Disk:     56.7,
		Tasks:    2,
		Version:  "v1.0.0",
		Hostname: "node-1",
		Uptime:   120,
		Health: &agentproto.HealthStatus{
			State: "ok",
			Since: &now,
		},
	}

	if err := handler.handleHeartbeat(context.Background(), agent, payload); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !repo.updated {
		t.Fatalf("expected repo update")
	}
	if repo.lastUpdate.Version != "v1.0.0" || repo.lastUpdate.Hostname != "node-1" {
		t.Fatalf("unexpected heartbeat update data")
	}
	if cacheStore.last == nil || cacheStore.last.Version != "v1.0.0" {
		t.Fatalf("expected cache update")
	}
	if cacheStore.last.Health == nil || cacheStore.last.Health.State != "ok" {
		t.Fatalf("expected cached health state")
	}
}

func TestUpdateRequiredSendsMessage(t *testing.T) {
	hub := ws.NewHub()
	go hub.Run()

	client := &ws.Client{
		AgentID: 1,
		Send:    make(chan []byte, 1),
		Hub:     hub,
	}
	hub.Register(client)

	handler := NewAgentWebSocketHandler(hub, &fakeAgentRepo{}, nil, "v2.0.0", "yyhuni/orbit-agent")

	handler.maybeSendUpdateRequired(context.Background(), 1, "v1.0.0")

	select {
	case msg := <-client.Send:
		if len(msg) == 0 {
			t.Fatalf("expected message payload")
		}
	default:
		t.Fatalf("expected update_required to be sent")
	}
}

func TestSendConfigUpdateSendsMessage(t *testing.T) {
	hub := ws.NewHub()
	go hub.Run()

	client := &ws.Client{
		AgentID: 1,
		Send:    make(chan []byte, 1),
		Hub:     hub,
	}
	hub.Register(client)

	for i := 0; i < 50; i++ {
		if hub.IsConnected(1) {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	handler := NewAgentWebSocketHandler(hub, &fakeAgentRepo{}, nil, "", "")
	agent := &model.Agent{
		ID:            1,
		MaxTasks:      5,
		CPUThreshold:  80,
		MemThreshold:  81,
		DiskThreshold: 82,
	}

	handler.sendConfigUpdate(agent)

	select {
	case msg := <-client.Send:
		var envelope agentproto.Message
		if err := json.Unmarshal(msg, &envelope); err != nil {
			t.Fatalf("invalid message envelope: %v", err)
		}
		if envelope.Type != agentproto.MessageTypeConfigUpdate {
			t.Fatalf("expected config_update, got %s", envelope.Type)
		}
		var payload agentproto.ConfigUpdatePayload
		if err := json.Unmarshal(envelope.Payload, &payload); err != nil {
			t.Fatalf("invalid config payload: %v", err)
		}
		if payload.MaxTasks == nil || *payload.MaxTasks != 5 {
			t.Fatalf("expected max tasks 5")
		}
		if payload.CPUThreshold == nil || *payload.CPUThreshold != 80 {
			t.Fatalf("expected cpu threshold 80")
		}
		if payload.MemThreshold == nil || *payload.MemThreshold != 81 {
			t.Fatalf("expected mem threshold 81")
		}
		if payload.DiskThreshold == nil || *payload.DiskThreshold != 82 {
			t.Fatalf("expected disk threshold 82")
		}
	default:
		t.Fatalf("expected config_update to be sent")
	}
}
