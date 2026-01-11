package model

import (
	"time"
)

// StatisticsHistory represents daily statistics history
type StatisticsHistory struct {
	ID              int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Date            time.Time `gorm:"column:date;type:date;uniqueIndex" json:"date"`
	TotalTargets    int       `gorm:"column:total_targets;default:0" json:"totalTargets"`
	TotalSubdomains int       `gorm:"column:total_subdomains;default:0" json:"totalSubdomains"`
	TotalIPs        int       `gorm:"column:total_ips;default:0" json:"totalIps"`
	TotalEndpoints  int       `gorm:"column:total_endpoints;default:0" json:"totalEndpoints"`
	TotalWebsites   int       `gorm:"column:total_websites;default:0" json:"totalWebsites"`
	TotalVulns      int       `gorm:"column:total_vulns;default:0" json:"totalVulns"`
	TotalAssets     int       `gorm:"column:total_assets;default:0" json:"totalAssets"`
	CreatedAt       time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt       time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

// TableName returns the table name for StatisticsHistory
func (StatisticsHistory) TableName() string {
	return "statistics_history"
}
