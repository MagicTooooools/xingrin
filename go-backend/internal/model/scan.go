package model

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/datatypes"
)

// Scan represents a scan job
type Scan struct {
	ID                int            `gorm:"primaryKey;autoIncrement" json:"id"`
	TargetID          int            `gorm:"column:target_id;not null;index:idx_scan_target" json:"targetId"`
	EngineIDs         pq.Int64Array  `gorm:"column:engine_ids;type:integer[]" json:"engineIds"`
	EngineNames       datatypes.JSON `gorm:"column:engine_names;type:jsonb" json:"engineNames"`
	YamlConfiguration string         `gorm:"column:yaml_configuration;type:text" json:"yamlConfiguration"`
	ScanMode          string         `gorm:"column:scan_mode;size:10;default:'full'" json:"scanMode"`
	Status            string         `gorm:"column:status;size:20;default:'initiated';index:idx_scan_status" json:"status"`
	ResultsDir        string         `gorm:"column:results_dir;size:100" json:"resultsDir"`
	ContainerIDs      pq.StringArray `gorm:"column:container_ids;type:varchar(100)[]" json:"containerIds"`
	WorkerID          *int           `gorm:"column:worker_id" json:"workerId"`
	ErrorMessage      string         `gorm:"column:error_message;size:2000" json:"errorMessage"`
	Progress          int            `gorm:"column:progress;default:0" json:"progress"`
	CurrentStage      string         `gorm:"column:current_stage;size:50" json:"currentStage"`
	StageProgress     datatypes.JSON `gorm:"column:stage_progress;type:jsonb" json:"stageProgress"`
	CreatedAt         time.Time      `gorm:"column:created_at;autoCreateTime;index:idx_scan_created_at" json:"createdAt"`
	StoppedAt         *time.Time     `gorm:"column:stopped_at" json:"stoppedAt"`
	DeletedAt         *time.Time     `gorm:"column:deleted_at;index:idx_scan_deleted_at" json:"-"`

	// Cached statistics
	CachedSubdomainsCount  int        `gorm:"column:cached_subdomains_count;default:0" json:"cachedSubdomainsCount"`
	CachedWebsitesCount    int        `gorm:"column:cached_websites_count;default:0" json:"cachedWebsitesCount"`
	CachedEndpointsCount   int        `gorm:"column:cached_endpoints_count;default:0" json:"cachedEndpointsCount"`
	CachedIPsCount         int        `gorm:"column:cached_ips_count;default:0" json:"cachedIpsCount"`
	CachedDirectoriesCount int        `gorm:"column:cached_directories_count;default:0" json:"cachedDirectoriesCount"`
	CachedScreenshotsCount int        `gorm:"column:cached_screenshots_count;default:0" json:"cachedScreenshotsCount"`
	CachedVulnsTotal       int        `gorm:"column:cached_vulns_total;default:0" json:"cachedVulnsTotal"`
	CachedVulnsCritical    int        `gorm:"column:cached_vulns_critical;default:0" json:"cachedVulnsCritical"`
	CachedVulnsHigh        int        `gorm:"column:cached_vulns_high;default:0" json:"cachedVulnsHigh"`
	CachedVulnsMedium      int        `gorm:"column:cached_vulns_medium;default:0" json:"cachedVulnsMedium"`
	CachedVulnsLow         int        `gorm:"column:cached_vulns_low;default:0" json:"cachedVulnsLow"`
	StatsUpdatedAt         *time.Time `gorm:"column:stats_updated_at" json:"statsUpdatedAt"`

	// Relationships
	Target *Target `gorm:"foreignKey:TargetID" json:"target,omitempty"`
}

// TableName returns the table name for Scan
func (Scan) TableName() string {
	return "scan"
}

// ScanStatus constants
const (
	ScanStatusInitiated = "initiated"
	ScanStatusRunning   = "running"
	ScanStatusCompleted = "completed"
	ScanStatusFailed    = "failed"
	ScanStatusStopped   = "stopped"
	ScanStatusPending   = "pending"
)

// ScanMode constants
const (
	ScanModeFull  = "full"
	ScanModeQuick = "quick"
)
