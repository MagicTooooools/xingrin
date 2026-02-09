package catalogwiring

import (
	catalogrepo "github.com/yyhuni/lunafox/server/internal/modules/catalog/repository"
	identityrepo "github.com/yyhuni/lunafox/server/internal/modules/identity/repository"
)

func NewEngineStoreAdapter(repo *catalogrepo.EngineRepository) *catalogEngineStoreAdapter {
	return newCatalogEngineStoreAdapter(repo)
}

func NewTargetStoreAdapter(repo *catalogrepo.TargetRepository) *catalogTargetStoreAdapter {
	return newCatalogTargetStoreAdapter(repo)
}

func NewOrganizationStoreAdapter(repo *identityrepo.OrganizationRepository) *catalogOrganizationStoreAdapter {
	return newCatalogOrganizationStoreAdapter(repo)
}

func NewWordlistStoreAdapter(repo *catalogrepo.WordlistRepository) *catalogWordlistStoreAdapter {
	return newCatalogWordlistStoreAdapter(repo)
}
