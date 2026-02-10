package scanlogwiring

import (
	scanapp "github.com/yyhuni/lunafox/server/internal/modules/scan/application"
	scanrepo "github.com/yyhuni/lunafox/server/internal/modules/scan/repository"
	"github.com/yyhuni/lunafox/server/internal/pkg/timeutil"
)

type scanLogStoreAdapter struct {
	repo *scanrepo.ScanLogRepository
}

func (adapter *scanLogStoreAdapter) FindByScanIDWithCursor(scanID int, afterID int64, limit int) ([]scanapp.ScanLogEntry, error) {
	logs, err := adapter.repo.FindByScanIDWithCursor(scanID, afterID, limit)
	if err != nil {
		return nil, err
	}
	results := make([]scanapp.ScanLogEntry, 0, len(logs))
	for _, item := range logs {
		results = append(results, scanapp.ScanLogEntry{
			ID:        item.ID,
			ScanID:    item.ScanID,
			Level:     item.Level,
			Content:   item.Content,
			CreatedAt: timeutil.ToUTC(item.CreatedAt),
		})
	}
	return results, nil
}

func (adapter *scanLogStoreAdapter) BulkCreate(logs []scanapp.ScanLogEntry) error {
	items := make([]scanrepo.ScanLogRecord, 0, len(logs))
	for _, item := range logs {
		items = append(items, scanrepo.ScanLogRecord{
			ID:        item.ID,
			ScanID:    item.ScanID,
			Level:     item.Level,
			Content:   item.Content,
			CreatedAt: timeutil.ToUTC(item.CreatedAt),
		})
	}
	return adapter.repo.BulkCreate(items)
}
