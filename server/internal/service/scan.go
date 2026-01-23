package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"slices"
	"strings"

	"github.com/yyhuni/orbit/server/internal/dto"
	"github.com/yyhuni/orbit/server/internal/engineschema"
	"github.com/yyhuni/orbit/server/internal/model"
	"github.com/yyhuni/orbit/server/internal/repository"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
)

var (
	ErrScanNotFound           = errors.New("scan not found")
	ErrScanCannotStop         = errors.New("scan cannot be stopped in current status")
	ErrNoTargetsForScan       = errors.New("no targets provided for scan")
	ErrScanInvalidConfig      = errors.New("invalid scan configuration")
	ErrScanInvalidEngineNames = errors.New("invalid engineNames")
)

// ScanService handles scan business logic
type ScanService struct {
	repo       *repository.ScanRepository
	targetRepo *repository.TargetRepository
	orgRepo    *repository.OrganizationRepository
}

// NewScanService creates a new scan service
func NewScanService(
	repo *repository.ScanRepository,
	scanLogRepo *repository.ScanLogRepository, // Keep for backward compatibility, but not used
	targetRepo *repository.TargetRepository,
	orgRepo *repository.OrganizationRepository,
) *ScanService {
	return &ScanService{
		repo:       repo,
		targetRepo: targetRepo,
		orgRepo:    orgRepo,
	}
}

// List returns paginated scans
func (s *ScanService) List(query *dto.ScanListQuery) ([]model.Scan, int64, error) {
	return s.repo.FindAll(query.GetPage(), query.GetPageSize(), query.TargetID, query.Status, query.Search)
}

// GetByID returns a scan by ID
func (s *ScanService) GetByID(id int) (*model.Scan, error) {
	scan, err := s.repo.FindByIDWithTarget(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrScanNotFound
		}
		return nil, err
	}
	return scan, nil
}

// Delete soft deletes a scan (two-phase delete)
func (s *ScanService) Delete(id int) (int64, []string, error) {
	// Check if scan exists
	_, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil, ErrScanNotFound
		}
		return 0, nil, err
	}

	// Soft delete
	return s.repo.BulkSoftDelete([]int{id})
}

// BulkDelete soft deletes multiple scans
func (s *ScanService) BulkDelete(ids []int) (int64, []string, error) {
	if len(ids) == 0 {
		return 0, nil, nil
	}
	return s.repo.BulkSoftDelete(ids)
}

// GetStatistics returns scan statistics
func (s *ScanService) GetStatistics() (*dto.ScanStatisticsResponse, error) {
	stats, err := s.repo.GetStatistics()
	if err != nil {
		return nil, err
	}

	return &dto.ScanStatisticsResponse{
		Total:           stats.Total,
		Running:         stats.Running,
		Completed:       stats.Completed,
		Failed:          stats.Failed,
		TotalVulns:      stats.TotalVulns,
		TotalSubdomains: stats.TotalSubdomains,
		TotalEndpoints:  stats.TotalEndpoints,
		TotalWebsites:   stats.TotalWebsites,
		TotalAssets:     stats.TotalAssets,
	}, nil
}

// Stop stops a running scan
func (s *ScanService) Stop(id int) (int, error) {
	scan, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, ErrScanNotFound
		}
		return 0, err
	}

	// Check if scan can be stopped (only pending or running can be stopped)
	if scan.Status != model.ScanStatusRunning && scan.Status != model.ScanStatusPending {
		return 0, ErrScanCannotStop
	}

	// Update status to cancelled
	if err := s.repo.UpdateStatus(id, model.ScanStatusCancelled); err != nil {
		return 0, err
	}

	// TODO: Revoke celery tasks when worker integration is implemented
	// For now, just return 0 revoked tasks
	return 0, nil
}

func normalizeEngineNames(engineNames []string) ([]string, error) {
	if len(engineNames) == 0 {
		return nil, nil
	}

	cleaned := make([]string, 0, len(engineNames))
	for _, name := range engineNames {
		n := strings.TrimSpace(name)
		if n == "" {
			continue
		}
		cleaned = append(cleaned, n)
	}

	// Deduplicate while preserving order
	seen := map[string]bool{}
	unique := make([]string, 0, len(cleaned))
	for _, name := range cleaned {
		if seen[name] {
			continue
		}
		seen[name] = true
		unique = append(unique, name)
	}

	if slices.Contains(unique, "") {
		return nil, ErrScanInvalidEngineNames
	}

	return unique, nil
}

func parseYAMLMapping(b []byte) (map[string]any, error) {
	var root map[string]any
	if err := yaml.Unmarshal(b, &root); err != nil {
		return nil, err
	}
	return root, nil
}

