package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xingrin/go-backend/internal/dto"
	"github.com/xingrin/go-backend/internal/model"
	"github.com/xingrin/go-backend/internal/pkg/csv"
	"github.com/xingrin/go-backend/internal/service"
)

// EndpointSnapshotHandler handles endpoint snapshot endpoints
type EndpointSnapshotHandler struct {
	svc *service.EndpointSnapshotService
}

// NewEndpointSnapshotHandler creates a new endpoint snapshot handler
func NewEndpointSnapshotHandler(svc *service.EndpointSnapshotService) *EndpointSnapshotHandler {
	return &EndpointSnapshotHandler{svc: svc}
}

// BulkUpsert creates endpoint snapshots and syncs to asset table
// POST /api/scans/:id/endpoints/bulk-upsert
func (h *EndpointSnapshotHandler) BulkUpsert(c *gin.Context) {
	scanID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		dto.BadRequest(c, "Invalid scan ID")
		return
	}

	var req dto.BulkUpsertEndpointSnapshotsRequest
	if !dto.BindJSON(c, &req) {
		return
	}

	snapshotCount, assetCount, err := h.svc.SaveAndSync(scanID, req.TargetID, req.Endpoints)
	if err != nil {
		if errors.Is(err, service.ErrScanNotFoundForSnapshot) {
			dto.NotFound(c, "Scan not found")
			return
		}
		dto.InternalError(c, "Failed to save endpoint snapshots")
		return
	}

	dto.Success(c, dto.BulkUpsertEndpointSnapshotsResponse{
		SnapshotCount: int(snapshotCount),
		AssetCount:    int(assetCount),
	})
}

// List returns paginated endpoint snapshots for a scan
// GET /api/scans/:id/endpoints
func (h *EndpointSnapshotHandler) List(c *gin.Context) {
	scanID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		dto.BadRequest(c, "Invalid scan ID")
		return
	}

	var query dto.EndpointSnapshotListQuery
	if !dto.BindQuery(c, &query) {
		return
	}

	snapshots, total, err := h.svc.ListByScan(scanID, &query)
	if err != nil {
		if errors.Is(err, service.ErrScanNotFoundForSnapshot) {
			dto.NotFound(c, "Scan not found")
			return
		}
		dto.InternalError(c, "Failed to list endpoint snapshots")
		return
	}

	var resp []dto.EndpointSnapshotResponse
	for _, s := range snapshots {
		resp = append(resp, toEndpointSnapshotResponse(&s))
	}

	dto.Paginated(c, resp, total, query.GetPage(), query.GetPageSize())
}

// Export exports endpoint snapshots as CSV
// GET /api/scans/:id/endpoints/export
func (h *EndpointSnapshotHandler) Export(c *gin.Context) {
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
		dto.InternalError(c, "Failed to export endpoint snapshots")
		return
	}

	rows, err := h.svc.StreamByScan(scanID)
	if err != nil {
		dto.InternalError(c, "Failed to export endpoint snapshots")
		return
	}

	headers := []string{
		"id", "scan_id", "url", "host", "title", "status_code",
		"content_length", "content_type", "webserver", "tech",
		"matched_gf_patterns", "created_at",
	}
	filename := fmt.Sprintf("scan-%d-endpoints.csv", scanID)

	mapper := func(rows *sql.Rows) ([]string, error) {
		snapshot, err := h.svc.ScanRow(rows)
		if err != nil {
			return nil, err
		}

		statusCode := ""
		if snapshot.StatusCode != nil {
			statusCode = strconv.Itoa(*snapshot.StatusCode)
		}

		contentLength := ""
		if snapshot.ContentLength != nil {
			contentLength = strconv.Itoa(*snapshot.ContentLength)
		}

		tech := ""
		if len(snapshot.Tech) > 0 {
			tech = strings.Join(snapshot.Tech, "|")
		}

		patterns := ""
		if len(snapshot.MatchedGFPatterns) > 0 {
			patterns = strings.Join(snapshot.MatchedGFPatterns, "|")
		}

		return []string{
			strconv.Itoa(snapshot.ID),
			strconv.Itoa(snapshot.ScanID),
			snapshot.URL,
			snapshot.Host,
			snapshot.Title,
			statusCode,
			contentLength,
			snapshot.ContentType,
			snapshot.Webserver,
			tech,
			patterns,
			snapshot.CreatedAt.Format("2006-01-02 15:04:05"),
		}, nil
	}

	if err := csv.StreamCSV(c, rows, headers, filename, mapper, count); err != nil {
		return
	}
}

func toEndpointSnapshotResponse(s *model.EndpointSnapshot) dto.EndpointSnapshotResponse {
	tech := []string(s.Tech)
	if tech == nil {
		tech = []string{}
	}
	patterns := []string(s.MatchedGFPatterns)
	if patterns == nil {
		patterns = []string{}
	}
	return dto.EndpointSnapshotResponse{
		ID:                s.ID,
		ScanID:            s.ScanID,
		URL:               s.URL,
		Host:              s.Host,
		Title:             s.Title,
		StatusCode:        s.StatusCode,
		ContentLength:     s.ContentLength,
		Location:          s.Location,
		Webserver:         s.Webserver,
		ContentType:       s.ContentType,
		Tech:              tech,
		ResponseBody:      s.ResponseBody,
		Vhost:             s.Vhost,
		MatchedGFPatterns: patterns,
		ResponseHeaders:   s.ResponseHeaders,
		CreatedAt:         s.CreatedAt,
	}
}
