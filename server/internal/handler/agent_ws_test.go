package handler

import (
	"context"
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
	handler := NewAgentWebSocketHandler(ws.NewHub(), repo, cacheStore, nil, "", "")

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

type fakeCanceller struct {
	called bool
	agent  int
}

func (f *fakeCanceller) CancelRunningTasksForAgent(ctx context.Context, agentID int) error {
	f.called = true
	f.agent = agentID
	return nil
}

func TestUpdateRequiredSendsMessageAndCancels(t *testing.T) {
	hub := ws.NewHub()
	go hub.Run()

	client := &ws.Client{
		AgentID: 1,
		Send:    make(chan []byte, 1),
		Hub:     hub,
	}
	hub.Register(client)

	canceller := &fakeCanceller{}
	handler := NewAgentWebSocketHandler(hub, &fakeAgentRepo{}, nil, canceller, "v2.0.0", "yyhuni/orbit-agent")

	handler.maybeSendUpdateRequired(context.Background(), 1, "v1.0.0")

	select {
	case msg := <-client.Send:
		if len(msg) == 0 {
			t.Fatalf("expected message payload")
		}
	default:
		t.Fatalf("expected update_required to be sent")
	}

	if !canceller.called || canceller.agent != 1 {
		t.Fatalf("expected canceller to be invoked")
	}
}
