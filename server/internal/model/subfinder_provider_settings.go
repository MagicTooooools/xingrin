package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// SubfinderProviderSettings stores API keys for subfinder data sources (singleton, id=1)
type SubfinderProviderSettings struct {
	ID        int             `gorm:"primaryKey" json:"id"`
	Providers ProviderConfigs `gorm:"column:providers;type:jsonb" json:"providers"`
	CreatedAt time.Time       `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time       `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (SubfinderProviderSettings) TableName() string {
	return "subfinder_provider_settings"
}

// ProviderConfigs maps provider name to its configuration
type ProviderConfigs map[string]ProviderConfig

// ProviderConfig holds credentials for a single provider
type ProviderConfig struct {
	Enabled   bool   `json:"enabled"`
	Email     string `json:"email,omitempty"`
	APIKey    string `json:"api_key,omitempty"`
	APIId     string `json:"api_id,omitempty"`
	APISecret string `json:"api_secret,omitempty"`
}

// Scan implements sql.Scanner for GORM
func (p *ProviderConfigs) Scan(value any) error {
	if value == nil {
		*p = make(ProviderConfigs)
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("ProviderConfigs: expected []byte from database")
	}
	return json.Unmarshal(bytes, p)
}

// Value implements driver.Valuer for GORM
func (p ProviderConfigs) Value() (driver.Value, error) {
	return json.Marshal(p)
}

// ProviderFormatType defines how provider credentials are formatted
type ProviderFormatType string

const (
	FormatTypeSingle    ProviderFormatType = "single"
	FormatTypeComposite ProviderFormatType = "composite"
)

// ProviderFormat defines the credential format for a provider
type ProviderFormat struct {
	Type   ProviderFormatType
	Format string // field name for single, template for composite (e.g., "{email}:{api_key}")
}

// ProviderFormats defines credential formats for generating subfinder config YAML
var ProviderFormats = map[string]ProviderFormat{
	"fofa":           {Type: FormatTypeComposite, Format: "{email}:{api_key}"},
	"censys":         {Type: FormatTypeComposite, Format: "{api_id}:{api_secret}"},
	"hunter":         {Type: FormatTypeSingle, Format: "api_key"},
	"shodan":         {Type: FormatTypeSingle, Format: "api_key"},
	"zoomeye":        {Type: FormatTypeSingle, Format: "api_key"},
	"securitytrails": {Type: FormatTypeSingle, Format: "api_key"},
	"threatbook":     {Type: FormatTypeSingle, Format: "api_key"},
	"quake":          {Type: FormatTypeSingle, Format: "api_key"},
}
