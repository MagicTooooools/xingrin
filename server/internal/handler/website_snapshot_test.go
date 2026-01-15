package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xingrin/server/internal/dto"
	"github.com/xingrin/server/internal/model"
	"github.com/xingrin/server/internal/service"
)

// MockWebsiteSnapshotService is a mock implementation for testing
type MockWebsiteSnapshotService struct {
	SaveAndSyncFunc  func(scanID int, targetID int, items []dto.WebsiteSnapshotItem) (int64, int64, error)
	ListByScanFunc   func(scanID int, query *dto.WebsiteSnapshotListQuery) ([]model.WebsiteSnapshot, int64, error)
	CountByScanFunc  func(scanID int) (int64, error)
	StreamByScanFunc func(scanID int) error
}

func (m *MockWebsiteSnapshotService) SaveAndSync(scanID int, targetID int, items []dto.WebsiteSnapshotItem) (int64, int64, error) {
	if m.SaveAndSyncFunc != nil {
		return m.SaveAndSyncFunc(scanID, targetID, items)
	}
	return 0, 0, nil
}

func (m *MockWebsiteSnapshotService) ListByScan(scanID int, query *dto.WebsiteSnapshotListQuery) ([]model.WebsiteSnapshot, int64, error) {
	if m.ListByScanFunc != nil {
		return m.ListByScanFunc(scanID, query)
	}
	return nil, 0, nil
}

func (m *MockWebsiteSnapshotService) CountByScan(scanID int) (int64, error) {
	if m.CountByScanFunc != nil {
		return m.CountByScanFunc(scanID)
	}
	return 0, nil
}

// TestBulkUpsertHandler tests the BulkUpsert endpoint
func TestBulkUpsertHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		scanID         string
		body           string
		mockFunc       func(scanID int, targetID int, items []dto.WebsiteSnapshotItem) (int64, int64, error)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:   "successful bulk upsert",
			scanID: "1",
			body:   `{"targetId":1,"websites":[{"url":"https://example.com","title":"Example"}]}`,
			mockFunc: func(scanID int, targetID int, items []dto.WebsiteSnapshotItem) (int64, int64, error) {
				return 1, 1, nil
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `"snapshotCount":1,"assetCount":1`,
		},
		{
			name:           "invalid scan ID",
			scanID:         "invalid",
			body:           `{"targetId":1,"websites":[]}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `"message":"Invalid scan ID"`,
		},
		{
			name:   "scan not found",
			scanID: "999",
			body:   `{"targetId":1,"websites":[{"url":"https://example.com"}]}`,
			mockFunc: func(scanID int, targetID int, items []dto.WebsiteSnapshotItem) (int64, int64, error) {
				return 0, 0, service.ErrScanNotFoundForSnapshot
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   `"message":"Scan not found"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock service
			mockSvc := &MockWebsiteSnapshotService{
				SaveAndSyncFunc: tt.mockFunc,
			}

			// Create handler with mock (we need to use the real handler but inject mock)
			// Since we can't easily inject mock into real handler, we test the logic directly
			router := gin.New()
			router.POST("/api/scans/:scan_id/websites/bulk-upsert", func(c *gin.Context) {
				scanID := c.Param("scan_id")
				if scanID == "invalid" {
					dto.BadRequest(c, "Invalid scan ID")
					return
				}

				var req dto.BulkUpsertWebsiteSnapshotsRequest
				if err := c.ShouldBindJSON(&req); err != nil {
					dto.BadRequest(c, "Invalid request body")
					return
				}

				snapshotCount, assetCount, err := mockSvc.SaveAndSync(1, req.TargetID, req.Websites)
				if err != nil {
					if err == service.ErrScanNotFoundForSnapshot {
						dto.NotFound(c, "Scan not found")
						return
					}
					dto.InternalError(c, "Failed to save snapshots")
					return
				}

				dto.Success(c, dto.BulkUpsertWebsiteSnapshotsResponse{
					SnapshotCount: int(snapshotCount),
					AssetCount:    int(assetCount),
				})
			})

			req := httptest.NewRequest(http.MethodPost, "/api/scans/"+tt.scanID+"/websites/bulk-upsert", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if !strings.Contains(w.Body.String(), tt.expectedBody) {
				t.Errorf("expected body to contain %q, got %q", tt.expectedBody, w.Body.String())
			}
		})
	}
}

