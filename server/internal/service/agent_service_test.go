package service

import (
	"context"
	"testing"
	"time"

	"github.com/yyhuni/orbit/server/internal/model"
	"github.com/yyhuni/orbit/server/internal/repository"
	"gorm.io/gorm"
)

type fakeAgentRepo struct {
	created *model.Agent
}

func (f *fakeAgentRepo) Create(ctx context.Context, agent *model.Agent) error {
	f.created = agent
	return nil
}

func (f *fakeAgentRepo) FindByID(ctx context.Context, id int) (*model.Agent, error) {
	return nil, gorm.ErrRecordNotFound
}

func (f *fakeAgentRepo) FindByAPIKey(ctx context.Context, apiKey string) (*model.Agent, error) {
	return nil, gorm.ErrRecordNotFound
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
	return nil
}

func (f *fakeAgentRepo) Delete(ctx context.Context, id int) error {
	return nil
}

type fakeTokenRepo struct {
	token *model.RegistrationToken
	err   error
}

func (f *fakeTokenRepo) Create(ctx context.Context, token *model.RegistrationToken) error {
	f.token = token
	return f.err
}

func (f *fakeTokenRepo) FindValid(ctx context.Context, token string, now time.Time) (*model.RegistrationToken, error) {
	if f.err != nil {
		return nil, f.err
	}
	if f.token == nil || f.token.Token != token {
		return nil, gorm.ErrRecordNotFound
	}
	return f.token, nil
}

func TestCreateRegistrationToken(t *testing.T) {
	agentRepo := &fakeAgentRepo{}
	tokenRepo := &fakeTokenRepo{}
	svc := NewAgentService(agentRepo, tokenRepo)

	token, err := svc.CreateRegistrationToken(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(token.Token) != 8 {
		t.Fatalf("expected 8 hex token, got %s", token.Token)
	}
	if token.ExpiresAt.Before(time.Now()) {
		t.Fatalf("expected future expiration")
	}
}

func TestRegisterAgentInvalidToken(t *testing.T) {
	agentRepo := &fakeAgentRepo{}
	tokenRepo := &fakeTokenRepo{}
	svc := NewAgentService(agentRepo, tokenRepo)

	_, err := svc.RegisterAgent(context.Background(), "badtoken", "host", "v1", "127.0.0.1")
	if err != ErrRegistrationTokenInvalid {
		t.Fatalf("expected invalid token error, got %v", err)
	}
}

func TestRegisterAgentCreatesOfflineAgent(t *testing.T) {
	agentRepo := &fakeAgentRepo{}
	tokenRepo := &fakeTokenRepo{
		token: &model.RegistrationToken{
			Token: "abcd1234",
		},
	}
	svc := NewAgentService(agentRepo, tokenRepo)

	agent, err := svc.RegisterAgent(context.Background(), "abcd1234", "host1", "v1", "127.0.0.1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if agent.Status != "offline" {
		t.Fatalf("expected offline status")
	}
	if len(agent.APIKey) != 8 {
		t.Fatalf("expected api key length 8")
	}
	if agentRepo.created == nil {
		t.Fatalf("expected agent to be created")
	}
}

func TestUpdateAgentConfig(t *testing.T) {
	agent := &model.Agent{
		ID:            1,
		MaxTasks:      5,
		CPUThreshold:  80,
		MemThreshold:  81,
		DiskThreshold: 82,
	}
	repo := &statefulAgentRepo{agent: agent}
	svc := NewAgentService(repo, &fakeTokenRepo{})

	maxTasks := 9
	cpu := 70
	update := AgentConfigUpdate{
		MaxTasks:     &maxTasks,
		CPUThreshold: &cpu,
	}

	updated, err := svc.UpdateAgentConfig(context.Background(), 1, update)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.MaxTasks != 9 || updated.CPUThreshold != 70 {
		t.Fatalf("expected updated config values")
	}
	if updated.MemThreshold != 81 || updated.DiskThreshold != 82 {
		t.Fatalf("expected untouched values to remain")
	}
	if repo.updated == nil {
		t.Fatalf("expected repository update")
	}
}

func TestRegenerateAPIKey(t *testing.T) {
	agent := &model.Agent{
		ID:     1,
		APIKey: "deadbeef",
	}
	repo := &statefulAgentRepo{agent: agent}
	svc := NewAgentService(repo, &fakeTokenRepo{})

	apiKey, err := svc.RegenerateAPIKey(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(apiKey) != 8 {
		t.Fatalf("expected 8 hex api key")
	}
	if repo.updated == nil || repo.updated.APIKey != apiKey {
		t.Fatalf("expected updated api key to persist")
	}
}

func TestDeleteAgentNotFound(t *testing.T) {
	repo := &statefulAgentRepo{deleteErr: gorm.ErrRecordNotFound}
	svc := NewAgentService(repo, &fakeTokenRepo{})

	if err := svc.DeleteAgent(context.Background(), 1); err != ErrAgentNotFound {
		t.Fatalf("expected ErrAgentNotFound, got %v", err)
	}
}

type statefulAgentRepo struct {
	agent     *model.Agent
	updated   *model.Agent
	deleteErr error
}

func (s *statefulAgentRepo) Create(ctx context.Context, agent *model.Agent) error {
	return nil
}

func (s *statefulAgentRepo) FindByID(ctx context.Context, id int) (*model.Agent, error) {
	if s.agent == nil {
		return nil, gorm.ErrRecordNotFound
	}
	return s.agent, nil
}

func (s *statefulAgentRepo) FindByAPIKey(ctx context.Context, apiKey string) (*model.Agent, error) {
	return nil, gorm.ErrRecordNotFound
}

func (s *statefulAgentRepo) List(ctx context.Context, page, pageSize int, status string) ([]*model.Agent, int64, error) {
	return nil, 0, nil
}

func (s *statefulAgentRepo) FindStaleOnline(ctx context.Context, before time.Time) ([]*model.Agent, error) {
	return nil, nil
}

func (s *statefulAgentRepo) Update(ctx context.Context, agent *model.Agent) error {
	s.updated = agent
	return nil
}

func (s *statefulAgentRepo) UpdateStatus(ctx context.Context, id int, status string) error {
	return nil
}

func (s *statefulAgentRepo) UpdateHeartbeat(ctx context.Context, id int, update repository.AgentHeartbeatUpdate) error {
	return nil
}

func (s *statefulAgentRepo) Delete(ctx context.Context, id int) error {
	return s.deleteErr
}
