package dto

import "time"

// ScanListQuery represents scan list query parameters
type ScanListQuery struct {
	PaginationQuery
	TargetID int    `form:"target" binding:"omitempty"`
	Status   string `form:"status" binding:"omitempty"`
	Search   string `form:"search" binding:"omitempty"`
}

// ScanResponse represents scan response
type ScanResponse struct {
	ID           int              `json:"id"`
	TargetID     int              `json:"targetId"`
	EngineIDs    []int64          `json:"engineIds"`
	EngineNames  []string         `json:"engineNames"`
	ScanMode     string           `json:"scanMode"`
	Status       string           `json:"status"`
	Progress     int              `json:"progress"`
	CurrentStage string           `json:"currentStage"`
	ErrorMessage string           `json:"errorMessage,omitempty"`
	CreatedAt    time.Time        `json:"createdAt"`
	StoppedAt    *time.Time       `json:"stoppedAt,omitempty"`
	Target       *TargetBrief     `json:"target,omitempty"`
	CachedStats  *ScanCachedStats `json:"cachedStats,omitempty"`
}

// TargetBrief represents brief target info for scan response
type TargetBrief struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

// ScanCachedStats represents cached statistics for scan
type ScanCachedStats struct {
	SubdomainsCount  int `json:"subdomainsCount"`
	WebsitesCount    int `json:"websitesCount"`
	EndpointsCount   int `json:"endpointsCount"`
	IPsCount         int `json:"ipsCount"`
	DirectoriesCount int `json:"directoriesCount"`
	ScreenshotsCount int `json:"screenshotsCount"`
	VulnsTotal       int `json:"vulnsTotal"`
	VulnsCritical    int `json:"vulnsCritical"`
	VulnsHigh        int `json:"vulnsHigh"`
	VulnsMedium      int `json:"vulnsMedium"`
	VulnsLow         int `json:"vulnsLow"`
}

// ScanDetailResponse represents detailed scan response
type ScanDetailResponse struct {
	ScanResponse
	YamlConfiguration string                 `json:"yamlConfiguration,omitempty"`
	ResultsDir        string                 `json:"resultsDir,omitempty"`
	WorkerID          *int                   `json:"workerId,omitempty"`
	StageProgress     map[string]interface{} `json:"stageProgress,omitempty"`
}

// InitiateScanRequest represents initiate scan request (deprecated, use CreateScanRequest)
type InitiateScanRequest struct {
	OrganizationID *int     `json:"organizationId" binding:"omitempty"`
	TargetID       *int     `json:"targetId" binding:"omitempty"`
	EngineIDs      []int    `json:"engineIds" binding:"required,min=1"`
	EngineNames    []string `json:"engineNames" binding:"required,min=1"`
	Configuration  string   `json:"configuration" binding:"required"`
}

// QuickScanRequest represents quick scan request (deprecated, use CreateScanRequest)
type QuickScanRequest struct {
	Targets       []QuickScanTarget `json:"targets" binding:"required,min=1"`
	EngineIDs     []int             `json:"engineIds"`
	EngineNames   []string          `json:"engineNames"`
	Configuration string            `json:"configuration" binding:"required"`
}

// CreateScanRequest represents unified scan creation request
// POST /api/scans
type CreateScanRequest struct {
	// Mode: "normal" (default) or "quick"
	Mode string `json:"mode" binding:"omitempty,oneof=normal quick"`

	// For mode=normal: target ID (required)
	TargetID int `json:"targetId" binding:"omitempty"`

	// For mode=quick: raw targets (required)
	Targets []string `json:"targets" binding:"omitempty"`

	// Common fields
	EngineIDs     []int  `json:"engineIds" binding:"omitempty"`
	EngineNames   []string `json:"engineNames" binding:"omitempty"`
	Configuration string `json:"configuration" binding:"omitempty"`
}

// QuickScanTarget represents a target in quick scan
type QuickScanTarget struct {
	Name string `json:"name" binding:"required"`
}

// QuickScanResponse represents quick scan response
type QuickScanResponse struct {
	Count       int            `json:"count"`
	TargetStats map[string]int `json:"targetStats"`
	AssetStats  map[string]int `json:"assetStats"`
	Errors      []string       `json:"errors,omitempty"`
	Scans       []ScanResponse `json:"scans"`
}

// StopScanResponse represents stop scan response
type StopScanResponse struct {
	RevokedTaskCount int `json:"revokedTaskCount"`
}

// ScanStatisticsResponse represents scan statistics response
type ScanStatisticsResponse struct {
	Total           int64 `json:"total"`
	Running         int64 `json:"running"`
	Completed       int64 `json:"completed"`
	Failed          int64 `json:"failed"`
	TotalVulns      int64 `json:"totalVulns"`
	TotalSubdomains int64 `json:"totalSubdomains"`
	TotalEndpoints  int64 `json:"totalEndpoints"`
	TotalWebsites   int64 `json:"totalWebsites"`
	TotalAssets     int64 `json:"totalAssets"`
}

// ScanLogListQuery represents scan log list query parameters (cursor pagination)
type ScanLogListQuery struct {
	AfterID int64 `form:"afterId" binding:"omitempty,min=0"`
	Limit   int   `form:"limit" binding:"omitempty,min=1,max=1000"`
}

// ScanLogListResponse represents scan log list response
type ScanLogListResponse struct {
	Results []ScanLogResponse `json:"results"`
	HasMore bool              `json:"hasMore"`
}

// ScanLogResponse represents scan log response
type ScanLogResponse struct {
	ID        int64     `json:"id"`
	ScanID    int       `json:"scanId"`
	Level     string    `json:"level"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
}

// ScanLogItem represents a single log item for bulk create
type ScanLogItem struct {
	Level   string `json:"level" binding:"required,oneof=info warning error"`
	Content string `json:"content" binding:"required"`
}

// BulkCreateScanLogsRequest represents bulk create scan logs request
type BulkCreateScanLogsRequest struct {
	Logs []ScanLogItem `json:"logs" binding:"required,min=1,max=1000,dive"`
}

// BulkCreateScanLogsResponse represents bulk create scan logs response
type BulkCreateScanLogsResponse struct {
	CreatedCount int `json:"createdCount"`
}
