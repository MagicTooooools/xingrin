package snapshotwiring

import (
	"database/sql"

	snapshotdomain "github.com/yyhuni/lunafox/server/internal/modules/snapshot/domain"
	snapshotrepo "github.com/yyhuni/lunafox/server/internal/modules/snapshot/repository"
)

type snapshotHostPortStoreAdapter struct {
	repo *snapshotrepo.HostPortSnapshotRepository
}

func newSnapshotHostPortStoreAdapter(repo *snapshotrepo.HostPortSnapshotRepository) *snapshotHostPortStoreAdapter {
	return &snapshotHostPortStoreAdapter{repo: repo}
}

func (adapter *snapshotHostPortStoreAdapter) BulkCreate(snapshots []snapshotdomain.HostPortSnapshot) (int64, error) {
	return adapter.repo.BulkCreate(snapshots)
}

func (adapter *snapshotHostPortStoreAdapter) FindByScanID(scanID int, page, pageSize int, filter string) ([]snapshotdomain.HostPortSnapshot, int64, error) {
	return adapter.repo.FindByScanID(scanID, page, pageSize, filter)
}

func (adapter *snapshotHostPortStoreAdapter) StreamByScanID(scanID int) (*sql.Rows, error) {
	return adapter.repo.StreamByScanID(scanID)
}

func (adapter *snapshotHostPortStoreAdapter) CountByScanID(scanID int) (int64, error) {
	return adapter.repo.CountByScanID(scanID)
}

func (adapter *snapshotHostPortStoreAdapter) ScanRow(rows *sql.Rows) (*snapshotdomain.HostPortSnapshot, error) {
	return adapter.repo.ScanRow(rows)
}
