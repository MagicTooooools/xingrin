package application

import (
	"context"
	"errors"
	"time"

	agentdomain "github.com/yyhuni/lunafox/server/internal/modules/agent/domain"
	"gorm.io/gorm"
)

func (service *CommandService) CreateRegistrationToken(ctx context.Context) (*agentdomain.RegistrationToken, error) {
	now := time.Now().UTC()
	if err := service.tokenStore.DeleteExpired(ctx, now); err != nil {
		return nil, err
	}

	token, err := generateHexToken(4)
	if err != nil {
		return nil, err
	}
	expiresAt := now.Add(1 * time.Hour)
	registration := agentdomain.NewRegistrationToken(token, expiresAt)
	if err := service.tokenStore.Create(ctx, registration); err != nil {
		return nil, err
	}
	return registration, nil
}

func (service *CommandService) ValidateRegistrationToken(ctx context.Context, token string) error {
	if token == "" {
		return ErrRegistrationTokenInvalid
	}
	registration, err := service.tokenStore.FindValid(ctx, token, time.Now().UTC())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrRegistrationTokenInvalid
		}
		return err
	}
	if registration == nil {
		return ErrRegistrationTokenInvalid
	}
	return nil
}

func (service *CommandService) RegisterAgent(ctx context.Context, token, hostname, version, ipAddress string, options AgentRegistrationOptions) (*agentdomain.Agent, error) {
	if token == "" {
		return nil, ErrRegistrationTokenInvalid
	}

	registration, err := service.tokenStore.FindValid(ctx, token, time.Now().UTC())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRegistrationTokenInvalid
		}
		return nil, err
	}
	if registration == nil {
		return nil, ErrRegistrationTokenInvalid
	}

	apiKey, err := generateHexToken(4)
	if err != nil {
		return nil, err
	}

	agent := agentdomain.NewRegisteredAgent(token, hostname, version, ipAddress, apiKey, options)

	if err := service.agentStore.Create(ctx, agent); err != nil {
		return nil, err
	}
	return agent, nil
}
