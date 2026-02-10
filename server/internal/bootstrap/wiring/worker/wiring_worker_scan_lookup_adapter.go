package workerwiring

import (
	"errors"

	catalogapp "github.com/yyhuni/lunafox/server/internal/modules/catalog/application"
	scanrepo "github.com/yyhuni/lunafox/server/internal/modules/scan/repository"
	"gorm.io/gorm"
)

type workerScanLookupAdapter struct {
	repo *scanrepo.ScanRepository
}

func (adapter *workerScanLookupAdapter) GetActiveByID(id int) (*catalogapp.WorkerScanRef, error) {
	scan, err := adapter.repo.GetActiveByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, catalogapp.ErrWorkerScanNotFound
		}
		return nil, err
	}
	return &catalogapp.WorkerScanRef{ID: scan.ID}, nil
}
