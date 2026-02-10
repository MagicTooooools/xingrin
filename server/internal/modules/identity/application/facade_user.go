package application

import (
	"context"
	"errors"

	"github.com/yyhuni/lunafox/server/internal/auth"
	identitydomain "github.com/yyhuni/lunafox/server/internal/modules/identity/domain"
	"github.com/yyhuni/lunafox/server/internal/modules/identity/dto"
	"gorm.io/gorm"
)

// UserFacade handles user business logic.
type UserFacade struct {
	queryService *UserQueryService
	cmdService   *UserCommandService
}

// NewUserFacade creates a new user service.
func NewUserFacade(store UserStore) *UserFacade {
	return &UserFacade{
		queryService: NewUserQueryService(store),
		cmdService:   NewUserCommandService(store, authPasswordHasher{}),
	}
}

// Create creates a new user.
func (service *UserFacade) Create(req *dto.CreateUserRequest) (*identitydomain.User, error) {
	user, err := service.cmdService.CreateUser(context.Background(), req.Username, req.Password, req.Email)
	if err != nil {
		if errors.Is(err, ErrUsernameExists) {
			return nil, ErrUsernameExists
		}
		return nil, err
	}
	return user, nil
}

// List returns paginated users.
func (service *UserFacade) List(query *dto.PaginationQuery) ([]identitydomain.User, int64, error) {
	users, total, err := service.queryService.ListUsers(context.Background(), query.GetPage(), query.GetPageSize())
	if err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

// GetByID returns a user by ID.
func (service *UserFacade) GetByID(id int) (*identitydomain.User, error) {
	user, err := service.queryService.GetUserByID(context.Background(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

// UpdatePassword updates user password.
func (service *UserFacade) UpdatePassword(id int, req *dto.UpdatePasswordRequest) error {
	err := service.cmdService.UpdateUserPassword(context.Background(), id, req.OldPassword, req.NewPassword)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		if errors.Is(err, ErrInvalidPassword) {
			return ErrInvalidPassword
		}
		return err
	}
	return nil
}

type authPasswordHasher struct{}

func (authPasswordHasher) HashPassword(password string) (string, error) {
	return auth.HashPassword(password)
}

func (authPasswordHasher) VerifyPassword(password, hashed string) bool {
	return auth.VerifyPassword(password, hashed)
}
