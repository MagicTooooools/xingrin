package repository

import (
	"context"
	"time"

	"github.com/yyhuni/lunafox/server/internal/model"
	"gorm.io/gorm"
)

// RegistrationTokenRepository defines data access for registration tokens.
type RegistrationTokenRepository interface {
	Create(ctx context.Context, token *model.RegistrationToken) error
	FindValid(ctx context.Context, token string, now time.Time) (*model.RegistrationToken, error)
	DeleteExpired(ctx context.Context, now time.Time) error
}

type registrationTokenRepository struct {
	db *gorm.DB
}

// NewRegistrationTokenRepository creates a new registration token repository.
func NewRegistrationTokenRepository(db *gorm.DB) RegistrationTokenRepository {
	return &registrationTokenRepository{db: db}
}

// Create inserts a new registration token.
func (r *registrationTokenRepository) Create(ctx context.Context, token *model.RegistrationToken) error {
	return r.db.WithContext(ctx).Create(token).Error
}

// FindValid returns a non-expired token by value.
func (r *registrationTokenRepository) FindValid(ctx context.Context, token string, now time.Time) (*model.RegistrationToken, error) {
	var result model.RegistrationToken
	err := r.db.WithContext(ctx).
		Where("token = ? AND expires_at > ?", token, now).
		First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteExpired removes all expired registration tokens.
func (r *registrationTokenRepository) DeleteExpired(ctx context.Context, now time.Time) error {
	return r.db.WithContext(ctx).
		Where("expires_at <= ?", now).
		Delete(&model.RegistrationToken{}).Error
}
