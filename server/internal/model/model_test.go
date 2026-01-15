package model

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
	"time"
)

// TableNamer interface for models with TableName method
type TableNamer interface {
	TableName() string
}

// TestTableNames tests that all models return correct table names
// Property 1: 数据库表名映射正确性
// *对于任意* Go 模型，其 TableName() 方法返回的表名应与 Django 模型的 db_table 一致。
// **验证: 需求 4.1**
func TestTableNames(t *testing.T) {
	tests := []struct {
		model    TableNamer
		expected string
	}{
		// Base models
		{Organization{}, "organization"},
		{Target{}, "target"},
		{Scan{}, "scan"},
		{Subdomain{}, "subdomain"},
		{Website{}, "website"},
		{WorkerNode{}, "worker_node"},
		{ScanEngine{}, "scan_engine"},
		// Asset models
		{Endpoint{}, "endpoint"},
		{Directory{}, "directory"},
		{HostPort{}, "host_port_mapping"},
		{Vulnerability{}, "vulnerability"},
		{Screenshot{}, "screenshot"},
		// Snapshot models
		{SubdomainSnapshot{}, "subdomain_snapshot"},
		{WebsiteSnapshot{}, "website_snapshot"},
		{EndpointSnapshot{}, "endpoint_snapshot"},
		{DirectorySnapshot{}, "directory_snapshot"},
		{HostPortSnapshot{}, "host_port_mapping_snapshot"},
		{VulnerabilitySnapshot{}, "vulnerability_snapshot"},
		{ScreenshotSnapshot{}, "screenshot_snapshot"},
		// Scan-related models
		{ScanLog{}, "scan_log"},
		{ScanInputTarget{}, "scan_input_target"},
		{ScheduledScan{}, "scheduled_scan"},
		{SubfinderProviderSettings{}, "subfinder_provider_settings"},
		// Engine models
		{Wordlist{}, "wordlist"},
		{NucleiTemplateRepo{}, "nuclei_template_repo"},
		// Notification models
		{Notification{}, "notification"},
		{NotificationSettings{}, "notification_settings"},
		// Config models
		{BlacklistRule{}, "blacklist_rule"},
		// Statistics models
		{AssetStatistics{}, "asset_statistics"},
		{StatisticsHistory{}, "statistics_history"},
		// Auth models
		{User{}, "auth_user"},
		{Session{}, "django_session"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.model.TableName(); got != tt.expected {
				t.Errorf("TableName() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// TestJSONFieldNames tests that JSON field names are camelCase
// Property 2: JSON 字段名转换正确性
// *对于任意* Go 模型序列化为 JSON，所有字段名应为 camelCase 格式。
// **验证: 需求 4.6**
func TestJSONFieldNames(t *testing.T) {
	// Test Target model
	target := Target{
		ID:            1,
		Name:          "test.com",
		Type:          "domain",
		CreatedAt:     time.Now(),
		LastScannedAt: nil,
	}

	jsonBytes, err := json.Marshal(target)
	if err != nil {
		t.Fatalf("Failed to marshal Target: %v", err)
	}
	jsonStr := string(jsonBytes)

	// Should contain camelCase
	camelCaseFields := []string{"id", "name", "type", "createdAt", "lastScannedAt"}
	for _, field := range camelCaseFields {
		if !strings.Contains(jsonStr, `"`+field+`"`) {
			t.Errorf("JSON should contain camelCase field %q, got: %s", field, jsonStr)
		}
	}

	// Should NOT contain snake_case
	snakeCaseFields := []string{"created_at", "last_scanned_at", "organization_id"}
	for _, field := range snakeCaseFields {
		if strings.Contains(jsonStr, `"`+field+`"`) {
			t.Errorf("JSON should NOT contain snake_case field %q, got: %s", field, jsonStr)
		}
	}
}

// TestScanJSONFieldNames tests Scan model JSON serialization
func TestScanJSONFieldNames(t *testing.T) {
	scan := Scan{
		ID:                    1,
		TargetID:              1,
		Status:                "running",
		Progress:              50,
		CurrentStage:          "subdomain_discovery",
		CachedSubdomainsCount: 100,
		CachedWebsitesCount:   50,
	}

	jsonBytes, err := json.Marshal(scan)
	if err != nil {
		t.Fatalf("Failed to marshal Scan: %v", err)
	}
	jsonStr := string(jsonBytes)

	// Should contain camelCase
	camelCaseFields := []string{
		"targetId", "status", "progress", "currentStage",
		"cachedSubdomainsCount", "cachedWebsitesCount",
	}
	for _, field := range camelCaseFields {
		if !strings.Contains(jsonStr, `"`+field+`"`) {
			t.Errorf("JSON should contain camelCase field %q, got: %s", field, jsonStr)
		}
	}

	// Should NOT contain snake_case
	snakeCaseFields := []string{
		"target_id", "current_stage", "cached_subdomains_count",
	}
	for _, field := range snakeCaseFields {
		if strings.Contains(jsonStr, `"`+field+`"`) {
			t.Errorf("JSON should NOT contain snake_case field %q, got: %s", field, jsonStr)
		}
	}
}

// TestGORMColumnTags tests that GORM column tags use snake_case
// Property 3: 数据库字段映射正确性
// *对于任意* Go 模型字段，其 gorm column tag 应与数据库实际列名（snake_case）一致。
// **验证: 需求 4.2**
func TestGORMColumnTags(t *testing.T) {
	// Test Target model
	targetType := reflect.TypeOf(Target{})
	expectedColumns := map[string]string{
		"ID":            "", // primaryKey, no explicit column
		"Name":          "name",
		"Type":          "type",
		"CreatedAt":     "created_at",
		"LastScannedAt": "last_scanned_at",
		"DeletedAt":     "deleted_at",
	}

	for fieldName, expectedColumn := range expectedColumns {
		field, found := targetType.FieldByName(fieldName)
		if !found {
			t.Errorf("Field %s not found in Target", fieldName)
			continue
		}

		gormTag := field.Tag.Get("gorm")
		if expectedColumn != "" && !strings.Contains(gormTag, "column:"+expectedColumn) {
			t.Errorf("Field %s: expected gorm column:%s, got tag: %s", fieldName, expectedColumn, gormTag)
		}
	}
}

// TestScanGORMColumnTags tests Scan model GORM tags
func TestScanGORMColumnTags(t *testing.T) {
	scanType := reflect.TypeOf(Scan{})
	expectedColumns := map[string]string{
		"TargetID":              "target_id",
		"EngineIDs":             "engine_ids",
		"EngineNames":           "engine_names",
		"YamlConfiguration":     "yaml_configuration",
		"ScanMode":              "scan_mode",
		"Status":                "status",
		"ResultsDir":            "results_dir",
		"ContainerIDs":          "container_ids",
		"WorkerID":              "worker_id",
		"ErrorMessage":          "error_message",
		"Progress":              "progress",
		"CurrentStage":          "current_stage",
		"StageProgress":         "stage_progress",
		"CreatedAt":             "created_at",
		"StoppedAt":             "stopped_at",
		"DeletedAt":             "deleted_at",
		"CachedSubdomainsCount": "cached_subdomains_count",
		"CachedWebsitesCount":   "cached_websites_count",
		"CachedEndpointsCount":  "cached_endpoints_count",
		"CachedIPsCount":        "cached_ips_count",
		"CachedVulnsTotal":      "cached_vulns_total",
	}

	for fieldName, expectedColumn := range expectedColumns {
		field, found := scanType.FieldByName(fieldName)
		if !found {
			t.Errorf("Field %s not found in Scan", fieldName)
			continue
		}

		gormTag := field.Tag.Get("gorm")
		if !strings.Contains(gormTag, "column:"+expectedColumn) {
			t.Errorf("Field %s: expected gorm column:%s, got tag: %s", fieldName, expectedColumn, gormTag)
		}
	}
}

// TestWorkerNodePasswordHidden tests that password is hidden from JSON
func TestWorkerNodePasswordHidden(t *testing.T) {
	worker := WorkerNode{
		ID:        1,
		Name:      "worker-1",
		IPAddress: "192.168.1.1",
		Password:  "secret123",
		Status:    "connected",
	}

	jsonBytes, err := json.Marshal(worker)
	if err != nil {
		t.Fatalf("Failed to marshal WorkerNode: %v", err)
	}
	jsonStr := string(jsonBytes)

	// Password should NOT appear in JSON
	if strings.Contains(jsonStr, "secret123") {
		t.Errorf("Password should be hidden from JSON, got: %s", jsonStr)
	}
	if strings.Contains(jsonStr, `"password"`) {
		t.Errorf("Password field should not appear in JSON, got: %s", jsonStr)
	}
}
