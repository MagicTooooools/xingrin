package repository

import (
	"strings"
	"testing"
)

func TestPullTaskSQLIncludesSkipLocked(t *testing.T) {
	if !strings.Contains(pullTaskSQL, "FOR UPDATE SKIP LOCKED") {
		t.Fatalf("pullTaskSQL missing FOR UPDATE SKIP LOCKED")
	}
	if !strings.Contains(pullTaskSQL, "ORDER BY stage DESC, created_at ASC") {
		t.Fatalf("pullTaskSQL missing expected ordering")
	}
}
