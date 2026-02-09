package assetwiring

import (
	"database/sql"

	assetdomain "github.com/yyhuni/lunafox/server/internal/modules/asset/domain"
	assetrepo "github.com/yyhuni/lunafox/server/internal/modules/asset/repository"
)

type assetSubdomainStoreAdapter struct {
	repo *assetrepo.SubdomainRepository
}

func newAssetSubdomainStoreAdapter(repo *assetrepo.SubdomainRepository) *assetSubdomainStoreAdapter {
	return &assetSubdomainStoreAdapter{repo: repo}
}

func (adapter *assetSubdomainStoreAdapter) FindByTargetID(targetID int, page, pageSize int, filter string) ([]assetdomain.Subdomain, int64, error) {
	items, total, err := adapter.repo.FindByTargetID(targetID, page, pageSize, filter)
	if err != nil {
		return nil, 0, err
	}
	results := make([]assetdomain.Subdomain, 0, len(items))
	for index := range items {
		results = append(results, *assetModelSubdomainToDomain(&items[index]))
	}
	return results, total, nil
}

func (adapter *assetSubdomainStoreAdapter) BulkCreate(items []assetdomain.Subdomain) (int, error) {
	return adapter.repo.BulkCreate(assetDomainSubdomainListToModel(items))
}

func (adapter *assetSubdomainStoreAdapter) BulkDelete(ids []int) (int64, error) {
	return adapter.repo.BulkDelete(ids)
}

func (adapter *assetSubdomainStoreAdapter) StreamByTargetID(targetID int) (*sql.Rows, error) {
	return adapter.repo.StreamByTargetID(targetID)
}

func (adapter *assetSubdomainStoreAdapter) CountByTargetID(targetID int) (int64, error) {
	return adapter.repo.CountByTargetID(targetID)
}

func (adapter *assetSubdomainStoreAdapter) ScanRow(rows *sql.Rows) (*assetdomain.Subdomain, error) {
	item, err := adapter.repo.ScanRow(rows)
	if err != nil {
		return nil, err
	}
	return assetModelSubdomainToDomain(item), nil
}
