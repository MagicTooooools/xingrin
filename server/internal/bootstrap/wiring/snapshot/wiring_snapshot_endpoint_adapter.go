package snapshotwiring

import (
	"database/sql"

	snapshotdomain "github.com/yyhuni/lunafox/server/internal/modules/snapshot/domain"
	snapshotrepo "github.com/yyhuni/lunafox/server/internal/modules/snapshot/repository"
)

type snapshotEndpointStoreAdapter struct {
	repo *snapshotrepo.EndpointSnapshotRepository
}

func newSnapshotEndpointStoreAdapter(repo *snapshotrepo.EndpointSnapshotRepository) *snapshotEndpointStoreAdapter {
	return &snapshotEndpointStoreAdapter{repo: repo}
}

func (adapter *snapshotEndpointStoreAdapter) BulkCreate(snapshots []snapshotdomain.EndpointSnapshot) (int64, error) {
	return adapter.repo.BulkCreate(snapshots)
}

func (adapter *snapshotEndpointStoreAdapter) FindByScanID(scanID int, page, pageSize int, filter string) ([]snapshotdomain.EndpointSnapshot, int64, error) {
	return adapter.repo.FindByScanID(scanID, page, pageSize, filter)
}

func (adapter *snapshotEndpointStoreAdapter) StreamByScanID(scanID int) (*sql.Rows, error) {
	return adapter.repo.StreamByScanID(scanID)
}

func (adapter *snapshotEndpointStoreAdapter) CountByScanID(scanID int) (int64, error) {
	return adapter.repo.CountByScanID(scanID)
}

func (adapter *snapshotEndpointStoreAdapter) ScanRow(rows *sql.Rows) (*snapshotdomain.EndpointSnapshot, error) {
	return adapter.repo.ScanRow(rows)
}
