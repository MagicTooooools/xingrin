package model

import (
	"time"
)

// Target represents a scan target (domain, IP, etc.)
type Target struct {
	ID             int        `gorm:"primaryKey;autoIncrement" json:"id"`
	Name           string     `gorm:"column:name;size:300" json:"name"`
	Type           string     `gorm:"column:type;size:20;default:'domain'" json:"type"`
	OrganizationID *int       `gorm:"column:organization_id" json:"organizationId"`
	CreatedAt      time.Time  `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	LastScannedAt  *time.Time `gorm:"column:last_scanned_at" json:"lastScannedAt"`
	DeletedAt      *time.Time `gorm:"column:deleted_at;index" json:"-"`

	// Relationships
	Organization *Organization `gorm:"foreignKey:OrganizationID" json:"organization,omitempty"`
}

// TableName returns the table name for Target
func (Target) TableName() string {
	return "target"
}
