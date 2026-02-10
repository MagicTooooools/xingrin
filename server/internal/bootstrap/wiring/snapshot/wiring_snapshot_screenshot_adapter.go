package snapshotwiring

import (
	snapshotdomain "github.com/yyhuni/lunafox/server/internal/modules/snapshot/domain"
	snapshotrepo "github.com/yyhuni/lunafox/server/internal/modules/snapshot/repository"
)

type snapshotScreenshotStoreAdapter struct {
	repo *snapshotrepo.ScreenshotSnapshotRepository
}

func newSnapshotScreenshotStoreAdapter(repo *snapshotrepo.ScreenshotSnapshotRepository) *snapshotScreenshotStoreAdapter {
	return &snapshotScreenshotStoreAdapter{repo: repo}
}

func (adapter *snapshotScreenshotStoreAdapter) BulkUpsert(snapshots []snapshotdomain.ScreenshotSnapshot) (int64, error) {
	return adapter.repo.BulkUpsert(snapshots)
}

func (adapter *snapshotScreenshotStoreAdapter) FindByScanID(scanID int, page, pageSize int, filter string) ([]snapshotdomain.ScreenshotSnapshot, int64, error) {
	return adapter.repo.FindByScanID(scanID, page, pageSize, filter)
}

func (adapter *snapshotScreenshotStoreAdapter) FindByIDAndScanID(id int, scanID int) (*snapshotdomain.ScreenshotSnapshot, error) {
	return adapter.repo.FindByIDAndScanID(id, scanID)
}

func (adapter *snapshotScreenshotStoreAdapter) CountByScanID(scanID int) (int64, error) {
	return adapter.repo.CountByScanID(scanID)
}
