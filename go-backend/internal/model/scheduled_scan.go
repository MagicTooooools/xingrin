package model

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/datatypes"
)

// ScheduledScan represents a scheduled scan task
type ScheduledScan struct {
	ID                int            `gorm:"primaryKey;autoIncrement" json:"id"`
	Name              string         `gorm:"column:name;size:200" json:"name"`
	EngineIDs         pq.Int64Array  `gorm:"column:engine_ids;type:integer[]" json:"engineIds"`
	EngineNames       datatypes.JSON `gorm:"column:engine_names;type:jsonb" json:"engineNames"`
	YamlConfiguration string         `gorm:"column:yaml_configuration;type:text" json:"yamlConfiguration"`
	OrganizationID    *int           `gorm:"column:organization_id" json:"organizationId"`
	TargetID          *int           `gorm:"column:target_id" json:"targetId"`
	CronExpression    string         `gorm:"column:cron_expression;size:100;default:'0 2 * * *'" json:"cronExpression"`
	IsEnabled         bool           `gorm:"column:is_enabled;default:true;index" json:"isEnabled"`
	RunCount          int            `gorm:"column:run_count;default:0" json:"runCount"`
	LastRunTime       *time.Time     `gorm:"column:last_run_time" json:"lastRunTime"`
	NextRunTime       *time.Time     `gorm:"column:next_run_time" json:"nextRunTime"`
	CreatedAt         time.Time      `gorm:"column:created_at;autoCreateTime;index" json:"createdAt"`
	UpdatedAt         time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`

	// Relationships
	Organization *Organization `gorm:"foreignKey:OrganizationID" json:"organization,omitempty"`
	Target       *Target       `gorm:"foreignKey:TargetID" json:"target,omitempty"`
}

// TableName returns the table name for ScheduledScan
func (ScheduledScan) TableName() string {
	return "scheduled_scan"
}
