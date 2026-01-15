package model

import (
	"time"
)

// Subdomain represents a subdomain asset
type Subdomain struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	TargetID  int       `gorm:"column:target_id;not null;index:idx_subdomain_target;uniqueIndex:unique_subdomain_name_target,priority:2" json:"targetId"`
	Name      string    `gorm:"column:name;size:1000;index:idx_subdomain_name;uniqueIndex:unique_subdomain_name_target,priority:1" json:"name"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime;index:idx_subdomain_created_at" json:"createdAt"`

	// Relationships
	Target *Target `gorm:"foreignKey:TargetID" json:"target,omitempty"`
}

// TableName returns the table name for Subdomain
func (Subdomain) TableName() string {
	return "subdomain"
}
