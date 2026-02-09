package securitywiring

import catalogrepo "github.com/yyhuni/lunafox/server/internal/modules/catalog/repository"

func NewTargetLookupAdapter(repo *catalogrepo.TargetRepository) *securityTargetLookupAdapter {
	return newSecurityTargetLookupAdapter(repo)
}
