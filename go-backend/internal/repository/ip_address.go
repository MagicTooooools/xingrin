package repository

import (
	"database/sql"
	"sort"
	"time"

	"github.com/xingrin/go-backend/internal/model"
	"github.com/xingrin/go-backend/internal/pkg/scope"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// IPAddressRepository handles IP address (host_port_mapping) database operations
type IPAddressRepository struct {
	db *gorm.DB
}

// NewIPAddressRepository creates a new IP address repository
func NewIPAddressRepository(db *gorm.DB) *IPAddressRepository {
	return &IPAddressRepository{db: db}
}

// IPAddressFilterMapping defines field mapping for filtering
var IPAddressFilterMapping = scope.FilterMapping{
	"host": {Column: "host"},
	"ip":   {Column: "ip", NeedsCast: true},
	"port": {Column: "port", IsNumeric: true},
}

// IPAggregationRow represents a row from IP aggregation query
type IPAggregationRow struct {
	IP        string
	CreatedAt time.Time
}

// GetIPAggregation returns IPs with their earliest created_at, ordered by created_at DESC
func (r *IPAddressRepository) GetIPAggregation(targetID int, filter string) ([]IPAggregationRow, int64, error) {
	// Build base query
	baseQuery := r.db.Model(&model.HostPortMapping{}).Where("target_id = ?", targetID)

	// Apply filter
	baseQuery = baseQuery.Scopes(scope.WithFilter(filter, IPAddressFilterMapping))

	// Get distinct IPs with MIN(created_at)
	var results []IPAggregationRow
	err := baseQuery.
		Select("ip, MIN(created_at) as created_at").
		Group("ip").
		Order("MIN(created_at) DESC").
		Scan(&results).Error
	if err != nil {
		return nil, 0, err
	}

	return results, int64(len(results)), nil
}

// GetHostsAndPortsByIP returns hosts and ports for a specific IP
func (r *IPAddressRepository) GetHostsAndPortsByIP(targetID int, ip string, filter string) ([]string, []int, error) {
	baseQuery := r.db.Model(&model.HostPortMapping{}).
		Where("target_id = ? AND ip = ?", targetID, ip)

	// Apply filter
	baseQuery = baseQuery.Scopes(scope.WithFilter(filter, IPAddressFilterMapping))

	// Get distinct host and port combinations
	var mappings []struct {
		Host string
		Port int
	}
	err := baseQuery.
		Select("DISTINCT host, port").
		Scan(&mappings).Error
	if err != nil {
		return nil, nil, err
	}

	// Collect unique hosts and ports
	hostSet := make(map[string]struct{})
	portSet := make(map[int]struct{})
	for _, m := range mappings {
		hostSet[m.Host] = struct{}{}
		portSet[m.Port] = struct{}{}
	}

	// Convert to sorted slices
	hosts := make([]string, 0, len(hostSet))
	for h := range hostSet {
		hosts = append(hosts, h)
	}
	sort.Strings(hosts)

	ports := make([]int, 0, len(portSet))
	for p := range portSet {
		ports = append(ports, p)
	}
	sort.Ints(ports)

	return hosts, ports, nil
}

// StreamByTargetID returns a sql.Rows cursor for streaming export (raw format)
func (r *IPAddressRepository) StreamByTargetID(targetID int) (*sql.Rows, error) {
	return r.db.Model(&model.HostPortMapping{}).
		Where("target_id = ?", targetID).
		Order("ip, host, port").
		Rows()
}

// StreamByTargetIDAndIPs returns a sql.Rows cursor for streaming export filtered by IPs
func (r *IPAddressRepository) StreamByTargetIDAndIPs(targetID int, ips []string) (*sql.Rows, error) {
	return r.db.Model(&model.HostPortMapping{}).
		Where("target_id = ? AND ip IN ?", targetID, ips).
		Order("ip, host, port").
		Rows()
}

// CountByTargetID returns the count of unique IPs for a target
func (r *IPAddressRepository) CountByTargetID(targetID int) (int64, error) {
	var count int64
	err := r.db.Model(&model.HostPortMapping{}).
		Where("target_id = ?", targetID).
		Distinct("ip").
		Count(&count).Error
	return count, err
}

// ScanRow scans a single row into HostPortMapping model
func (r *IPAddressRepository) ScanRow(rows *sql.Rows) (*model.HostPortMapping, error) {
	var mapping model.HostPortMapping
	if err := r.db.ScanRows(rows, &mapping); err != nil {
		return nil, err
	}
	return &mapping, nil
}

// BulkUpsert creates multiple mappings, ignoring duplicates (ON CONFLICT DO NOTHING)
func (r *IPAddressRepository) BulkUpsert(mappings []model.HostPortMapping) (int64, error) {
	if len(mappings) == 0 {
		return 0, nil
	}

	var totalAffected int64

	// Process in batches to avoid PostgreSQL parameter limits
	batchSize := 100
	for i := 0; i < len(mappings); i += batchSize {
		end := min(i+batchSize, len(mappings))
		batch := mappings[i:end]

		// Use ON CONFLICT DO NOTHING since all fields are in unique constraint
		result := r.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&batch)
		if result.Error != nil {
			return totalAffected, result.Error
		}
		totalAffected += result.RowsAffected
	}

	return totalAffected, nil
}

// DeleteByIPs deletes all mappings for the given IPs
func (r *IPAddressRepository) DeleteByIPs(ips []string) (int64, error) {
	if len(ips) == 0 {
		return 0, nil
	}
	result := r.db.Where("ip IN ?", ips).Delete(&model.HostPortMapping{})
	return result.RowsAffected, result.Error
}
