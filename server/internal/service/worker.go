package service

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/yyhuni/orbit/server/internal/dto"
	"github.com/yyhuni/orbit/server/internal/model"
	"github.com/yyhuni/orbit/server/internal/repository"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
)

var (
	ErrWorkerScanNotFound = errors.New("scan not found")
	ErrWorkerToolRequired = errors.New("tool parameter required for provider_config")
)

// WorkerService handles scan data for workers
type WorkerService struct {
	scanRepo                       *repository.ScanRepository
	subfinderProviderSettingsRepo  *repository.SubfinderProviderSettingsRepository
}

// NewWorkerService creates a new worker service
func NewWorkerService(
	scanRepo *repository.ScanRepository,
	subfinderProviderSettingsRepo *repository.SubfinderProviderSettingsRepository,
) *WorkerService {
	return &WorkerService{
		scanRepo:                      scanRepo,
		subfinderProviderSettingsRepo: subfinderProviderSettingsRepo,
	}
}

// TargetInfo contains target name and type
type TargetInfo struct {
	Name string
	Type string
}

// GetTargetName returns target name and type for a scan
func (s *WorkerService) GetTargetName(scanID int) (*TargetInfo, error) {
	target, err := s.scanRepo.GetTargetByScanID(scanID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrWorkerScanNotFound
		}
		return nil, err
	}

	return &TargetInfo{
		Name: target.Name,
		Type: target.Type,
	}, nil
}

// GetProviderConfig returns provider config (API keys) for a tool
func (s *WorkerService) GetProviderConfig(scanID int, toolName string) (*dto.WorkerProviderConfigResponse, error) {
	if toolName == "" {
		return nil, ErrWorkerToolRequired
	}

	// Check scan exists
	_, err := s.scanRepo.FindByID(scanID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrWorkerScanNotFound
		}
		return nil, err
	}

	config, err := s.generateProviderConfig(toolName)
	if err != nil {
		return nil, err
	}

	return &dto.WorkerProviderConfigResponse{
		Content: config,
	}, nil
}

// generateProviderConfig generates provider config YAML for a tool
func (s *WorkerService) generateProviderConfig(toolName string) (string, error) {
	switch toolName {
	case "subfinder":
		return s.generateSubfinderConfig()
	default:
		return "", nil
	}
}

// generateSubfinderConfig generates subfinder provider-config.yaml content
func (s *WorkerService) generateSubfinderConfig() (string, error) {
	settings, err := s.subfinderProviderSettingsRepo.GetInstance()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil // No settings configured
		}
		return "", err
	}

	config := make(map[string][]string)
	hasEnabled := false

	for providerName, formatInfo := range model.ProviderFormats {
		providerConfig, exists := settings.Providers[providerName]
		if !exists || !providerConfig.Enabled {
			config[providerName] = []string{}
			continue
		}

		value := s.buildProviderValue(providerConfig, formatInfo)
		if value != "" {
			config[providerName] = []string{value}
			hasEnabled = true
		} else {
			config[providerName] = []string{}
		}
	}

	if !hasEnabled {
		return "", nil
	}

	yamlBytes, err := yaml.Marshal(config)
	if err != nil {
		return "", fmt.Errorf("failed to marshal provider config: %w", err)
	}

	return string(yamlBytes), nil
}

// buildProviderValue builds the config value string for a provider
func (s *WorkerService) buildProviderValue(config model.ProviderConfig, formatInfo model.ProviderFormat) string {
	if formatInfo.Type == model.FormatTypeComposite {
		// Handle composite formats like "email:api_key"
		result := formatInfo.Format
		result = strings.ReplaceAll(result, "{email}", config.Email)
		result = strings.ReplaceAll(result, "{api_key}", config.APIKey)
		result = strings.ReplaceAll(result, "{api_id}", config.APIId)
		result = strings.ReplaceAll(result, "{api_secret}", config.APISecret)

		// Check if all placeholders were replaced (no empty values)
		if strings.Contains(result, "{}") || result == formatInfo.Format {
			return ""
		}
		// Check for empty segments (e.g., ":key" or "email:")
		if slices.Contains(strings.Split(result, ":"), "") {
			return ""
		}
		return result
	}

	// Single field format
	return config.APIKey
}


