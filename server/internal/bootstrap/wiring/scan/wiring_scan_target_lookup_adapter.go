package scanwiring

import (
	catalogrepo "github.com/yyhuni/lunafox/server/internal/modules/catalog/repository"
	scanapp "github.com/yyhuni/lunafox/server/internal/modules/scan/application"
)

type scanCreateTargetLookupAdapter struct {
	repo *catalogrepo.TargetRepository
}

func newScanCreateTargetLookupAdapter(repo *catalogrepo.TargetRepository) *scanCreateTargetLookupAdapter {
	return &scanCreateTargetLookupAdapter{repo: repo}
}

func (adapter *scanCreateTargetLookupAdapter) FindByID(id int) (*scanapp.TargetRef, error) {
	target, err := adapter.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return &scanapp.TargetRef{ID: target.ID, Name: target.Name, Type: target.Type, CreatedAt: target.CreatedAt, LastScannedAt: target.LastScannedAt, DeletedAt: target.DeletedAt}, nil
}
