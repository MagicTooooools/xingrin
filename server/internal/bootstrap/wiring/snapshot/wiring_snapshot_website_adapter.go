package snapshotwiring

import (
	"database/sql"

	snapshotdomain "github.com/yyhuni/lunafox/server/internal/modules/snapshot/domain"
	snapshotrepo "github.com/yyhuni/lunafox/server/internal/modules/snapshot/repository"
)

type snapshotWebsiteStoreAdapter struct {
	repo *snapshotrepo.WebsiteSnapshotRepository
}

func newSnapshotWebsiteStoreAdapter(repo *snapshotrepo.WebsiteSnapshotRepository) *snapshotWebsiteStoreAdapter {
	return &snapshotWebsiteStoreAdapter{repo: repo}
}

func (adapter *snapshotWebsiteStoreAdapter) BulkCreate(snapshots []snapshotdomain.WebsiteSnapshot) (int64, error) {
	return adapter.repo.BulkCreate(snapshots)
}

func (adapter *snapshotWebsiteStoreAdapter) FindByScanID(scanID int, page, pageSize int, filter string) ([]snapshotdomain.WebsiteSnapshot, int64, error) {
	return adapter.repo.FindByScanID(scanID, page, pageSize, filter)
}

func (adapter *snapshotWebsiteStoreAdapter) StreamByScanID(scanID int) (*sql.Rows, error) {
	return adapter.repo.StreamByScanID(scanID)
}

func (adapter *snapshotWebsiteStoreAdapter) CountByScanID(scanID int) (int64, error) {
	return adapter.repo.CountByScanID(scanID)
}

func (adapter *snapshotWebsiteStoreAdapter) ScanRow(rows *sql.Rows) (*snapshotdomain.WebsiteSnapshot, error) {
	return adapter.repo.ScanRow(rows)
}
