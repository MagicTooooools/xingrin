package service

import (
	"database/sql"
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/orbit/server/internal/dto"
	"github.com/orbit/server/internal/model"
	"github.com/orbit/server/internal/pkg/validator"
	"github.com/orbit/server/internal/repository"
	"gorm.io/gorm"
)

// HostPortSnapshotService handles host-port snapshot business logic
type HostPortSnapshotService struct {
	snapshotRepo    *repository.HostPortSnapshotRepository
	scanRepo        *repository.ScanRepository
	hostPortService *HostPortService
}

// NewHostPortSnapshotService creates a new host-port snapshot service
func NewHostPortSnapshotService(
	snapshotRepo *repository.HostPortSnapshotRepository,
	scanRepo *repository.ScanRepository,
	hostPortService *HostPortService,
) *HostPortSnapshotService {
	return &HostPortSnapshotService{
		snapshotRepo:    snapshotRepo,
		scanRepo:        scanRepo,
		hostPortService: hostPortService,
	}
}

// SaveAndSync saves host-port snapshots and syncs to asset table
// 1. Validates scan exists and is not soft-deleted
// 2. Validates host/ip match target (filters invalid items)
// 3. Saves to host_port_mapping_snapshot table
// 4. TODO: Sync to host_port_mapping table (when asset service is implemented)
func (s *HostPortSnapshotService) SaveAndSync(scanID int, targetID int, items []dto.HostPortSnapshotItem) (snapshotCount int64, assetCount int64, err error) {
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

	// Filter valid host-port mappings
	snapshots := make([]model.HostPortSnapshot, 0, len(items))

	for _, item := range items {
		if isHostPortMatchTarget(item.Host, item.IP, target.Name, target.Type) {
			snapshots = append(snapshots, model.HostPortSnapshot{
				ScanID: scanID,
				Host:   item.Host,
				IP:     item.IP,
				Port:   item.Port,
			})
		}
	}

	if len(snapshots) == 0 {
		return 0, 0, nil
	}

	// Save to snapshot table
	snapshotCount, err = s.snapshotRepo.BulkCreate(snapshots)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to bulk create snapshots: %w", err)
	}

	// Sync to asset table (HostPortSnapshotItem is an alias of HostPortItem, no conversion needed)
	assetCount, err = s.hostPortService.BulkUpsert(targetID, items)
	if err != nil {
		return snapshotCount, 0, fmt.Errorf("failed to sync to asset table: %w", err)
	}

	return snapshotCount, assetCount, nil
}

// ListByScan returns paginated host-port snapshots for a scan
func (s *HostPortSnapshotService) ListByScan(scanID int, query *dto.HostPortSnapshotListQuery) ([]model.HostPortSnapshot, int64, error) {
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
func (s *HostPortSnapshotService) StreamByScan(scanID int) (*sql.Rows, error) {
	_, err := s.scanRepo.FindByID(scanID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrScanNotFoundForSnapshot
		}
		return nil, err
	}

	return s.snapshotRepo.StreamByScanID(scanID)
}

// CountByScan returns the count of host-port snapshots for a scan
func (s *HostPortSnapshotService) CountByScan(scanID int) (int64, error) {
	_, err := s.scanRepo.FindByID(scanID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, ErrScanNotFoundForSnapshot
		}
		return 0, err
	}

	return s.snapshotRepo.CountByScanID(scanID)
}

// ScanRow scans a row into HostPortSnapshot model
func (s *HostPortSnapshotService) ScanRow(rows *sql.Rows) (*model.HostPortSnapshot, error) {
	return s.snapshotRepo.ScanRow(rows)
}

// isHostPortMatchTarget checks if host/ip belongs to target
// Matching rules by target type:
//   - domain: host equals target or ends with .target
//   - ip: ip must exactly equal target
//   - cidr: ip must be within the CIDR range
func isHostPortMatchTarget(host, ip, targetName, targetType string) bool {
	host = strings.ToLower(strings.TrimSpace(host))
	ip = strings.TrimSpace(ip)
	targetName = strings.ToLower(strings.TrimSpace(targetName))

	if host == "" || ip == "" || targetName == "" {
		return false
	}

	switch targetType {
	case validator.TargetTypeDomain:
		// Check if host matches target domain
		return host == targetName || strings.HasSuffix(host, "."+targetName)

	case validator.TargetTypeIP:
		// IP must exactly match target
		return ip == targetName

	case validator.TargetTypeCIDR:
		// IP must be within CIDR range
		ipAddr := net.ParseIP(ip)
		if ipAddr == nil {
			return false
		}
		_, network, err := net.ParseCIDR(targetName)
		if err != nil {
			return false
		}
		return network.Contains(ipAddr)

	default:
		return false
	}
}
