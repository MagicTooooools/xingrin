package application

import (
	"context"

	scandomain "github.com/yyhuni/lunafox/server/internal/modules/scan/domain"
)

type ScanListFilter struct {
	Page     int
	PageSize int
	TargetID int
	Status   string
	Search   string
}

type QueryTargetRef = scandomain.QueryTargetRef

type QueryScan = scandomain.QueryScan

type QueryStatistics = scandomain.QueryStatistics

type ScanQueryService struct{ store ScanQueryStore }

func NewScanQueryService(store ScanQueryStore) *ScanQueryService {
	return &ScanQueryService{store: store}
}

func (service *ScanQueryService) ListScans(ctx context.Context, filter ScanListFilter) ([]QueryScan, int64, error) {
	_ = ctx
	return service.store.FindAll(filter.Page, filter.PageSize, filter.TargetID, filter.Status, filter.Search)
}

func (service *ScanQueryService) GetScanByID(ctx context.Context, id int) (*QueryScan, error) {
	_ = ctx
	return service.store.GetQueryByID(id)
}

func (service *ScanQueryService) GetStatistics(ctx context.Context) (*QueryStatistics, error) {
	_ = ctx
	return service.store.GetStatistics()
}
