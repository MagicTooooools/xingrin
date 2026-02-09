package catalogwiring

import identityrepo "github.com/yyhuni/lunafox/server/internal/modules/identity/repository"

type catalogOrganizationStoreAdapter struct {
	repo *identityrepo.OrganizationRepository
}

func newCatalogOrganizationStoreAdapter(repo *identityrepo.OrganizationRepository) *catalogOrganizationStoreAdapter {
	return &catalogOrganizationStoreAdapter{repo: repo}
}

func (adapter *catalogOrganizationStoreAdapter) Exists(id int) (bool, error) {
	return adapter.repo.Exists(id)
}

func (adapter *catalogOrganizationStoreAdapter) BulkAddTargets(organizationID int, targetIDs []int) error {
	return adapter.repo.BulkAddTargets(organizationID, targetIDs)
}
