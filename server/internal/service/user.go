package service

import (
	"errors"

	"github.com/xingrin/server/internal/auth"
	"github.com/xingrin/server/internal/dto"
	"github.com/xingrin/server/internal/model"
	"github.com/xingrin/server/internal/repository"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrUsernameExists  = errors.New("username already exists")
	ErrInvalidPassword = errors.New("invalid password")
)

// UserService handles user business logic
type UserService struct {
	repo *repository.UserRepository
}

// NewUserService creates a new user service
func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// Create creates a new user
func (s *UserService) Create(req *dto.CreateUserRequest) (*model.User, error) {
	exists, err := s.repo.ExistsByUsername(req.Username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrUsernameExists
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Username: req.Username,
		Password: hashedPassword,
		Email:    req.Email,
		IsActive: true,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

// List returns paginated users
func (s *UserService) List(query *dto.PaginationQuery) ([]model.User, int64, error) {
	return s.repo.FindAll(query.GetPage(), query.GetPageSize())
}

// GetByID returns a user by ID
func (s *UserService) GetByID(id int) (*model.User, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

// UpdatePassword updates user password
func (s *UserService) UpdatePassword(id int, req *dto.UpdatePasswordRequest) error {
	user, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return err
	}

	if !auth.VerifyPassword(req.OldPassword, user.Password) {
		return ErrInvalidPassword
	}

	hashedPassword, err := auth.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	user.Password = hashedPassword
	return s.repo.Update(user)
}
