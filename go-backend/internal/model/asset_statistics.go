package model

import (
	"time"
)

// AssetStatistics represents asset statistics (singleton)
type AssetStatistics struct {
	ID int `gorm:"primaryKey" json:"id"`

	// Current statistics
	TotalTargets    int `gorm:"column:total_targets;default:0" json:"totalTargets"`
	TotalSubdomains int `gorm:"column:total_subdomains;default:0" json:"totalSubdomains"`
	TotalIPs        int `gorm:"column:total_ips;default:0" json:"totalIps"`
	TotalEndpoints  int `gorm:"column:total_endpoints;default:0" json:"totalEndpoints"`
	TotalWebsites   int `gorm:"column:total_websites;default:0" json:"totalWebsites"`
	TotalVulns      int `gorm:"column:total_vulns;default:0" json:"totalVulns"`
	TotalAssets     int `gorm:"column:total_assets;default:0" json:"totalAssets"`

	// Previous statistics (for trend calculation)
	PrevTargets    int `gorm:"column:prev_targets;default:0" json:"prevTargets"`
	PrevSubdomains int `gorm:"column:prev_subdomains;default:0" json:"prevSubdomains"`
	PrevIPs        int `gorm:"column:prev_ips;default:0" json:"prevIps"`
	PrevEndpoints  int `gorm:"column:prev_endpoints;default:0" json:"prevEndpoints"`
	PrevWebsites   int `gorm:"column:prev_websites;default:0" json:"prevWebsites"`
	PrevVulns      int `gorm:"column:prev_vulns;default:0" json:"prevVulns"`
	PrevAssets     int `gorm:"column:prev_assets;default:0" json:"prevAssets"`

	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

// TableName returns the table name for AssetStatistics
func (AssetStatistics) TableName() string {
	return "asset_statistics"
}
