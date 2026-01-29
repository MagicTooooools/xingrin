package repository

import (
	"strings"
	"testing"
)

func TestPullTaskSQLIncludesSkipLocked(t *testing.T) {
	if !strings.Contains(pullTaskSQL, "FOR UPDATE OF st SKIP LOCKED") {
		t.Fatalf("pullTaskSQL missing FOR UPDATE OF st SKIP LOCKED")
	}
	if !strings.Contains(pullTaskSQL, "ORDER BY st.stage DESC, st.created_at ASC") {
		t.Fatalf("pullTaskSQL missing expected ordering (stage DESC to prioritize completing existing scans)")
	}
	if !strings.Contains(pullTaskSQL, "s.deleted_at IS NULL") {
		t.Fatalf("pullTaskSQL missing soft-delete check")
	}
}

func TestFailTasksSQLChecksScanStatus(t *testing.T) {
	// Should check scan deleted_at
	if !strings.Contains(failTasksSQL, "s.deleted_at IS NOT NULL") {
		t.Fatalf("failTasksSQL missing deleted_at check")
	}
	// Should check scan status
	if !strings.Contains(failTasksSQL, "s.status NOT IN ('pending', 'running')") {
		t.Fatalf("failTasksSQL missing scan status check")
	}
	// Should have appropriate error messages
	if !strings.Contains(failTasksSQL, "Scan deleted") {
		t.Fatalf("failTasksSQL missing 'Scan deleted' error message")
	}
	if !strings.Contains(failTasksSQL, "Scan already ended") {
		t.Fatalf("failTasksSQL missing 'Scan already ended' error message")
	}
	if !strings.Contains(failTasksSQL, "Agent offline") {
		t.Fatalf("failTasksSQL missing 'Agent offline' error message")
	}
}
