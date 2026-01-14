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

// SubdomainSnapshotService handles subdomain snapshot business logic
type SubdomainSnapshotService struct {
	snapshotRepo     *repository.SubdomainSnapshotRepository
	scanRepo         *repository.ScanRepository
	subdomainService *SubdomainService
}

// NewSubdomainSnapshotService creates a new subdomain snapshot service
func NewSubdomainSnapshotService(
	snapshotRepo *repository.SubdomainSnapshotRepository,
	scanRepo *repository.ScanRepository,
	subdomainService *SubdomainService,
) *SubdomainSnapshotService {
	return &SubdomainSnapshotService{
		snapshotRepo:     snapshotRepo,
		scanRepo:         scanRepo,
		subdomainService: subdomainService,
	}
}

// SaveAndSync saves subdomain snapshots and syncs to asset table
// 1. Validates scan exists and is not soft-deleted
// 2. Validates target type is "domain" (only domains can have subdomains)
// 3. Validates subdomain names match target (filters invalid items)
// 4. Saves to subdomain_snapshot table
// 5. Calls SubdomainService.BulkCreate to sync to subdomain table
func (s *SubdomainSnapshotService) SaveAndSync(scanID int, targetID int, items []dto.SubdomainSnapshotItem) (snapshotCount int64, assetCount int64, err error) {
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

	// Only domain type targets can have subdomains
	if target.Type != "domain" {
		return 0, 0, ErrInvalidTargetType
	}

	// Filter valid subdomains
	snapshots := make([]model.SubdomainSnapshot, 0, len(items))
	validNames := make([]string, 0, len(items))

	for _, item := range items {
		if validator.IsSubdomainMatchTarget(item.Name, target.Name) {
			snapshots = append(snapshots, model.SubdomainSnapshot{
				ScanID: scanID,
				Name:   item.Name,
			})
			validNames = append(validNames, item.Name)
		}
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
	assetCountInt, err := s.subdomainService.BulkCreate(targetID, validNames)
	if err != nil {
		return snapshotCount, 0, nil
	}

	return snapshotCount, int64(assetCountInt), nil
}

// ListByScan returns paginated subdomain snapshots for a scan
func (s *SubdomainSnapshotService) ListByScan(scanID int, query *dto.SubdomainSnapshotListQuery) ([]model.SubdomainSnapshot, int64, error) {
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
func (s *SubdomainSnapshotService) StreamByScan(scanID int) (*sql.Rows, error) {
	_, err := s.scanRepo.FindByID(scanID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrScanNotFoundForSnapshot
		}
		return nil, err
	}

	return s.snapshotRepo.StreamByScanID(scanID)
}

// CountByScan returns the count of subdomain snapshots for a scan
func (s *SubdomainSnapshotService) CountByScan(scanID int) (int64, error) {
	_, err := s.scanRepo.FindByID(scanID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, ErrScanNotFoundForSnapshot
		}
		return 0, err
	}

	return s.snapshotRepo.CountByScanID(scanID)
}

// ScanRow scans a row into SubdomainSnapshot model
func (s *SubdomainSnapshotService) ScanRow(rows *sql.Rows) (*model.SubdomainSnapshot, error) {
	return s.snapshotRepo.ScanRow(rows)
}
