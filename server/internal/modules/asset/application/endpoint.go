package application

import (
	"context"
	"database/sql"
	"errors"

	assetdomain "github.com/yyhuni/lunafox/server/internal/modules/asset/domain"
	"gorm.io/gorm"
)

var ErrEndpointTargetNotFound = errors.New("target not found")

type EndpointQueryStore interface {
	FindByTargetID(targetID int, page, pageSize int, filter string) ([]assetdomain.Endpoint, int64, error)
	FindByID(id int) (*assetdomain.Endpoint, error)
	StreamByTargetID(targetID int) (*sql.Rows, error)
	CountByTargetID(targetID int) (int64, error)
	ScanRow(rows *sql.Rows) (*assetdomain.Endpoint, error)
}

type EndpointTargetLookup interface {
	FindByID(id int) (*assetdomain.TargetRef, error)
}

type EndpointQueryService struct {
	store        EndpointQueryStore
	targetLookup EndpointTargetLookup
}

func NewEndpointQueryService(store EndpointQueryStore, targetLookup EndpointTargetLookup) *EndpointQueryService {
	return &EndpointQueryService{store: store, targetLookup: targetLookup}
}

func (service *EndpointQueryService) ListByTarget(ctx context.Context, targetID, page, pageSize int, filter string) ([]assetdomain.Endpoint, int64, error) {
	_ = ctx

	if _, err := service.targetLookup.FindByID(targetID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, ErrEndpointTargetNotFound
		}
		return nil, 0, err
	}

	return service.store.FindByTargetID(targetID, page, pageSize, filter)
}

func (service *EndpointQueryService) GetByID(ctx context.Context, id int) (*assetdomain.Endpoint, error) {
	_ = ctx
	return service.store.FindByID(id)
}

func (service *EndpointQueryService) StreamByTarget(ctx context.Context, targetID int) (*sql.Rows, error) {
	_ = ctx

	if _, err := service.targetLookup.FindByID(targetID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrEndpointTargetNotFound
		}
		return nil, err
	}

	return service.store.StreamByTargetID(targetID)
}

func (service *EndpointQueryService) CountByTarget(ctx context.Context, targetID int) (int64, error) {
	_ = ctx

	if _, err := service.targetLookup.FindByID(targetID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, ErrEndpointTargetNotFound
		}
		return 0, err
	}

	return service.store.CountByTargetID(targetID)
}

func (service *EndpointQueryService) ScanRow(ctx context.Context, rows *sql.Rows) (*assetdomain.Endpoint, error) {
	_ = ctx
	return service.store.ScanRow(rows)
}

var ErrEndpointNotFound = errors.New("endpoint not found")

type EndpointUpsertItem struct {
	URL               string
	Host              string
	Location          string
	Title             string
	Webserver         string
	ContentType       string
	StatusCode        *int
	ContentLength     *int
	ResponseBody      string
	Tech              []string
	Vhost             *bool
	ResponseHeaders   string
}

type EndpointCommandStore interface {
	FindByID(id int) (*assetdomain.Endpoint, error)
	BulkCreate(endpoints []assetdomain.Endpoint) (int, error)
	Delete(id int) error
	BulkDelete(ids []int) (int64, error)
	BulkUpsert(endpoints []assetdomain.Endpoint) (int64, error)
}

type EndpointCommandService struct {
	store        EndpointCommandStore
	targetLookup EndpointTargetLookup
}

func NewEndpointCommandService(store EndpointCommandStore, targetLookup EndpointTargetLookup) *EndpointCommandService {
	return &EndpointCommandService{store: store, targetLookup: targetLookup}
}

func (service *EndpointCommandService) BulkCreate(ctx context.Context, targetID int, urls []string) (int, error) {
	_ = ctx

	target, err := service.targetLookup.FindByID(targetID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, ErrEndpointTargetNotFound
		}
		return 0, err
	}

	endpoints := make([]assetdomain.Endpoint, 0, len(urls))
	for _, rawURL := range urls {
		if assetdomain.IsURLMatchTarget(rawURL, *target) {
			endpoints = append(endpoints, assetdomain.Endpoint{
				TargetID: targetID,
				URL:      rawURL,
				Host:     assetdomain.ExtractHostFromURL(rawURL),
			})
		}
	}

	if len(endpoints) == 0 {
		return 0, nil
	}

	return service.store.BulkCreate(endpoints)
}

func (service *EndpointCommandService) Delete(ctx context.Context, id int) error {
	_ = ctx

	if _, err := service.store.FindByID(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrEndpointNotFound
		}
		return err
	}

	return service.store.Delete(id)
}

func (service *EndpointCommandService) BulkDelete(ctx context.Context, ids []int) (int64, error) {
	_ = ctx

	if len(ids) == 0 {
		return 0, nil
	}

	return service.store.BulkDelete(ids)
}

func (service *EndpointCommandService) BulkUpsert(ctx context.Context, targetID int, items []EndpointUpsertItem) (int64, error) {
	_ = ctx

	target, err := service.targetLookup.FindByID(targetID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, ErrEndpointTargetNotFound
		}
		return 0, err
	}

	endpoints := make([]assetdomain.Endpoint, 0, len(items))
	for _, item := range items {
		if !assetdomain.IsURLMatchTarget(item.URL, *target) {
			continue
		}

		host := item.Host
		if host == "" {
			host = assetdomain.ExtractHostFromURL(item.URL)
		}

		endpoints = append(endpoints, assetdomain.Endpoint{
			TargetID:          targetID,
			URL:               item.URL,
			Host:              host,
			Location:          item.Location,
			Title:             item.Title,
			Webserver:         item.Webserver,
			ContentType:       item.ContentType,
			StatusCode:        item.StatusCode,
			ContentLength:     item.ContentLength,
			ResponseBody:      item.ResponseBody,
			Tech:              item.Tech,
			Vhost:             item.Vhost,
			ResponseHeaders:   item.ResponseHeaders,
		})
	}

	if len(endpoints) == 0 {
		return 0, nil
	}

	return service.store.BulkUpsert(endpoints)
}
