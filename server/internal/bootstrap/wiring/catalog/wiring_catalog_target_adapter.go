package catalogwiring

import (
	catalogdomain "github.com/yyhuni/lunafox/server/internal/modules/catalog/domain"
	catalogrepo "github.com/yyhuni/lunafox/server/internal/modules/catalog/repository"
)

type catalogTargetStoreAdapter struct {
	repo *catalogrepo.TargetRepository
}

func newCatalogTargetStoreAdapter(repo *catalogrepo.TargetRepository) *catalogTargetStoreAdapter {
	return &catalogTargetStoreAdapter{repo: repo}
}

func (adapter *catalogTargetStoreAdapter) FindByID(id int) (*catalogdomain.Target, error) {
	target, err := adapter.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return catalogModelTargetToDomain(target), nil
}

func (adapter *catalogTargetStoreAdapter) FindAll(page, pageSize int, targetType, filter string) ([]catalogdomain.Target, int64, error) {
	targets, total, err := adapter.repo.FindAll(page, pageSize, targetType, filter)
	if err != nil {
		return nil, 0, err
	}
	return catalogModelTargetListToDomain(targets), total, nil
}

func (adapter *catalogTargetStoreAdapter) GetAssetCounts(targetID int) (*catalogdomain.TargetAssetCounts, error) {
	counts, err := adapter.repo.GetAssetCounts(targetID)
	if err != nil {
		return nil, err
	}
	return &catalogdomain.TargetAssetCounts{
		Subdomains:  counts.Subdomains,
		Websites:    counts.Websites,
		Endpoints:   counts.Endpoints,
		IPs:         counts.IPs,
		Directories: counts.Directories,
		Screenshots: counts.Screenshots,
	}, nil
}

func (adapter *catalogTargetStoreAdapter) GetVulnerabilityCounts(targetID int) (*catalogdomain.VulnerabilityCounts, error) {
	counts, err := adapter.repo.GetVulnerabilityCounts(targetID)
	if err != nil {
		return nil, err
	}
	return &catalogdomain.VulnerabilityCounts{
		Total:    counts.Total,
		Critical: counts.Critical,
		High:     counts.High,
		Medium:   counts.Medium,
		Low:      counts.Low,
	}, nil
}

func (adapter *catalogTargetStoreAdapter) ExistsByName(name string, excludeID ...int) (bool, error) {
	return adapter.repo.ExistsByName(name, excludeID...)
}

func (adapter *catalogTargetStoreAdapter) Create(target *catalogdomain.Target) error {
	modelTarget := catalogDomainTargetToModel(target)
	if err := adapter.repo.Create(modelTarget); err != nil {
		return err
	}
	*target = *catalogModelTargetToDomain(modelTarget)
	return nil
}

func (adapter *catalogTargetStoreAdapter) Update(target *catalogdomain.Target) error {
	return adapter.repo.Update(catalogDomainTargetToModel(target))
}

func (adapter *catalogTargetStoreAdapter) SoftDelete(id int) error {
	return adapter.repo.SoftDelete(id)
}

func (adapter *catalogTargetStoreAdapter) BulkSoftDelete(ids []int) (int64, error) {
	return adapter.repo.BulkSoftDelete(ids)
}

func (adapter *catalogTargetStoreAdapter) BulkCreateIgnoreConflicts(targets []catalogdomain.Target) (int, error) {
	modelTargets := catalogDomainTargetListToModel(targets)
	return adapter.repo.BulkCreateIgnoreConflicts(modelTargets)
}

func (adapter *catalogTargetStoreAdapter) FindByNames(names []string) ([]catalogdomain.Target, error) {
	targets, err := adapter.repo.FindByNames(names)
	if err != nil {
		return nil, err
	}
	return catalogModelTargetListToDomain(targets), nil
}
