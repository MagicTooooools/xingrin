package snapshotwiring

import (
	"database/sql"

	snapshotdomain "github.com/yyhuni/lunafox/server/internal/modules/snapshot/domain"
	snapshotrepo "github.com/yyhuni/lunafox/server/internal/modules/snapshot/repository"
)

type snapshotDirectoryStoreAdapter struct {
	repo *snapshotrepo.DirectorySnapshotRepository
}

func newSnapshotDirectoryStoreAdapter(repo *snapshotrepo.DirectorySnapshotRepository) *snapshotDirectoryStoreAdapter {
	return &snapshotDirectoryStoreAdapter{repo: repo}
}

func (adapter *snapshotDirectoryStoreAdapter) BulkCreate(snapshots []snapshotdomain.DirectorySnapshot) (int64, error) {
	return adapter.repo.BulkCreate(snapshots)
}

func (adapter *snapshotDirectoryStoreAdapter) FindByScanID(scanID int, page, pageSize int, filter string) ([]snapshotdomain.DirectorySnapshot, int64, error) {
	return adapter.repo.FindByScanID(scanID, page, pageSize, filter)
}

func (adapter *snapshotDirectoryStoreAdapter) StreamByScanID(scanID int) (*sql.Rows, error) {
	return adapter.repo.StreamByScanID(scanID)
}

func (adapter *snapshotDirectoryStoreAdapter) CountByScanID(scanID int) (int64, error) {
	return adapter.repo.CountByScanID(scanID)
}

func (adapter *snapshotDirectoryStoreAdapter) ScanRow(rows *sql.Rows) (*snapshotdomain.DirectorySnapshot, error) {
	return adapter.repo.ScanRow(rows)
}
