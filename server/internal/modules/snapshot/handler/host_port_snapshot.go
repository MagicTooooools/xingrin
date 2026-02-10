package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	service "github.com/yyhuni/lunafox/server/internal/modules/snapshot/application"
	"github.com/yyhuni/lunafox/server/internal/modules/snapshot/dto"
	"github.com/yyhuni/lunafox/server/internal/pkg/csv"
	"github.com/yyhuni/lunafox/server/internal/pkg/timeutil"
)

// HostPortSnapshotHandler handles host-port snapshot endpoints
type HostPortSnapshotHandler struct {
	svc *service.HostPortSnapshotFacade
}

// NewHostPortSnapshotHandler creates a new host-port snapshot handler
func NewHostPortSnapshotHandler(svc *service.HostPortSnapshotFacade) *HostPortSnapshotHandler {
	return &HostPortSnapshotHandler{svc: svc}
}

// BulkUpsert creates host-port snapshots and syncs to asset table
// POST /api/scans/:id/host-ports/bulk-upsert
func (h *HostPortSnapshotHandler) BulkUpsert(c *gin.Context) {
	scanID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		dto.BadRequest(c, "Invalid scan ID")
		return
	}

	var req dto.BulkUpsertHostPortSnapshotsRequest
	if !dto.BindJSON(c, &req) {
		return
	}

	snapshotCount, assetCount, err := h.svc.SaveAndSync(scanID, req.TargetID, req.HostPorts)
	if err != nil {
		if errors.Is(err, service.ErrScanNotFoundForSnapshot) {
			dto.NotFound(c, "Scan not found")
			return
		}
		if errors.Is(err, service.ErrTargetMismatch) {
			dto.BadRequest(c, "targetId does not match scan's target")
			return
		}
		dto.InternalError(c, "Failed to save snapshots")
		return
	}

	dto.Success(c, dto.BulkUpsertHostPortSnapshotsResponse{
		SnapshotCount: int(snapshotCount),
		AssetCount:    int(assetCount),
	})
}

// List returns paginated host-port snapshots for a scan
// GET /api/scans/:id/host-ports
func (h *HostPortSnapshotHandler) List(c *gin.Context) {
	scanID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		dto.BadRequest(c, "Invalid scan ID")
		return
	}

	var query dto.HostPortSnapshotListQuery
	if !dto.BindQuery(c, &query) {
		return
	}

	snapshots, total, err := h.svc.ListByScan(scanID, &query)
	if err != nil {
		if errors.Is(err, service.ErrScanNotFoundForSnapshot) {
			dto.NotFound(c, "Scan not found")
			return
		}
		dto.InternalError(c, "Failed to list snapshots")
		return
	}

	var resp []dto.HostPortSnapshotResponse
	for _, s := range snapshots {
		resp = append(resp, toHostPortSnapshotResponse(&s))
	}

	dto.Paginated(c, resp, total, query.GetPage(), query.GetPageSize())
}

// Export exports host-port snapshots as CSV
// GET /api/scans/:id/host-ports/export
func (h *HostPortSnapshotHandler) Export(c *gin.Context) {
	scanID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		dto.BadRequest(c, "Invalid scan ID")
		return
	}

	count, err := h.svc.CountByScan(scanID)
	if err != nil {
		if errors.Is(err, service.ErrScanNotFoundForSnapshot) {
			dto.NotFound(c, "Scan not found")
			return
		}
		dto.InternalError(c, "Failed to export snapshots")
		return
	}

	rows, err := h.svc.StreamByScan(scanID)
	if err != nil {
		dto.InternalError(c, "Failed to export snapshots")
		return
	}

	headers := []string{"id", "scan_id", "host", "ip", "port", "created_at"}
	filename := fmt.Sprintf("scan-%d-host-ports.csv", scanID)

	mapper := func(rows *sql.Rows) ([]string, error) {
		snapshot, err := h.svc.ScanRow(rows)
		if err != nil {
			return nil, err
		}

		return []string{
			strconv.Itoa(snapshot.ID),
			strconv.Itoa(snapshot.ScanID),
			snapshot.Host,
			snapshot.IP,
			strconv.Itoa(snapshot.Port),
			timeutil.FormatRFC3339NanoUTC(snapshot.CreatedAt),
		}, nil
	}

	if err := csv.StreamCSV(c, rows, headers, filename, mapper, count); err != nil {
		return
	}
}

// toHostPortSnapshotResponse converts model to response DTO
func toHostPortSnapshotResponse(s *service.HostPortSnapshot) dto.HostPortSnapshotResponse {
	return dto.HostPortSnapshotResponse{
		ID:        s.ID,
		ScanID:    s.ScanID,
		Host:      s.Host,
		IP:        s.IP,
		Port:      s.Port,
		CreatedAt: timeutil.ToUTC(s.CreatedAt),
	}
}
