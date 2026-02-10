package catalogwiring

import (
	catalogapp "github.com/yyhuni/lunafox/server/internal/modules/catalog/application"
	catalogrepo "github.com/yyhuni/lunafox/server/internal/modules/catalog/repository"
	identityrepo "github.com/yyhuni/lunafox/server/internal/modules/identity/repository"
)

func NewCatalogEngineStoreAdapter(repo *catalogrepo.EngineRepository) catalogapp.EngineStore {
	return newCatalogEngineStoreAdapter(repo)
}

func NewCatalogTargetStoreAdapter(repo *catalogrepo.TargetRepository) catalogapp.TargetStore {
	return newCatalogTargetStoreAdapter(repo)
}

func NewCatalogOrganizationStoreAdapter(repo *identityrepo.OrganizationRepository) catalogapp.OrganizationStore {
	return newCatalogOrganizationStoreAdapter(repo)
}

func NewCatalogWordlistStoreAdapter(repo *catalogrepo.WordlistRepository) catalogapp.WordlistStore {
	return newCatalogWordlistStoreAdapter(repo)
}
