package application

import (
	"context"
	"errors"
	"testing"

	"gorm.io/gorm"
)

type scanLogStoreStub struct {
	rows          []ScanLogEntry
	err           error
	lastScanID    int
	lastAfterID   int64
	lastLimit     int
	bulkCreateErr error
}

func (stub *scanLogStoreStub) FindByScanIDWithCursor(scanID int, afterID int64, limit int) ([]ScanLogEntry, error) {
	stub.lastScanID = scanID
	stub.lastAfterID = afterID
	stub.lastLimit = limit
	if stub.err != nil {
		return nil, stub.err
	}
	items := make([]ScanLogEntry, len(stub.rows))
	copy(items, stub.rows)
	return items, nil
}

func (stub *scanLogStoreStub) BulkCreate(logs []ScanLogEntry) error {
	_ = logs
	return stub.bulkCreateErr
}

type scanLookupStub struct {
	err error
}

func (stub *scanLookupStub) FindByID(id int) (*ScanLogScanRef, error) {
	if stub.err != nil {
		return nil, stub.err
	}
	return &ScanLogScanRef{ID: id}, nil
}

func TestScanLogServiceListByAfterID(t *testing.T) {
	store := &scanLogStoreStub{rows: []ScanLogEntry{{ID: 10}, {ID: 11}, {ID: 12}}}
	service := NewScanLogService(store, &scanLookupStub{})

	items, hasMore, err := service.ListByScanID(context.Background(), 7, 0, 2)
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	if !hasMore || len(items) != 2 {
		t.Fatalf("unexpected page: hasMore=%v len=%d", hasMore, len(items))
	}
	if store.lastAfterID != 0 || store.lastLimit != 3 {
		t.Fatalf("unexpected store args afterID=%d limit=%d", store.lastAfterID, store.lastLimit)
	}

	_, _, err = service.ListByScanID(context.Background(), 7, 9, 2)
	if err != nil {
		t.Fatalf("list with afterID failed: %v", err)
	}
	if store.lastAfterID != 9 {
		t.Fatalf("expected afterID 9, got %d", store.lastAfterID)
	}
}

func TestScanLogServiceNegativeAfterIDClamped(t *testing.T) {
	store := &scanLogStoreStub{rows: []ScanLogEntry{{ID: 1}}}
	service := NewScanLogService(store, &scanLookupStub{})

	_, _, err := service.ListByScanID(context.Background(), 7, -10, 20)
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	if store.lastAfterID != 0 {
		t.Fatalf("expected afterID clamped to 0, got %d", store.lastAfterID)
	}
}

func TestScanLogServiceScanNotFound(t *testing.T) {
	service := NewScanLogService(&scanLogStoreStub{}, &scanLookupStub{err: gorm.ErrRecordNotFound})

	_, _, err := service.ListByScanID(context.Background(), 7, 0, 20)
	if !errors.Is(err, ErrScanLogScanNotFound) {
		t.Fatalf("expected ErrScanLogScanNotFound, got %v", err)
	}
}
