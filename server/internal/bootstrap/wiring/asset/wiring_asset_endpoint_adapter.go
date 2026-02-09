package assetwiring

import (
	"database/sql"

	assetdomain "github.com/yyhuni/lunafox/server/internal/modules/asset/domain"
	assetrepo "github.com/yyhuni/lunafox/server/internal/modules/asset/repository"
)

type assetEndpointStoreAdapter struct {
	repo *assetrepo.EndpointRepository
}

func newAssetEndpointStoreAdapter(repo *assetrepo.EndpointRepository) *assetEndpointStoreAdapter {
	return &assetEndpointStoreAdapter{repo: repo}
}

func (adapter *assetEndpointStoreAdapter) FindByTargetID(targetID int, page, pageSize int, filter string) ([]assetdomain.Endpoint, int64, error) {
	items, total, err := adapter.repo.FindByTargetID(targetID, page, pageSize, filter)
	if err != nil {
		return nil, 0, err
	}
	results := make([]assetdomain.Endpoint, 0, len(items))
	for index := range items {
		results = append(results, *assetModelEndpointToDomain(&items[index]))
	}
	return results, total, nil
}

func (adapter *assetEndpointStoreAdapter) FindByID(id int) (*assetdomain.Endpoint, error) {
	item, err := adapter.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return assetModelEndpointToDomain(item), nil
}

func (adapter *assetEndpointStoreAdapter) BulkCreate(items []assetdomain.Endpoint) (int, error) {
	return adapter.repo.BulkCreate(assetDomainEndpointListToModel(items))
}

func (adapter *assetEndpointStoreAdapter) Delete(id int) error {
	return adapter.repo.Delete(id)
}

func (adapter *assetEndpointStoreAdapter) BulkDelete(ids []int) (int64, error) {
	return adapter.repo.BulkDelete(ids)
}

func (adapter *assetEndpointStoreAdapter) BulkUpsert(items []assetdomain.Endpoint) (int64, error) {
	return adapter.repo.BulkUpsert(assetDomainEndpointListToModel(items))
}

func (adapter *assetEndpointStoreAdapter) StreamByTargetID(targetID int) (*sql.Rows, error) {
	return adapter.repo.StreamByTargetID(targetID)
}

func (adapter *assetEndpointStoreAdapter) CountByTargetID(targetID int) (int64, error) {
	return adapter.repo.CountByTargetID(targetID)
}

func (adapter *assetEndpointStoreAdapter) ScanRow(rows *sql.Rows) (*assetdomain.Endpoint, error) {
	item, err := adapter.repo.ScanRow(rows)
	if err != nil {
		return nil, err
	}
	return assetModelEndpointToDomain(item), nil
}
