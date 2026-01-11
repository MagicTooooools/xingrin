package model

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/datatypes"
)

// Scan represents a scan job
type Scan struct {
	ID                int            `gorm:"primaryKey;autoIncrement" json:"id"`
	TargetID          int            `gorm:"column:target_id;not null" json:"targetId"`
	EngineIDs         pq.Int64Array  `gorm:"column:engine_ids;type:integer[]" json:"engineIds"`
	EngineNames       datatypes.JSON `gorm:"column:engine_names;type:jsonb" json:"engineNames"`
	YamlConfiguration string         `gorm:"column:yaml_configuration;type:text" json:"yamlConfiguration"`
	ScanMode          string         `gorm:"column:scan_mode;size:10;default:'full'" json:"scanMode"`
	Status            string         `gorm:"column:status;size:20;default:'initiated'" json:"status"`
	ResultsDir        string         `gorm:"column:results_dir;size:100" json:"resultsDir"`
	ContainerIDs      pq.StringArray `gorm:"column:container_ids;type:varchar(100)[]" json:"containerIds"`
	WorkerID          *int           `gorm:"column:worker_id" json:"workerId"`
	ErrorMessage      string         `gorm:"column:error_message;size:2000" json:"errorMessage"`
	Progress          int            `gorm:"column:progress;default:0" json:"progress"`
	CurrentStage      string         `gorm:"column:current_stage;size:50" json:"currentStage"`
	StageProgress     datatypes.JSON `gorm:"column:stage_progress;type:jsonb" json:"stageProgress"`
	CreatedAt         time.Time      `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	StoppedAt         *time.Time     `gorm:"column:stopped_at" json:"stoppedAt"`
	DeletedAt         *time.Time     `gorm:"column:deleted_at;index" json:"-"`

	// Cached statistics
	CachedSubdomainsCount   int        `gorm:"column:cached_subdomains_count" json:"cachedSubdomainsCount"`
	CachedWebsitesCount     int        `gorm:"column:cached_websites_count" json:"cachedWebsitesCount"`
	CachedEndpointsCount    int        `gorm:"column:cached_endpoints_count" json:"cachedEndpointsCount"`
	CachedIPsCount          int        `gorm:"column:cached_ips_count" json:"cachedIpsCount"`
	CachedDirectoriesCount  int        `gorm:"column:cached_directories_count" json:"cachedDirectoriesCount"`
	CachedScreenshotsCount  int        `gorm:"column:cached_screenshots_count" json:"cachedScreenshotsCount"`
	CachedVulnsTotal        int        `gorm:"column:cached_vulns_total" json:"cachedVulnsTotal"`
	CachedVulnsCritical     int        `gorm:"column:cached_vulns_critical" json:"cachedVulnsCritical"`
	CachedVulnsHigh         int        `gorm:"column:cached_vulns_high" json:"cachedVulnsHigh"`
	CachedVulnsMedium       int        `gorm:"column:cached_vulns_medium" json:"cachedVulnsMedium"`
	CachedVulnsLow          int        `gorm:"column:cached_vulns_low" json:"cachedVulnsLow"`
	StatsUpdatedAt          *time.Time `gorm:"column:stats_updated_at" json:"statsUpdatedAt"`

	// Relationships
	Target *Target `gorm:"foreignKey:TargetID" json:"target,omitempty"`
}

// TableName returns the table name for Scan
func (Scan) TableName() string {
	return "scan"
}

// ScanStatus constants
const (
	ScanStatusInitiated  = "initiated"
	ScanStatusRunning    = "running"
	ScanStatusCompleted  = "completed"
	ScanStatusFailed     = "failed"
	ScanStatusStopped    = "stopped"
	ScanStatusPending    = "pending"
)

// ScanMode constants
const (
	ScanModeFull        = "full"
	ScanModeIncremental = "incremental"
)
