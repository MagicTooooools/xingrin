package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xingrin/go-backend/internal/dto"
	"github.com/xingrin/go-backend/internal/pkg/csv"
	"github.com/xingrin/go-backend/internal/service"
)

// IPAddressHandler handles IP address endpoints
type IPAddressHandler struct {
	svc *service.IPAddressService
}

// NewIPAddressHandler creates a new IP address handler
func NewIPAddressHandler(svc *service.IPAddressService) *IPAddressHandler {
	return &IPAddressHandler{svc: svc}
}

// List returns paginated IP addresses aggregated by IP
// GET /api/targets/:id/ip-addresses
func (h *IPAddressHandler) List(c *gin.Context) {
	targetID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		dto.BadRequest(c, "Invalid target ID")
		return
	}

	var query dto.IPAddressListQuery
	if !dto.BindQuery(c, &query) {
		return
	}

	results, total, err := h.svc.ListByTarget(targetID, &query)
	if err != nil {
		if errors.Is(err, service.ErrTargetNotFound) {
			dto.NotFound(c, "Target not found")
			return
		}
		dto.InternalError(c, "Failed to list IP addresses")
		return
	}

	// Ensure empty arrays instead of null
	for i := range results {
		if results[i].Hosts == nil {
			results[i].Hosts = []string{}
		}
		if results[i].Ports == nil {
			results[i].Ports = []int{}
		}
	}

	dto.Paginated(c, results, total, query.GetPage(), query.GetPageSize())
}

// Export exports IP addresses as CSV (raw format)
// GET /api/targets/:id/ip-addresses/export
// Query params: ips (optional, comma-separated IP list for filtering)
func (h *IPAddressHandler) Export(c *gin.Context) {
	targetID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		dto.BadRequest(c, "Invalid target ID")
		return
	}

	// Parse optional IP filter
	var ips []string
	if ipsParam := c.Query("ips"); ipsParam != "" {
		ips = strings.Split(ipsParam, ",")
	}

	var rows *sql.Rows
	var count int64

	if len(ips) > 0 {
		// Export selected IPs only
		rows, err = h.svc.StreamByTargetAndIPs(targetID, ips)
		count = 0 // Unknown count for filtered export
	} else {
		// Export all
		count, err = h.svc.CountByTarget(targetID)
		if err != nil {
			if errors.Is(err, service.ErrTargetNotFound) {
				dto.NotFound(c, "Target not found")
				return
			}
			dto.InternalError(c, "Failed to export IP addresses")
			return
		}
		rows, err = h.svc.StreamByTarget(targetID)
	}

	if err != nil {
		if errors.Is(err, service.ErrTargetNotFound) {
			dto.NotFound(c, "Target not found")
			return
		}
		dto.InternalError(c, "Failed to export IP addresses")
		return
	}

	headers := []string{"ip", "host", "port", "created_at"}
	filename := fmt.Sprintf("target-%d-ip-addresses.csv", targetID)

	mapper := func(rows *sql.Rows) ([]string, error) {
		mapping, err := h.svc.ScanRow(rows)
		if err != nil {
			return nil, err
		}

		return []string{
			mapping.IP,
			mapping.Host,
			strconv.Itoa(mapping.Port),
			mapping.CreatedAt.Format("2006-01-02 15:04:05"),
		}, nil
	}

	if err := csv.StreamCSV(c, rows, headers, filename, mapper, count); err != nil {
		return
	}
}

// BulkUpsert creates multiple IP address mappings (ignores duplicates)
// POST /api/targets/:id/ip-addresses/bulk-upsert
func (h *IPAddressHandler) BulkUpsert(c *gin.Context) {
	targetID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		dto.BadRequest(c, "Invalid target ID")
		return
	}

	var req dto.BulkUpsertIPAddressesRequest
	if !dto.BindJSON(c, &req) {
		return
	}

	upsertedCount, err := h.svc.BulkUpsert(targetID, req.Mappings)
	if err != nil {
		if errors.Is(err, service.ErrTargetNotFound) {
			dto.NotFound(c, "Target not found")
			return
		}
		dto.InternalError(c, "Failed to upsert IP addresses")
		return
	}

	dto.Success(c, dto.BulkUpsertIPAddressesResponse{
		UpsertedCount: int(upsertedCount),
	})
}

// BulkDelete deletes IP address mappings by IP list
// POST /api/ip-addresses/bulk-delete
func (h *IPAddressHandler) BulkDelete(c *gin.Context) {
	var req dto.BulkDeleteIPAddressesRequest
	if !dto.BindJSON(c, &req) {
		return
	}

	deletedCount, err := h.svc.BulkDeleteByIPs(req.IPs)
	if err != nil {
		dto.InternalError(c, "Failed to delete IP addresses")
		return
	}

	dto.Success(c, dto.BulkDeleteIPAddressesResponse{
		DeletedCount: deletedCount,
	})
}
