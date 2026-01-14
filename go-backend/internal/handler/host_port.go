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

// HostPortHandler handles host-port endpoints
type HostPortHandler struct {
	svc *service.HostPortService
}

// NewHostPortHandler creates a new host-port handler
func NewHostPortHandler(svc *service.HostPortService) *HostPortHandler {
	return &HostPortHandler{svc: svc}
}

// List returns paginated host-ports aggregated by IP
// GET /api/targets/:id/host-ports
func (h *HostPortHandler) List(c *gin.Context) {
	targetID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		dto.BadRequest(c, "Invalid target ID")
		return
	}

	var query dto.HostPortListQuery
	if !dto.BindQuery(c, &query) {
		return
	}

	results, total, err := h.svc.ListByTarget(targetID, &query)
	if err != nil {
		if errors.Is(err, service.ErrTargetNotFound) {
			dto.NotFound(c, "Target not found")
			return
		}
		dto.InternalError(c, "Failed to list host-ports")
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

// Export exports host-ports as CSV (raw format)
// GET /api/targets/:id/host-ports/export
// Query params: ips (optional, comma-separated IP list for filtering)
func (h *HostPortHandler) Export(c *gin.Context) {
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
			dto.InternalError(c, "Failed to export host-ports")
			return
		}
		rows, err = h.svc.StreamByTarget(targetID)
	}

	if err != nil {
		if errors.Is(err, service.ErrTargetNotFound) {
			dto.NotFound(c, "Target not found")
			return
		}
		dto.InternalError(c, "Failed to export host-ports")
		return
	}

	headers := []string{"ip", "host", "port", "created_at"}
	filename := fmt.Sprintf("target-%d-host-ports.csv", targetID)

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

// BulkUpsert creates multiple host-port mappings (ignores duplicates)
// POST /api/targets/:id/host-ports/bulk-upsert
func (h *HostPortHandler) BulkUpsert(c *gin.Context) {
	targetID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		dto.BadRequest(c, "Invalid target ID")
		return
	}

	var req dto.BulkUpsertHostPortsRequest
	if !dto.BindJSON(c, &req) {
		return
	}

	upsertedCount, err := h.svc.BulkUpsert(targetID, req.Mappings)
	if err != nil {
		if errors.Is(err, service.ErrTargetNotFound) {
			dto.NotFound(c, "Target not found")
			return
		}
		dto.InternalError(c, "Failed to upsert host-ports")
		return
	}

	dto.Success(c, dto.BulkUpsertHostPortsResponse{
		UpsertedCount: int(upsertedCount),
	})
}

// BulkDelete deletes host-port mappings by IP list
// POST /api/host-ports/bulk-delete
func (h *HostPortHandler) BulkDelete(c *gin.Context) {
	var req dto.BulkDeleteHostPortsRequest
	if !dto.BindJSON(c, &req) {
		return
	}

	deletedCount, err := h.svc.BulkDeleteByIPs(req.IPs)
	if err != nil {
		dto.InternalError(c, "Failed to delete host-ports")
		return
	}

	dto.Success(c, dto.BulkDeleteHostPortsResponse{
		DeletedCount: deletedCount,
	})
}
