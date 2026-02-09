package catalogwiring

import (
	catalogdomain "github.com/yyhuni/lunafox/server/internal/modules/catalog/domain"
	catalogrepo "github.com/yyhuni/lunafox/server/internal/modules/catalog/repository"
)

type catalogEngineStoreAdapter struct {
	repo *catalogrepo.EngineRepository
}

func newCatalogEngineStoreAdapter(repo *catalogrepo.EngineRepository) *catalogEngineStoreAdapter {
	return &catalogEngineStoreAdapter{repo: repo}
}

func (adapter *catalogEngineStoreAdapter) FindByID(id int) (*catalogdomain.ScanEngine, error) {
	engine, err := adapter.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return catalogModelEngineToDomain(engine), nil
}

func (adapter *catalogEngineStoreAdapter) FindAll(page, pageSize int) ([]catalogdomain.ScanEngine, int64, error) {
	engines, total, err := adapter.repo.FindAll(page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	return catalogModelEngineListToDomain(engines), total, nil
}

func (adapter *catalogEngineStoreAdapter) ExistsByName(name string, excludeID ...int) (bool, error) {
	return adapter.repo.ExistsByName(name, excludeID...)
}

func (adapter *catalogEngineStoreAdapter) Create(engine *catalogdomain.ScanEngine) error {
	modelEngine := catalogDomainEngineToModel(engine)
	if err := adapter.repo.Create(modelEngine); err != nil {
		return err
	}
	*engine = *catalogModelEngineToDomain(modelEngine)
	return nil
}

func (adapter *catalogEngineStoreAdapter) Update(engine *catalogdomain.ScanEngine) error {
	return adapter.repo.Update(catalogDomainEngineToModel(engine))
}

func (adapter *catalogEngineStoreAdapter) Delete(id int) error {
	return adapter.repo.Delete(id)
}
