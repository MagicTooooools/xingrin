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

// DirectorySnapshotService handles directory snapshot business logic
type DirectorySnapshotService struct {
	snapshotRepo     *repository.DirectorySnapshotRepository
	scanRepo         *repository.ScanRepository
	directoryService *DirectoryService
}

// NewDirectorySnapshotService creates a new directory snapshot service
func NewDirectorySnapshotService(
	snapshotRepo *repository.DirectorySnapshotRepository,
	scanRepo *repository.ScanRepository,
	directoryService *DirectoryService,
) *DirectorySnapshotService {
	return &DirectorySnapshotService{
		snapshotRepo:     snapshotRepo,
		scanRepo:         scanRepo,
		directoryService: directoryService,
	}
}

// SaveAndSync saves directory snapshots and syncs to asset table
func (s *DirectorySnapshotService) SaveAndSync(scanID int, targetID int, items []dto.DirectorySnapshotItem) (snapshotCount int64, assetCount int64, err error) {
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

	// Filter valid directories
	snapshots := make([]model.DirectorySnapshot, 0, len(items))
	validItems := make([]dto.DirectoryUpsertItem, 0, len(items))

	for _, item := range items {
		if !validator.IsURLMatchTarget(item.URL, target.Name, target.Type) {
			continue
		}

		snapshots = append(snapshots, model.DirectorySnapshot{
			ScanID:        scanID,
			URL:           item.URL,
			Status:        item.Status,
			ContentLength: item.ContentLength,
			ContentType:   item.ContentType,
			Duration:      item.Duration,
		})

		validItems = append(validItems, dto.DirectoryUpsertItem{
			URL:           item.URL,
			Status:        item.Status,
			ContentLength: item.ContentLength,
			ContentType:   item.ContentType,
			Duration:      item.Duration,
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
	assetCount, err = s.directoryService.BulkUpsert(targetID, validItems)
	if err != nil {
		return snapshotCount, 0, nil
	}

	return snapshotCount, assetCount, nil
}

// ListByScan returns paginated directory snapshots for a scan
func (s *DirectorySnapshotService) ListByScan(scanID int, query *dto.DirectorySnapshotListQuery) ([]model.DirectorySnapshot, int64, error) {
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
func (s *DirectorySnapshotService) StreamByScan(scanID int) (*sql.Rows, error) {
	_, err := s.scanRepo.FindByID(scanID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrScanNotFoundForSnapshot
		}
		return nil, err
	}

	return s.snapshotRepo.StreamByScanID(scanID)
}

// CountByScan returns the count of directory snapshots for a scan
func (s *DirectorySnapshotService) CountByScan(scanID int) (int64, error) {
	_, err := s.scanRepo.FindByID(scanID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, ErrScanNotFoundForSnapshot
		}
		return 0, err
	}

	return s.snapshotRepo.CountByScanID(scanID)
}

// ScanRow scans a row into DirectorySnapshot model
func (s *DirectorySnapshotService) ScanRow(rows *sql.Rows) (*model.DirectorySnapshot, error) {
	return s.snapshotRepo.ScanRow(rows)
}
