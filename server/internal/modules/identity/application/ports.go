package application

import (
	"errors"

	"github.com/yyhuni/lunafox/server/internal/auth"
	identitydomain "github.com/yyhuni/lunafox/server/internal/modules/identity/domain"
)

var (
	ErrUserNotFound    = identitydomain.ErrUserNotFound
	ErrUsernameExists  = identitydomain.ErrUsernameExists
	ErrInvalidPassword = identitydomain.ErrInvalidPassword

	ErrOrganizationNotFound = identitydomain.ErrOrganizationNotFound
	ErrOrganizationExists   = identitydomain.ErrOrganizationNameExist
	ErrTargetNotFound       = identitydomain.ErrTargetNotFound

	ErrInvalidCredentials  = errors.New("invalid username or password")
	ErrUserDisabled        = errors.New("user account is disabled")
	ErrInvalidRefreshToken = errors.New("invalid or expired refresh token")
)

type UserQueryStore interface {
	GetUserByID(id int) (*identitydomain.User, error)
	FindAll(page, pageSize int) ([]identitydomain.User, int64, error)
}

type UserCommandStore interface {
	GetUserByID(id int) (*identitydomain.User, error)
	ExistsByUsername(username string) (bool, error)
	Create(user *identitydomain.User) error
	Update(user *identitydomain.User) error
}

type OrganizationQueryStore interface {
	GetActiveByID(id int) (*identitydomain.Organization, error)
	FindByIDWithCount(id int) (*identitydomain.OrganizationWithTargetCount, error)
	FindAll(page, pageSize int, filter string) ([]identitydomain.OrganizationWithTargetCount, int64, error)
	FindTargets(organizationID int, page, pageSize int, targetType, filter string) ([]identitydomain.OrganizationTargetRef, int64, error)
}

type OrganizationCommandStore interface {
	GetActiveByID(id int) (*identitydomain.Organization, error)
	ExistsByName(name string, excludeID ...int) (bool, error)
	Create(org *identitydomain.Organization) error
	Update(org *identitydomain.Organization) error
	SoftDelete(id int) error
	BulkSoftDelete(ids []int) (int64, error)
	BulkAddTargets(organizationID int, targetIDs []int) error
	UnlinkTargets(organizationID int, targetIDs []int) (int64, error)
}

type AuthUserStore interface {
	GetAuthUserByID(id int) (*identitydomain.User, error)
	FindAuthUserByUsername(username string) (*identitydomain.User, error)
}

type PasswordVerifier interface {
	VerifyPassword(password, hashed string) bool
}

type TokenProvider interface {
	GenerateTokenPair(userID int, username string) (*auth.TokenPair, error)
	ValidateToken(token string) (*auth.Claims, error)
	GenerateAccessToken(userID int, username string) (string, int64, error)
}

type User = identitydomain.User
type Organization = identitydomain.Organization

type OrganizationTargetRef = identitydomain.OrganizationTargetRef