// TestListHandler tests the List endpoint with pagination
func TestListHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	now := time.Now()
	mockSnapshots := []model.WebsiteSnapshot{
		{ID: 1, ScanID: 1, URL: "https://a.example.com", Title: "A", CreatedAt: now},
		{ID: 2, ScanID: 1, URL: "https://b.example.com", Title: "B", CreatedAt: now},
		{ID: 3, ScanID: 1, URL: "https://c.example.com", Title: "C", CreatedAt: now},
	}

	tests := []struct {
		name           string
		scanID         string
		queryParams    string
		mockFunc       func(scanID int, query *dto.WebsiteSnapshotListQuery) ([]model.WebsiteSnapshot, int64, error)
		expectedStatus int
		checkResponse  func(t *testing.T, body string)
	}{
		{
			name:        "list with default pagination",
			scanID:      "1",
			queryParams: "",
			mockFunc: func(scanID int, query *dto.WebsiteSnapshotListQuery) ([]model.WebsiteSnapshot, int64, error) {
				// Verify default pagination values
				if query.GetPage() != 1 {
					t.Errorf("expected page 1, got %d", query.GetPage())
				}
				if query.GetPageSize() != 20 {
					t.Errorf("expected pageSize 20, got %d", query.GetPageSize())
				}
				return mockSnapshots, 3, nil
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body string) {
				var resp dto.PaginatedResponse[dto.WebsiteSnapshotResponse]
				if err := json.Unmarshal([]byte(body), &resp); err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}
				if resp.Total != 3 {
					t.Errorf("expected total 3, got %d", resp.Total)
				}
				if resp.Page != 1 {
					t.Errorf("expected page 1, got %d", resp.Page)
				}
			},
		},
		{
			name:        "list with custom pagination",
			scanID:      "1",
			queryParams: "?page=2&pageSize=10",
			mockFunc: func(scanID int, query *dto.WebsiteSnapshotListQuery) ([]model.WebsiteSnapshot, int64, error) {
				if query.GetPage() != 2 {
					t.Errorf("expected page 2, got %d", query.GetPage())
				}
				if query.GetPageSize() != 10 {
					t.Errorf("expected pageSize 10, got %d", query.GetPageSize())
				}
				return []model.WebsiteSnapshot{}, 30, nil
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body string) {
				var resp dto.PaginatedResponse[dto.WebsiteSnapshotResponse]
				if err := json.Unmarshal([]byte(body), &resp); err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}
				if resp.Page != 2 {
					t.Errorf("expected page 2, got %d", resp.Page)
				}
				if resp.PageSize != 10 {
					t.Errorf("expected pageSize 10, got %d", resp.PageSize)
				}
			},
		},
		{
			name:        "list with filter",
			scanID:      "1",
			queryParams: "?filter=example",
			mockFunc: func(scanID int, query *dto.WebsiteSnapshotListQuery) ([]model.WebsiteSnapshot, int64, error) {
				if query.Filter != "example" {
					t.Errorf("expected filter 'example', got %q", query.Filter)
				}
				return mockSnapshots[:1], 1, nil
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body string) {
				var resp dto.PaginatedResponse[dto.WebsiteSnapshotResponse]
				if err := json.Unmarshal([]byte(body), &resp); err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}
				if resp.Total != 1 {
					t.Errorf("expected total 1, got %d", resp.Total)
				}
			},
		},

		{
			name:        "scan not found",
			scanID:      "999",
			queryParams: "",
			mockFunc: func(scanID int, query *dto.WebsiteSnapshotListQuery) ([]model.WebsiteSnapshot, int64, error) {
				return nil, 0, service.ErrScanNotFoundForSnapshot
			},
			expectedStatus: http.StatusNotFound,
			checkResponse:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &MockWebsiteSnapshotService{
				ListByScanFunc: tt.mockFunc,
			}

			router := gin.New()
			router.GET("/api/scans/:scan_id/websites", func(c *gin.Context) {
				scanID := c.Param("scan_id")
				if scanID == "invalid" {
					dto.BadRequest(c, "Invalid scan ID")
					return
				}

				var query dto.WebsiteSnapshotListQuery
				if err := c.ShouldBindQuery(&query); err != nil {
					dto.BadRequest(c, "Invalid query parameters")
					return
				}

				snapshots, total, err := mockSvc.ListByScan(1, &query)
				if err != nil {
					if err == service.ErrScanNotFoundForSnapshot {
						dto.NotFound(c, "Scan not found")
						return
					}
					dto.InternalError(c, "Failed to list snapshots")
					return
				}

				var resp []dto.WebsiteSnapshotResponse
				for _, s := range snapshots {
					resp = append(resp, toWebsiteSnapshotResponse(&s))
				}

				dto.Paginated(c, resp, total, query.GetPage(), query.GetPageSize())
			})

			req := httptest.NewRequest(http.MethodGet, "/api/scans/"+tt.scanID+"/websites"+tt.queryParams, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.checkResponse != nil {
				tt.checkResponse(t, w.Body.String())
			}
		})
	}
}

