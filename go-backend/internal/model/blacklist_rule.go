package model

import (
	"time"
)

// BlacklistRule represents a blacklist rule
type BlacklistRule struct {
	ID          int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Pattern     string    `gorm:"column:pattern;size:255" json:"pattern"`
	RuleType    string    `gorm:"column:rule_type;size:20" json:"ruleType"`
	Scope       string    `gorm:"column:scope;size:20;index:idx_blacklist_scope" json:"scope"`
	TargetID    *int      `gorm:"column:target_id;index:idx_blacklist_target" json:"targetId"`
	Description string    `gorm:"column:description;size:500" json:"description"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime;index:idx_blacklist_created_at" json:"createdAt"`

	// Relationships
	Target *Target `gorm:"foreignKey:TargetID" json:"target,omitempty"`
}

// TableName returns the table name for BlacklistRule
func (BlacklistRule) TableName() string {
	return "blacklist_rule"
}

// RuleType constants
const (
	RuleTypeDomain  = "domain"
	RuleTypeIP      = "ip"
	RuleTypeCIDR    = "cidr"
	RuleTypeKeyword = "keyword"
)

// Scope constants
const (
	ScopeGlobal = "global"
	ScopeTarget = "target"
)
