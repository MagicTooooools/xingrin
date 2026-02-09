package snapshotwiring

import (
	"database/sql"

	snapshotdomain "github.com/yyhuni/lunafox/server/internal/modules/snapshot/domain"
	snapshotrepo "github.com/yyhuni/lunafox/server/internal/modules/snapshot/repository"
	snapshotmodel "github.com/yyhuni/lunafox/server/internal/modules/snapshot/repository/persistence"
)

type snapshotEndpointStoreAdapter struct {
	repo *snapshotrepo.EndpointSnapshotRepository
}

func newSnapshotEndpointStoreAdapter(repo *snapshotrepo.EndpointSnapshotRepository) *snapshotEndpointStoreAdapter {
	return &snapshotEndpointStoreAdapter{repo: repo}
}

func (adapter *snapshotEndpointStoreAdapter) BulkCreate(snapshots []snapshotdomain.EndpointSnapshot) (int64, error) {
	items := make([]snapshotmodel.EndpointSnapshot, 0, len(snapshots))
	for index := range snapshots {
		items = append(items, *snapshotEndpointDomainToModel(&snapshots[index]))
	}
	return adapter.repo.BulkCreate(items)
}

func (adapter *snapshotEndpointStoreAdapter) FindByScanID(scanID int, page, pageSize int, filter string) ([]snapshotdomain.EndpointSnapshot, int64, error) {
	items, total, err := adapter.repo.FindByScanID(scanID, page, pageSize, filter)
	if err != nil {
		return nil, 0, err
	}
	results := make([]snapshotdomain.EndpointSnapshot, 0, len(items))
	for index := range items {
		results = append(results, *snapshotEndpointModelToDomain(&items[index]))
	}
	return results, total, nil
}

func (adapter *snapshotEndpointStoreAdapter) StreamByScanID(scanID int) (*sql.Rows, error) {
	return adapter.repo.StreamByScanID(scanID)
}

func (adapter *snapshotEndpointStoreAdapter) CountByScanID(scanID int) (int64, error) {
	return adapter.repo.CountByScanID(scanID)
}

func (adapter *snapshotEndpointStoreAdapter) ScanRow(rows *sql.Rows) (*snapshotdomain.EndpointSnapshot, error) {
	item, err := adapter.repo.ScanRow(rows)
	if err != nil {
		return nil, err
	}
	return snapshotEndpointModelToDomain(item), nil
}
