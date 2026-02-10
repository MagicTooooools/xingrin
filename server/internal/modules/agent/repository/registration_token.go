package repository

import (
	"context"
	"time"

	agentdomain "github.com/yyhuni/lunafox/server/internal/modules/agent/domain"
	"gorm.io/gorm"
)

// RegistrationTokenRepository defines data access for registration tokens.
type RegistrationTokenRepository interface {
	Create(ctx context.Context, token *agentdomain.RegistrationToken) error
	FindValid(ctx context.Context, token string, now time.Time) (*agentdomain.RegistrationToken, error)
	DeleteExpired(ctx context.Context, now time.Time) error
}

type registrationTokenRepository struct {
	db *gorm.DB
}

// NewRegistrationTokenRepository creates a new registration token repository.
func NewRegistrationTokenRepository(db *gorm.DB) RegistrationTokenRepository {
	return &registrationTokenRepository{db: db}
}
