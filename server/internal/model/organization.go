package model

import (
	"time"
)

// Organization represents an organization
type Organization struct {
	ID          int        `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string     `gorm:"column:name;size:300;index:idx_org_name" json:"name"`
	Description string     `gorm:"column:description;size:1000" json:"description"`
	CreatedAt   time.Time  `gorm:"column:created_at;autoCreateTime;index:idx_org_created_at" json:"createdAt"`
	DeletedAt   *time.Time `gorm:"column:deleted_at;index:idx_org_deleted_at" json:"-"`

	// Many-to-many relationship with Target
	Targets []Target `gorm:"many2many:organization_target;" json:"targets,omitempty"`
}

// TableName returns the table name for Organization
func (Organization) TableName() string {
	return "organization"
}
