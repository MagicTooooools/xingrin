package application

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

var (
	ErrScanLogScanNotFound = errors.New("scan not found")
)

type ScanLogEntry struct {
	ID        int64
	ScanID    int
	Level     string
	Content   string
	CreatedAt time.Time
}

type ScanLogScanRef struct{ ID int }

type ScanLogStore interface {
	FindByScanIDWithCursor(scanID int, afterID int64, limit int) ([]ScanLogEntry, error)
	BulkCreate(logs []ScanLogEntry) error
}

type ScanLookup interface {
	FindByID(id int) (*ScanLogScanRef, error)
}

type ScanLogCreateItem struct {
	Level   string
	Content string
}

type ScanLogService struct {
	logStore   ScanLogStore
	scanLookup ScanLookup
}

func NewScanLogService(logStore ScanLogStore, scanLookup ScanLookup) *ScanLogService {
	return &ScanLogService{logStore: logStore, scanLookup: scanLookup}
}

func (service *ScanLogService) ListByScanID(ctx context.Context, scanID int, afterID int64, limit int) ([]ScanLogEntry, bool, error) {
	_ = ctx
	_, err := service.scanLookup.FindByID(scanID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, ErrScanLogScanNotFound
		}
		return nil, false, err
	}

	if afterID < 0 {
		afterID = 0
	}
	if limit <= 0 {
		limit = 200
	}
	if limit > 1000 {
		limit = 1000
	}

	logs, err := service.logStore.FindByScanIDWithCursor(scanID, afterID, limit+1)
	if err != nil {
		return nil, false, err
	}
	hasMore := len(logs) > limit
	if hasMore {
		logs = logs[:limit]
	}
	return logs, hasMore, nil
}

func (service *ScanLogService) BulkCreate(ctx context.Context, scanID int, items []ScanLogCreateItem) (int, error) {
	_ = ctx
	_, err := service.scanLookup.FindByID(scanID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, ErrScanLogScanNotFound
		}
		return 0, err
	}
	if len(items) == 0 {
		return 0, nil
	}
	logs := make([]ScanLogEntry, len(items))
	for index, item := range items {
		logs[index] = ScanLogEntry{ScanID: scanID, Level: item.Level, Content: item.Content}
	}
	if err := service.logStore.BulkCreate(logs); err != nil {
		return 0, err
	}
	return len(logs), nil
}
