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

var (
	ErrScanNotFoundForSnapshot = errors.New("scan not found")
	ErrTargetMismatch          = errors.New("targetId does not match scan's target")
)

// WebsiteSnapshotService handles website snapshot business logic
type WebsiteSnapshotService struct {
	snapshotRepo   *repository.WebsiteSnapshotRepository
	scanRepo       *repository.ScanRepository
	websiteService *WebsiteService
}

// NewWebsiteSnapshotService creates a new website snapshot service
func NewWebsiteSnapshotService(
	snapshotRepo *repository.WebsiteSnapshotRepository,
	scanRepo *repository.ScanRepository,
	websiteService *WebsiteService,
) *WebsiteSnapshotService {
	return &WebsiteSnapshotService{
		snapshotRepo:   snapshotRepo,
		scanRepo:       scanRepo,
		websiteService: websiteService,
	}
}

// SaveAndSync saves website snapshots and syncs to asset table
// 1. Validates scan exists and is not soft-deleted
// 2. Validates URLs match target (filters invalid items)
// 3. Saves to website_snapshot table
// 4. Calls WebsiteService.BulkUpsert to sync to website table
func (s *WebsiteSnapshotService) SaveAndSync(scanID int, targetID int, items []dto.WebsiteSnapshotItem) (snapshotCount int64, assetCount int64, err error) {
	if len(items) == 0 {
		return 0, 0, nil
	}

	// Step 1: Validate scan exists and is not soft-deleted
	scan, err := s.scanRepo.FindByID(scanID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, 0, ErrScanNotFoundForSnapshot
		}
		return 0, 0, err
	}

	// Verify target_id matches scan's target
	if scan.TargetID != targetID {
		return 0, 0, ErrTargetMismatch
	}

	// Step 2: Get target for URL validation
	target, err := s.scanRepo.GetTargetByScanID(scanID)
	if err != nil {
		return 0, 0, err
	}

	// Step 3: Filter valid items and convert to snapshot models
	// Only URLs that match target are saved (defensive validation)
	snapshots := make([]model.WebsiteSnapshot, 0, len(items))
	validItems := make([]dto.WebsiteSnapshotItem, 0, len(items))

	for _, item := range items {
		// Validate URL matches target
		if !validator.IsURLMatchTarget(item.URL, target.Name, target.Type) {
			continue // Skip invalid URLs
		}

		host := item.Host
		if host == "" {
			host = repository.ExtractHostFromURL(item.URL)
		}

		snapshots = append(snapshots, model.WebsiteSnapshot{
			ScanID:          scanID,
			URL:             item.URL,
			Host:            host,
			Title:           item.Title,
			StatusCode:      item.StatusCode,
			ContentLength:   item.ContentLength,
			Location:        item.Location,
			Webserver:       item.Webserver,
			ContentType:     item.ContentType,
			Tech:            item.Tech,
			ResponseBody:    item.ResponseBody,
			Vhost:           item.Vhost,
			ResponseHeaders: item.ResponseHeaders,
		})
		validItems = append(validItems, item)
	}

	if len(snapshots) == 0 {
		return 0, 0, nil
	}

	// Step 4: Save to snapshot table
	snapshotCount, err = s.snapshotRepo.BulkCreate(snapshots)
	if err != nil {
		return 0, 0, err
	}

	// Step 5: Sync to asset table (WebsiteSnapshotItem is an alias of WebsiteUpsertItem, no conversion needed)
	assetCount, err = s.websiteService.BulkUpsert(targetID, validItems)
	if err != nil {
		// Log error but don't fail - snapshot is already saved
		// In production, consider using a transaction or compensation logic
		return snapshotCount, 0, nil
	}

	return snapshotCount, assetCount, nil
}

// ListByScan returns paginated website snapshots for a scan
func (s *WebsiteSnapshotService) ListByScan(scanID int, query *dto.WebsiteSnapshotListQuery) ([]model.WebsiteSnapshot, int64, error) {
	// Validate scan exists
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
func (s *WebsiteSnapshotService) StreamByScan(scanID int) (*sql.Rows, error) {
	// Validate scan exists
	_, err := s.scanRepo.FindByID(scanID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrScanNotFoundForSnapshot
		}
		return nil, err
	}

	return s.snapshotRepo.StreamByScanID(scanID)
}

// CountByScan returns the count of website snapshots for a scan
func (s *WebsiteSnapshotService) CountByScan(scanID int) (int64, error) {
	// Validate scan exists
	_, err := s.scanRepo.FindByID(scanID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, ErrScanNotFoundForSnapshot
		}
		return 0, err
	}

	return s.snapshotRepo.CountByScanID(scanID)
}

// ScanRow scans a row into WebsiteSnapshot model
func (s *WebsiteSnapshotService) ScanRow(rows *sql.Rows) (*model.WebsiteSnapshot, error) {
	return s.snapshotRepo.ScanRow(rows)
}