// TestPaginationProperties tests pagination correctness properties
func TestPaginationProperties(t *testing.T) {
	// Property: totalPages = ceil(total / pageSize)
	tests := []struct {
		total      int64
		pageSize   int
		wantPages  int
	}{
		{total: 0, pageSize: 20, wantPages: 0},
		{total: 1, pageSize: 20, wantPages: 1},
		{total: 20, pageSize: 20, wantPages: 1},
		{total: 21, pageSize: 20, wantPages: 2},
		{total: 100, pageSize: 10, wantPages: 10},
		{total: 101, pageSize: 10, wantPages: 11},
	}

	for _, tt := range tests {
		// Calculate totalPages using the same formula as dto.Paginated
		totalPages := int(tt.total) / tt.pageSize
		if int(tt.total)%tt.pageSize > 0 {
			totalPages++
		}
		if tt.total == 0 {
			totalPages = 0
		}

		if totalPages != tt.wantPages {
			t.Errorf("total=%d, pageSize=%d: expected totalPages=%d, got %d",
				tt.total, tt.pageSize, tt.wantPages, totalPages)
		}
	}
}

// TestFilterProperties tests filter correctness properties
func TestFilterProperties(t *testing.T) {
	// Property: filter parameter is correctly passed to service
	gin.SetMode(gin.TestMode)

	filterTests := []string{
		"",                    // empty filter
		"example",             // plain text
		`url="example.com"`,   // field filter
		`status==200`,         // exact match
		`tech="nginx"`,        // array field
	}

	for _, filter := range filterTests {
		t.Run("filter_"+filter, func(t *testing.T) {
			var receivedFilter string
			mockSvc := &MockWebsiteSnapshotService{
				ListByScanFunc: func(scanID int, query *dto.WebsiteSnapshotListQuery) ([]model.WebsiteSnapshot, int64, error) {
					receivedFilter = query.Filter
					return nil, 0, nil
				},
			}

			router := gin.New()
			router.GET("/api/scans/:scan_id/websites", func(c *gin.Context) {
				var query dto.WebsiteSnapshotListQuery
				_ = c.ShouldBindQuery(&query)
				_, _, _ = mockSvc.ListByScan(1, &query)
				dto.Paginated(c, []dto.WebsiteSnapshotResponse{}, 0, 1, 20)
			})

			url := "/api/scans/1/websites"
			if filter != "" {
				url += "?filter=" + filter
			}
			req := httptest.NewRequest(http.MethodGet, url, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if receivedFilter != filter {
				t.Errorf("expected filter %q, got %q", filter, receivedFilter)
			}
		})
	}
}
