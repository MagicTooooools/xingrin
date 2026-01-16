package service

import (
	"errors"

	"github.com/orbit/server/internal/dto"
	"github.com/orbit/server/internal/model"
	"github.com/orbit/server/internal/pkg/validator"
	"github.com/orbit/server/internal/repository"
	"gorm.io/gorm"
)

var (
	ErrScreenshotSnapshotNotFound = errors.New("screenshot snapshot not found")
)

// ScreenshotSnapshotService handles screenshot snapshot business logic
type ScreenshotSnapshotService struct {
	snapshotRepo      *repository.ScreenshotSnapshotRepository
	scanRepo          *repository.ScanRepository
	screenshotService *ScreenshotService
}

// NewScreenshotSnapshotService creates a new screenshot snapshot service
func NewScreenshotSnapshotService(
	snapshotRepo *repository.ScreenshotSnapshotRepository,
	scanRepo *repository.ScanRepository,
	screenshotService *ScreenshotService,
) *ScreenshotSnapshotService {
	return &ScreenshotSnapshotService{
		snapshotRepo:      snapshotRepo,
		scanRepo:          scanRepo,
		screenshotService: screenshotService,
	}
}

// SaveAndSync saves screenshot snapshots and syncs to asset table
// 1. Validates scan exists and is not soft-deleted
// 2. Validates URLs match target (filters invalid items)
// 3. Upserts into screenshot_snapshot table
// 4. Calls ScreenshotService.BulkUpsert to sync to screenshot table
func (s *ScreenshotSnapshotService) SaveAndSync(scanID int, targetID int, items []dto.ScreenshotSnapshotItem) (snapshotCount int64, assetCount int64, err error) {
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
		return 0, 0, ErrTargetMismatch
	}

	// Get target for validation
	target, err := s.scanRepo.GetTargetByScanID(scanID)
	if err != nil {
		return 0, 0, err
	}

	// Filter valid items
	snapshots := make([]model.ScreenshotSnapshot, 0, len(items)) // snapshot table models
	assetItems := make([]dto.ScreenshotItem, 0, len(items))      // asset table items

	for _, item := range items {
		if !validator.IsURLMatchTarget(item.URL, target.Name, target.Type) {
			continue
		}

		snapshots = append(snapshots, model.ScreenshotSnapshot{
			ScanID:     scanID,
			URL:        item.URL,
			StatusCode: item.StatusCode,
			Image:      item.Image,
		})

		assetItems = append(assetItems, dto.ScreenshotItem(item))
	}

	if len(snapshots) == 0 {
		return 0, 0, nil
	}

	// Save to snapshot table
	snapshotCount, err = s.snapshotRepo.BulkUpsert(snapshots)
	if err != nil {
		return 0, 0, err
	}

	// Sync to asset table
	assetCount, err = s.screenshotService.BulkUpsert(targetID, &dto.BulkUpsertScreenshotRequest{Screenshots: assetItems})
	if err != nil {
		// Snapshot is already saved; don't fail the request
		return snapshotCount, 0, nil
	}

	return snapshotCount, assetCount, nil
}

// ListByScan returns paginated screenshot snapshots for a scan
func (s *ScreenshotSnapshotService) ListByScan(scanID int, query *dto.ScreenshotSnapshotListQuery) ([]model.ScreenshotSnapshot, int64, error) {
	_, err := s.scanRepo.FindByID(scanID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, ErrScanNotFoundForSnapshot
		}
		return nil, 0, err
	}

	return s.snapshotRepo.FindByScanID(scanID, query.GetPage(), query.GetPageSize(), query.Filter)
}

// GetByID returns a screenshot snapshot by ID under a scan (including image data)
func (s *ScreenshotSnapshotService) GetByID(scanID int, id int) (*model.ScreenshotSnapshot, error) {
	// Validate scan exists
	_, err := s.scanRepo.FindByID(scanID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrScanNotFoundForSnapshot
		}
		return nil, err
	}

	snapshot, err := s.snapshotRepo.FindByIDAndScanID(id, scanID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrScreenshotSnapshotNotFound
		}
		return nil, err
	}

	return snapshot, nil
}
