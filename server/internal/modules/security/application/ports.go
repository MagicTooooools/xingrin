package application

import (
	securitydomain "github.com/yyhuni/lunafox/server/internal/modules/security/domain"
	"gorm.io/datatypes"
)

var (
	ErrVulnerabilityNotFound = securitydomain.ErrVulnerabilityNotFound
	ErrTargetNotFound        = securitydomain.ErrTargetNotFound
)

type VulnerabilityQueryStore interface {
	FindAll(page, pageSize int, filter, severity string, isReviewed *bool) ([]securitydomain.Vulnerability, int64, error)
	GetByID(id int) (*securitydomain.Vulnerability, error)
	FindByTargetID(targetID, page, pageSize int, filter, severity string, isReviewed *bool) ([]securitydomain.Vulnerability, int64, error)
	GetStats() (total, pending, reviewed int64, err error)
	GetStatsByTargetID(targetID int) (total, pending, reviewed int64, err error)
}

type VulnerabilityCommandStore interface {
	BulkCreate(vulnerabilities []securitydomain.Vulnerability) (int64, error)
	BulkDelete(ids []int) (int64, error)
	MarkAsReviewed(id int) error
	MarkAsUnreviewed(id int) error
	BulkMarkAsReviewed(ids []int) (int64, error)
	BulkMarkAsUnreviewed(ids []int) (int64, error)
}

type VulnerabilityStore interface {
	VulnerabilityQueryStore
	VulnerabilityCommandStore
}

type VulnerabilityTargetLookup interface {
	GetActiveByID(id int) (*securitydomain.TargetRef, error)
}

type VulnerabilityRawOutputCodec interface {
	Encode(rawOutput map[string]any) (datatypes.JSON, error)
}
