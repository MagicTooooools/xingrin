package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yyhuni/lunafox/server/internal/agentproto"
	"github.com/yyhuni/lunafox/server/internal/model"
	"github.com/yyhuni/lunafox/server/internal/repository"
	"github.com/yyhuni/lunafox/server/internal/service"
	"gorm.io/gorm"
)

type configAgentRepo struct {
	agent *model.Agent
}

func (r *configAgentRepo) Create(ctx context.Context, agent *model.Agent) error {
	r.agent = agent
	return nil
}

func (r *configAgentRepo) FindByID(ctx context.Context, id int) (*model.Agent, error) {
	if r.agent != nil && r.agent.ID == id {
		return r.agent, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (r *configAgentRepo) FindByAPIKey(ctx context.Context, apiKey string) (*model.Agent, error) {
	return nil, gorm.ErrRecordNotFound
}

func (r *configAgentRepo) List(ctx context.Context, page, pageSize int, status string) ([]*model.Agent, int64, error) {
	return nil, 0, nil
}

func (r *configAgentRepo) FindStaleOnline(ctx context.Context, before time.Time) ([]*model.Agent, error) {
	return nil, nil
}

func (r *configAgentRepo) Update(ctx context.Context, agent *model.Agent) error {
	r.agent = agent
	return nil
}

func (r *configAgentRepo) UpdateStatus(ctx context.Context, id int, status string) error {
	return nil
}

func (r *configAgentRepo) UpdateHeartbeat(ctx context.Context, id int, update repository.AgentHeartbeatUpdate) error {
	return nil
}

func (r *configAgentRepo) Delete(ctx context.Context, id int) error {
	return nil
}

type configTokenRepo struct{}

func (r *configTokenRepo) Create(ctx context.Context, token *model.RegistrationToken) error {
	return nil
}

func (r *configTokenRepo) FindValid(ctx context.Context, token string, now time.Time) (*model.RegistrationToken, error) {
	return nil, gorm.ErrRecordNotFound
}

type fakeConfigNotifier struct {
	called  bool
	agentID int
	payload agentproto.ConfigUpdatePayload
}

func (f *fakeConfigNotifier) SendConfigUpdate(agentID int, payload agentproto.ConfigUpdatePayload) {
	f.called = true
	f.agentID = agentID
	f.payload = payload
}

func TestUpdateAgentConfigSendsConfigUpdate(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &configAgentRepo{agent: &model.Agent{
		ID:            1,
		MaxTasks:      5,
		CPUThreshold:  85,
		MemThreshold:  85,
		DiskThreshold: 90,
	}}
	svc := service.NewAgentService(repo, &configTokenRepo{})
	notifier := &fakeConfigNotifier{}
	handler := NewAgentHandler(svc, "", "", "", "", nil, notifier)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	body := `{"maxTasks":10,"cpuThreshold":70,"memThreshold":71,"diskThreshold":72}`
	req := httptest.NewRequest(http.MethodPut, "/api/agents/1/config", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	c.Params = []gin.Param{{Key: "id", Value: "1"}}

	handler.UpdateAgentConfig(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if !notifier.called {
		t.Fatalf("expected config notifier to be called")
	}
	if notifier.agentID != 1 {
		t.Fatalf("expected agentID 1, got %d", notifier.agentID)
	}
	if notifier.payload.MaxTasks == nil || *notifier.payload.MaxTasks != 10 {
		t.Fatalf("expected max tasks 10")
	}
	if notifier.payload.CPUThreshold == nil || *notifier.payload.CPUThreshold != 70 {
		t.Fatalf("expected cpu threshold 70")
	}
	if notifier.payload.MemThreshold == nil || *notifier.payload.MemThreshold != 71 {
		t.Fatalf("expected mem threshold 71")
	}
	if notifier.payload.DiskThreshold == nil || *notifier.payload.DiskThreshold != 72 {
		t.Fatalf("expected disk threshold 72")
	}
}
