package dto

import "github.com/yyhuni/lunafox/server/internal/preset"

// PresetResponse represents the API response for a preset engine.
// Note: enabledFeatures is parsed by frontend from configuration,
// keeping consistency with user engines.
type PresetResponse struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description,omitempty"`
	Configuration string `json:"configuration"`
}

// NewPresetResponse creates a PresetResponse from a preset.Preset.
func NewPresetResponse(p *preset.Preset) PresetResponse {
	return PresetResponse{
		ID:            p.ID,
		Name:          p.Name,
		Description:   p.Description,
		Configuration: p.Configuration,
	}
}

// NewPresetListResponse creates a slice of PresetResponse from a slice of presets.
func NewPresetListResponse(presets []preset.Preset) []PresetResponse {
	responses := make([]PresetResponse, len(presets))
	for i := range presets {
		responses[i] = NewPresetResponse(&presets[i])
	}
	return responses
}
