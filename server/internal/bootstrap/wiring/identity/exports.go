package identitywiring

import (
	identityapp "github.com/yyhuni/lunafox/server/internal/modules/identity/application"
	identityrepo "github.com/yyhuni/lunafox/server/internal/modules/identity/repository"
)

func NewIdentityUserStoreAdapter(repo *identityrepo.UserRepository) identityapp.UserStore {
	return newIdentityUserStoreAdapter(repo)
}

func NewIdentityAuthUserStoreAdapter(repo *identityrepo.UserRepository) identityapp.AuthUserStore {
	return newIdentityUserStoreAdapter(repo)
}

func NewIdentityOrganizationStoreAdapter(repo *identityrepo.OrganizationRepository) identityapp.OrganizationStore {
	return newIdentityOrganizationStoreAdapter(repo)
}
