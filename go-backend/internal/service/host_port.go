package service

import (
	"database/sql"
	"errors"

	"github.com/xingrin/go-backend/internal/dto"
	"github.com/xingrin/go-backend/internal/model"
	"github.com/xingrin/go-backend/internal/repository"
	"gorm.io/gorm"
)

// HostPortService handles host-port business logic
type HostPortService struct {
	repo       *repository.HostPortRepository
	targetRepo *repository.TargetRepository
}

// NewHostPortService creates a new host-port service
func NewHostPortService(repo *repository.HostPortRepository, targetRepo *repository.TargetRepository) *HostPortService {
	return &HostPortService{repo: repo, targetRepo: targetRepo}
}

// ListByTarget returns paginated host-ports aggregated by IP
func (s *HostPortService) ListByTarget(targetID int, query *dto.HostPortListQuery) ([]dto.HostPortResponse, int64, error) {
	// Check if target exists
	_, err := s.targetRepo.FindByID(targetID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, ErrTargetNotFound
		}
		return nil, 0, err
	}

	// Get IP aggregation (all IPs with their earliest created_at)
	ipRows, total, err := s.repo.GetIPAggregation(targetID, query.Filter)
	if err != nil {
		return nil, 0, err
	}

	// Apply pagination to IP list
	page := query.GetPage()
	pageSize := query.GetPageSize()
	start := (page - 1) * pageSize
	end := start + pageSize

	if start >= len(ipRows) {
		return []dto.HostPortResponse{}, total, nil
	}
	if end > len(ipRows) {
		end = len(ipRows)
	}

	pagedIPs := ipRows[start:end]

	// For each IP, get its hosts and ports
	results := make([]dto.HostPortResponse, 0, len(pagedIPs))
	for _, row := range pagedIPs {
		hosts, ports, err := s.repo.GetHostsAndPortsByIP(targetID, row.IP, query.Filter)
		if err != nil {
			return nil, 0, err
		}

		results = append(results, dto.HostPortResponse{
			IP:        row.IP,
			Hosts:     hosts,
			Ports:     ports,
			CreatedAt: row.CreatedAt,
		})
	}

	return results, total, nil
}

// StreamByTarget returns a cursor for streaming export (raw format)
func (s *HostPortService) StreamByTarget(targetID int) (*sql.Rows, error) {
	_, err := s.targetRepo.FindByID(targetID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTargetNotFound
		}
		return nil, err
	}
	return s.repo.StreamByTargetID(targetID)
}

// StreamByTargetAndIPs returns a cursor for streaming export filtered by IPs
func (s *HostPortService) StreamByTargetAndIPs(targetID int, ips []string) (*sql.Rows, error) {
	_, err := s.targetRepo.FindByID(targetID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTargetNotFound
		}
		return nil, err
	}
	return s.repo.StreamByTargetIDAndIPs(targetID, ips)
}

// CountByTarget returns the count of unique IPs for a target
func (s *HostPortService) CountByTarget(targetID int) (int64, error) {
	_, err := s.targetRepo.FindByID(targetID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, ErrTargetNotFound
		}
		return 0, err
	}
	return s.repo.CountByTargetID(targetID)
}

// ScanRow scans a row into HostPort model
func (s *HostPortService) ScanRow(rows *sql.Rows) (*model.HostPort, error) {
	return s.repo.ScanRow(rows)
}

// BulkUpsert creates multiple mappings for a target (ignores duplicates)
func (s *HostPortService) BulkUpsert(targetID int, items []dto.HostPortItem) (int64, error) {
	_, err := s.targetRepo.FindByID(targetID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, ErrTargetNotFound
		}
		return 0, err
	}

	// Convert DTOs to models
	mappings := make([]model.HostPort, 0, len(items))
	for _, item := range items {
		mappings = append(mappings, model.HostPort{
			TargetID: targetID,
			Host:     item.Host,
			IP:       item.IP,
			Port:     item.Port,
		})
	}

	if len(mappings) == 0 {
		return 0, nil
	}

	return s.repo.BulkUpsert(mappings)
}

// BulkDeleteByIPs deletes all mappings for the given IPs
func (s *HostPortService) BulkDeleteByIPs(ips []string) (int64, error) {
	if len(ips) == 0 {
		return 0, nil
	}
	return s.repo.DeleteByIPs(ips)
}
