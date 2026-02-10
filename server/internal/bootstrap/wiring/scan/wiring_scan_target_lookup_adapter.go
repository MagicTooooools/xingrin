package scanwiring

import (
	catalogrepo "github.com/yyhuni/lunafox/server/internal/modules/catalog/repository"
	scanapp "github.com/yyhuni/lunafox/server/internal/modules/scan/application"
	"github.com/yyhuni/lunafox/server/internal/pkg/timeutil"
)

type scanCreateTargetLookupAdapter struct {
	repo *catalogrepo.TargetRepository
}

func newScanCreateTargetLookupAdapter(repo *catalogrepo.TargetRepository) *scanCreateTargetLookupAdapter {
	return &scanCreateTargetLookupAdapter{repo: repo}
}

func (adapter *scanCreateTargetLookupAdapter) GetActiveByID(id int) (*scanapp.TargetRef, error) {
	target, err := adapter.repo.GetActiveByID(id)
	if err != nil {
		return nil, err
	}
	return &scanapp.TargetRef{ID: target.ID, Name: target.Name, Type: target.Type, CreatedAt: timeutil.ToUTC(target.CreatedAt), LastScannedAt: timeutil.ToUTCPtr(target.LastScannedAt), DeletedAt: timeutil.ToUTCPtr(target.DeletedAt)}, nil
}
