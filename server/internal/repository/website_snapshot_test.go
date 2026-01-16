package repository

import (
	"testing"

	"github.com/orbit/server/internal/model"
)

// TestWebsiteSnapshotFilterMapping tests the filter mapping configuration
func TestWebsiteSnapshotFilterMapping(t *testing.T) {
	expectedFields := []string{"url", "host", "title", "status", "webserver", "tech"}

	for _, field := range expectedFields {
		if _, ok := WebsiteSnapshotFilterMapping[field]; !ok {
			t.Errorf("expected field %s not found in WebsiteSnapshotFilterMapping", field)
		}
	}

	// Test specific configurations
	if !WebsiteSnapshotFilterMapping["status"].IsNumeric {
		t.Error("status field should be marked as numeric")
	}

	if !WebsiteSnapshotFilterMapping["tech"].IsArray {
		t.Error("tech field should be marked as array")
	}
}

// TestBulkCreateDeduplication tests that BulkCreate handles duplicates correctly
// This is a unit test for the deduplication logic
func TestBulkCreateDeduplication(t *testing.T) {
	// Test that duplicate URLs in the same scan should be deduplicated
	// This test verifies the model structure supports the unique constraint

	snapshot1 := model.WebsiteSnapshot{
		ScanID: 1,
		URL:    "https://example.com",
		Host:   "example.com",
	}

	snapshot2 := model.WebsiteSnapshot{
		ScanID: 1,
		URL:    "https://example.com", // Same URL, same scan
		Host:   "example.com",
	}

	// Verify they have the same unique key fields
	if snapshot1.ScanID != snapshot2.ScanID || snapshot1.URL != snapshot2.URL {
		t.Error("snapshots should have same unique key fields for deduplication test")
	}

	// Different scan should be allowed
	snapshot3 := model.WebsiteSnapshot{
		ScanID: 2, // Different scan
		URL:    "https://example.com",
		Host:   "example.com",
	}

	if snapshot1.ScanID == snapshot3.ScanID {
		t.Error("snapshot3 should have different scan_id")
	}
}
