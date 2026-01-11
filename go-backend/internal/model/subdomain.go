package model

import (
	"time"
)

// Subdomain represents a discovered subdomain
type Subdomain struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	TargetID  int       `gorm:"column:target_id;not null" json:"targetId"`
	Name      string    `gorm:"column:name;size:1000" json:"name"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`

	// Relationships
	Target *Target `gorm:"foreignKey:TargetID" json:"target,omitempty"`
}

// TableName returns the table name for Subdomain
func (Subdomain) TableName() string {
	return "subdomain"
}
