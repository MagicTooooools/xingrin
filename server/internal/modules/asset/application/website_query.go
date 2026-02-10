package application

import (
	"context"
	"database/sql"
	"errors"

	assetdomain "github.com/yyhuni/lunafox/server/internal/modules/asset/domain"
	"github.com/yyhuni/lunafox/server/internal/pkg/dberrors"
)

var ErrWebsiteTargetNotFound = errors.New("target not found")

type WebsiteQueryStore interface {
	FindByTargetID(targetID int, page, pageSize int, filter string) ([]assetdomain.Website, int64, error)
	StreamByTargetID(targetID int) (*sql.Rows, error)
	CountByTargetID(targetID int) (int64, error)
	ScanRow(rows *sql.Rows) (*assetdomain.Website, error)
}

type WebsiteTargetLookup interface {
	GetActiveByID(id int) (*assetdomain.TargetRef, error)
}

type WebsiteQueryService struct {
	store        WebsiteQueryStore
	targetLookup WebsiteTargetLookup
}

func NewWebsiteQueryService(store WebsiteQueryStore, targetLookup WebsiteTargetLookup) *WebsiteQueryService {
	return &WebsiteQueryService{store: store, targetLookup: targetLookup}
}

func (service *WebsiteQueryService) ListByTarget(ctx context.Context, targetID, page, pageSize int, filter string) ([]assetdomain.Website, int64, error) {
	_ = ctx
	return service.store.FindByTargetID(targetID, page, pageSize, filter)
}

func (service *WebsiteQueryService) StreamByTarget(ctx context.Context, targetID int) (*sql.Rows, error) {
	_ = ctx

	if _, err := service.targetLookup.GetActiveByID(targetID); err != nil {
		if dberrors.IsRecordNotFound(err) {
			return nil, ErrWebsiteTargetNotFound
		}
		return nil, err
	}

	return service.store.StreamByTargetID(targetID)
}

func (service *WebsiteQueryService) CountByTarget(ctx context.Context, targetID int) (int64, error) {
	_ = ctx

	if _, err := service.targetLookup.GetActiveByID(targetID); err != nil {
		if dberrors.IsRecordNotFound(err) {
			return 0, ErrWebsiteTargetNotFound
		}
		return 0, err
	}

	return service.store.CountByTargetID(targetID)
}

func (service *WebsiteQueryService) ScanRow(ctx context.Context, rows *sql.Rows) (*assetdomain.Website, error) {
	_ = ctx
	return service.store.ScanRow(rows)
}
