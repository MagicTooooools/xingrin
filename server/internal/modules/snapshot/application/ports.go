package application

import (
	"github.com/shopspring/decimal"
	snapshotdomain "github.com/yyhuni/lunafox/server/internal/modules/snapshot/domain"
	"gorm.io/datatypes"
)

var (
	ErrSnapshotScanNotFound   = snapshotdomain.ErrSnapshotScanNotFound
	ErrSnapshotTargetMismatch = snapshotdomain.ErrSnapshotTargetMismatch
)

type SnapshotScanRefLookup interface {
	GetScanRefByID(id int) (*snapshotdomain.ScanRef, error)
	GetTargetRefByScanID(scanID int) (*snapshotdomain.ScanTargetRef, error)
}

type VulnerabilityRawOutputCodec interface {
	Encode(rawOutput map[string]any) (datatypes.JSON, error)
}

type VulnerabilitySnapshotItem struct {
	URL         string
	VulnType    string
	Severity    string
	Source      string
	CVSSScore   *decimal.Decimal
	Description string
	RawOutput   map[string]any
}

type VulnerabilityAssetCreateItem struct {
	URL         string
	VulnType    string
	Severity    string
	Source      string
	CVSSScore   *decimal.Decimal
	Description string
	RawOutput   map[string]any
}
