package model

import (
	"time"
)

// Target represents a scan target
type Target struct {
	ID            int        `gorm:"primaryKey;autoIncrement" json:"id"`
	Name          string     `gorm:"column:name;size:300;index:idx_target_name" json:"name"`
	Type          string     `gorm:"column:type;size:20;default:'domain';index:idx_target_type" json:"type"`
	CreatedAt     time.Time  `gorm:"column:created_at;autoCreateTime;index:idx_target_created_at" json:"createdAt"`
	LastScannedAt *time.Time `gorm:"column:last_scanned_at" json:"lastScannedAt"`
	DeletedAt     *time.Time `gorm:"column:deleted_at;index:idx_target_deleted_at" json:"-"`

	// Relationships
	Organizations []Organization `gorm:"many2many:organization_target;" json:"organizations,omitempty"`
}

// TableName returns the table name for Target
func (Target) TableName() string {
	return "target"
}

// TargetType constants
const (
	TargetTypeDomain = "domain"
	TargetTypeIP     = "ip"
	TargetTypeCIDR   = "cidr"
)
