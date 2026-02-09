package identitywiring

import identityrepo "github.com/yyhuni/lunafox/server/internal/modules/identity/repository"

func NewUserStoreAdapter(repo *identityrepo.UserRepository) *identityUserStoreAdapter {
	return newIdentityUserStoreAdapter(repo)
}

func NewOrganizationStoreAdapter(repo *identityrepo.OrganizationRepository) *identityOrganizationStoreAdapter {
	return newIdentityOrganizationStoreAdapter(repo)
}
