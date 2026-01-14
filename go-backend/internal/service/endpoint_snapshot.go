package service

import (
	"database/sql"
	"errors"

	"github.com/xingrin/go-backend/internal/dto"
	"github.com/xingrin/go-backend/internal/model"
	"github.com/xingrin/go-backend/internal/pkg/validator"
	"github.com/xingrin/go-backend/internal/repository"
	"gorm.io/gorm"
)

// EndpointSnapshotService handles endpoint snapshot business logic
type EndpointSnapshotService struct {
	snapshotRepo    *repository.EndpointSnapshotRepository
	scanRepo        *repository.ScanRepository
	endpointService *EndpointService
}

// NewEndpointSnapshotService creates a new endpoint snapshot service
func NewEndpointSnapshotService(
	snapshotRepo *repository.EndpointSnapshotRepository,
	scanRepo *repository.ScanRepository,
	endpointService *EndpointService,
) *EndpointSnapshotService {
	return &EndpointSnapshotService{
		snapshotRepo:    snapshotRepo,
		scanRepo:        scanRepo,
		endpointService: endpointService,
	}
}

// SaveAndSync saves endpoint snapshots and syncs to asset table
func (s *EndpointSnapshotService) SaveAndSync(scanID int, targetID int, items []dto.EndpointSnapshotItem) (snapshotCount int64, assetCount int64, err error) {
	if len(items) == 0 {
		return 0, 0, nil
	}

	// Validate scan exists
	scan, err := s.scanRepo.FindByID(scanID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, 0, ErrScanNotFoundForSnapshot
		}
		return 0, 0, err
	}

	if scan.TargetID != targetID {
		return 0, 0, errors.New("target_id does not match scan's target")
	}

	// Get target for validation
	target, err := s.scanRepo.GetTargetByScanID(scanID)
	if err != nil {
		return 0, 0, err
	}

	// Filter valid endpoints
	snapshots := make([]model.EndpointSnapshot, 0, len(items))
	validItems := make([]dto.EndpointUpsertItem, 0, len(items))

	for _, item := range items {
		if !validator.IsURLMatchTarget(item.URL, target.Name, target.Type) {
			continue
		}

		host := item.Host
		if host == "" {
			host = repository.ExtractHostFromURL(item.URL)
		}

		snapshots = append(snapshots, model.EndpointSnapshot{
			ScanID:            scanID,
			URL:               item.URL,
			Host:              host,
			Title:             item.Title,
			StatusCode:        item.StatusCode,
			ContentLength:     item.ContentLength,
			Location:          item.Location,
			Webserver:         item.Webserver,
			ContentType:       item.ContentType,
			Tech:              item.Tech,
			ResponseBody:      item.ResponseBody,
			Vhost:             item.Vhost,
			MatchedGFPatterns: item.MatchedGFPatterns,
			ResponseHeaders:   item.ResponseHeaders,
		})

		validItems = append(validItems, dto.EndpointUpsertItem{
			URL:               item.URL,
			Host:              item.Host,
			Title:             item.Title,
			StatusCode:        item.StatusCode,
			ContentLength:     item.ContentLength,
			Location:          item.Location,
			Webserver:         item.Webserver,
			ContentType:       item.ContentType,
			Tech:              item.Tech,
			ResponseBody:      item.ResponseBody,
			Vhost:             item.Vhost,
			MatchedGFPatterns: item.MatchedGFPatterns,
			ResponseHeaders:   item.ResponseHeaders,
		})
	}

	if len(snapshots) == 0 {
		return 0, 0, nil
	}

	// Save to snapshot table
	snapshotCount, err = s.snapshotRepo.BulkCreate(snapshots)
	if err != nil {
		return 0, 0, err
	}

	// Sync to asset table
	assetCount, err = s.endpointService.BulkUpsert(targetID, validItems)
	if err != nil {
		return snapshotCount, 0, nil
	}

	return snapshotCount, assetCount, nil
}

// ListByScan returns paginated endpoint snapshots for a scan
func (s *EndpointSnapshotService) ListByScan(scanID int, query *dto.EndpointSnapshotListQuery) ([]model.EndpointSnapshot, int64, error) {
	_, err := s.scanRepo.FindByID(scanID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, ErrScanNotFoundForSnapshot
		}
		return nil, 0, err
	}

	return s.snapshotRepo.FindByScanID(scanID, query.GetPage(), query.GetPageSize(), query.Filter)
}

// StreamByScan returns a cursor for streaming export
func (s *EndpointSnapshotService) StreamByScan(scanID int) (*sql.Rows, error) {
	_, err := s.scanRepo.FindByID(scanID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrScanNotFoundForSnapshot
		}
		return nil, err
	}

	return s.snapshotRepo.StreamByScanID(scanID)
}

// CountByScan returns the count of endpoint snapshots for a scan
func (s *EndpointSnapshotService) CountByScan(scanID int) (int64, error) {
	_, err := s.scanRepo.FindByID(scanID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, ErrScanNotFoundForSnapshot
		}
		return 0, err
	}

	return s.snapshotRepo.CountByScanID(scanID)
}

// ScanRow scans a row into EndpointSnapshot model
func (s *EndpointSnapshotService) ScanRow(rows *sql.Rows) (*model.EndpointSnapshot, error) {
	return s.snapshotRepo.ScanRow(rows)
}
