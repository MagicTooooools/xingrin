package application

import (
	"io"
	"time"

	catalogdomain "github.com/yyhuni/lunafox/server/internal/modules/catalog/domain"
)

var (
	ErrEngineNotFound = catalogdomain.ErrEngineNotFound
	ErrEngineExists   = catalogdomain.ErrEngineExists
	ErrInvalidEngine  = catalogdomain.ErrInvalidEngine

	ErrTargetNotFound       = catalogdomain.ErrTargetNotFound
	ErrTargetExists         = catalogdomain.ErrTargetExists
	ErrInvalidTarget        = catalogdomain.ErrInvalidTarget
	ErrTargetOrgNotFound    = catalogdomain.ErrTargetOrgNotFound
	ErrTargetOrgBindingFail = catalogdomain.ErrTargetOrgBindingFail

	ErrWordlistNotFound = catalogdomain.ErrWordlistNotFound
	ErrWordlistExists   = catalogdomain.ErrWordlistExists
	ErrEmptyName        = catalogdomain.ErrWordlistNameEmpty
	ErrNameTooLong      = catalogdomain.ErrWordlistNameTooLong
	ErrInvalidName      = catalogdomain.ErrWordlistNameInvalid
	ErrFileNotFound     = catalogdomain.ErrWordlistFileNotFound
	ErrInvalidFileType  = catalogdomain.ErrWordlistInvalidFileType
)

type EngineQueryStore interface {
	GetByID(id int) (*catalogdomain.ScanEngine, error)
	FindAll(page, pageSize int) ([]catalogdomain.ScanEngine, int64, error)
}

type EngineCommandStore interface {
	GetByID(id int) (*catalogdomain.ScanEngine, error)
	ExistsByName(name string, excludeID ...int) (bool, error)
	Create(engine *catalogdomain.ScanEngine) error
	Update(engine *catalogdomain.ScanEngine) error
	Delete(id int) error
}

type TargetQueryStore interface {
	GetActiveByID(id int) (*catalogdomain.Target, error)
	FindAll(page, pageSize int, targetType, filter string) ([]catalogdomain.Target, int64, error)
	GetAssetCounts(targetID int) (*catalogdomain.TargetAssetCounts, error)
	GetVulnerabilityCounts(targetID int) (*catalogdomain.VulnerabilityCounts, error)
}

type TargetCommandStore interface {
	GetActiveByID(id int) (*catalogdomain.Target, error)
	ExistsByName(name string, excludeID ...int) (bool, error)
	Create(target *catalogdomain.Target) error
	Update(target *catalogdomain.Target) error
	SoftDelete(id int) error
	BulkSoftDelete(ids []int) (int64, error)
	BulkCreateIgnoreConflicts(targets []catalogdomain.Target) (int, error)
	FindByNames(names []string) ([]catalogdomain.Target, error)
}

type OrganizationTargetBindingStore interface {
	ExistsByID(id int) (bool, error)
	BulkAddTargets(organizationID int, targetIDs []int) error
}

type WordlistFileMetadata struct {
	FilePath  string
	FileSize  int64
	LineCount int
	FileHash  string
}

type WordlistFileStore interface {
	Save(basePath, filename string, content io.Reader) (*WordlistFileMetadata, error)
	Write(path, content string) (*WordlistFileMetadata, error)
	Read(path string) (string, error)
	Remove(path string) error
	Exists(path string) bool
	RefreshMetadata(path string, knownSize int64, knownUpdatedAt time.Time) (*WordlistFileMetadata, bool, error)
}

type WordlistQueryStore interface {
	FindAll(page, pageSize int) ([]catalogdomain.Wordlist, int64, error)
	List() ([]catalogdomain.Wordlist, error)
	GetByID(id int) (*catalogdomain.Wordlist, error)
	FindByName(name string) (*catalogdomain.Wordlist, error)
	Update(wordlist *catalogdomain.Wordlist) error
}

type WordlistCommandStore interface {
	GetByID(id int) (*catalogdomain.Wordlist, error)
	ExistsByName(name string) (bool, error)
	Create(wordlist *catalogdomain.Wordlist) error
	Update(wordlist *catalogdomain.Wordlist) error
	Delete(id int) error
}