// CreateNormal creates a scan record for an existing target (create-only, no scheduling).
// It validates the configuration YAML, and validates known engines (those with embedded schemas).
func (s *ScanService) CreateNormal(req *dto.CreateScanRequest) (*model.Scan, error) {
	if req == nil {
		return nil, ErrScanInvalidConfig
	}
	if req.TargetID == 0 {
		return nil, ErrTargetNotFound
	}

	configYAML := strings.TrimSpace(req.Configuration)
	if configYAML == "" {
		return nil, ErrScanInvalidConfig
	}

	root, err := parseYAMLMapping([]byte(configYAML))
	if err != nil {
		return nil, fmt.Errorf("%w: parse yaml: %v", ErrScanInvalidConfig, err)
	}
	if root == nil {
		return nil, fmt.Errorf("%w: yaml must be a mapping", ErrScanInvalidConfig)
	}

	if len(req.EngineNames) != 1 {
		return nil, ErrScanInvalidEngineNames
	}

	engineNames, err := normalizeEngineNames(req.EngineNames)
	if err != nil {
		return nil, err
	}
	if len(engineNames) != 1 {
		return nil, ErrScanInvalidEngineNames
	}

	// Validate known engines.
	candidates := engineNames
	for _, engine := range candidates {
		engine = strings.TrimSpace(engine)
		if engine == "" {
			continue
		}
		if err := engineschema.ValidateYAML(engine, []byte(configYAML)); err != nil {
			// If schema doesn't exist yet for this engine, skip validation.
			if errors.Is(err, fs.ErrNotExist) {
				continue
			}
			return nil, fmt.Errorf("%w: %s: %v", ErrScanInvalidConfig, engine, err)
		}
	}

	// Ensure target exists.
	target, err := s.targetRepo.FindByID(req.TargetID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTargetNotFound
		}
		return nil, err
	}

	// Prepare scan payload.
	engineNamesJSON, err := json.Marshal(engineNames)
	if err != nil {
		return nil, err
	}

	engineIDs := make([]int64, 0, len(req.EngineIDs))
	for _, id := range req.EngineIDs {
		engineIDs = append(engineIDs, int64(id))
	}

	scan := &model.Scan{
		TargetID:          req.TargetID,
		EngineIDs:         engineIDs,
		EngineNames:       engineNamesJSON,
		YamlConfiguration: configYAML,
		ScanMode:          model.ScanModeFull,
		Status:            model.ScanStatusPending,
	}

	inputs := []model.ScanInputTarget{{
		Value:     target.Name,
		InputType: target.Type,
	}}
	if err := s.repo.CreateWithInputTargets(scan, inputs); err != nil {
		return nil, err
	}

	// Attach target for response shaping.
	scan.Target = target

	return scan, nil
}

// ToScanResponse converts scan model to response DTO
func (s *ScanService) ToScanResponse(scan *model.Scan) *dto.ScanResponse {
	resp := &dto.ScanResponse{
		ID:           scan.ID,
		TargetID:     scan.TargetID,
		ScanMode:     scan.ScanMode,
		Status:       scan.Status,
		Progress:     scan.Progress,
		CurrentStage: scan.CurrentStage,
		ErrorMessage: scan.ErrorMessage,
		CreatedAt:    scan.CreatedAt,
		StoppedAt:    scan.StoppedAt,
	}

	// Convert engine IDs
	if scan.EngineIDs != nil {
		resp.EngineIDs = scan.EngineIDs
	} else {
		resp.EngineIDs = []int64{}
	}

	// Convert engine names from JSON
	if scan.EngineNames != nil {
		var names []string
		if err := json.Unmarshal(scan.EngineNames, &names); err == nil {
			resp.EngineNames = names
		} else {
			resp.EngineNames = []string{}
		}
	} else {
		resp.EngineNames = []string{}
	}

	// Add target info
	if scan.Target != nil {
		resp.Target = &dto.TargetBrief{
			ID:   scan.Target.ID,
			Name: scan.Target.Name,
			Type: scan.Target.Type,
		}
	}

	// Add cached stats
	resp.CachedStats = &dto.ScanCachedStats{
		SubdomainsCount:  scan.CachedSubdomainsCount,
		WebsitesCount:    scan.CachedWebsitesCount,
		EndpointsCount:   scan.CachedEndpointsCount,
		IPsCount:         scan.CachedIPsCount,
		DirectoriesCount: scan.CachedDirectoriesCount,
		ScreenshotsCount: scan.CachedScreenshotsCount,
		VulnsTotal:       scan.CachedVulnsTotal,
		VulnsCritical:    scan.CachedVulnsCritical,
		VulnsHigh:        scan.CachedVulnsHigh,
		VulnsMedium:      scan.CachedVulnsMedium,
		VulnsLow:         scan.CachedVulnsLow,
	}

	return resp
}

// ToScanDetailResponse converts scan model to detailed response DTO
func (s *ScanService) ToScanDetailResponse(scan *model.Scan) *dto.ScanDetailResponse {
	base := s.ToScanResponse(scan)

	resp := &dto.ScanDetailResponse{
		ScanResponse:      *base,
		YamlConfiguration: scan.YamlConfiguration,
		ResultsDir:        scan.ResultsDir,
		WorkerID:          scan.WorkerID,
	}

	// Convert stage progress from JSON
	if scan.StageProgress != nil {
		var progress map[string]interface{}
		if err := json.Unmarshal(scan.StageProgress, &progress); err == nil {
			resp.StageProgress = progress
		}
	}

	return resp
}
