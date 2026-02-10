package application

import (
	"errors"
	"strings"
	"testing"

	catalogdomain "github.com/yyhuni/lunafox/server/internal/modules/catalog/domain"
)

type workerScanGuardStub struct {
	called bool
	lastID int
	err    error
}

func (stub *workerScanGuardStub) EnsureActiveByID(id int) error {
	stub.called = true
	stub.lastID = id
	return stub.err
}

type workerSettingsStoreStub struct {
	settings *catalogdomain.SubfinderProviderSettings
	err      error
}

func (stub *workerSettingsStoreStub) GetInstance() (*catalogdomain.SubfinderProviderSettings, error) {
	if stub.err != nil {
		return nil, stub.err
	}
	return stub.settings, nil
}

func TestWorkerServiceGetProviderConfigToolRequired(t *testing.T) {
	service := NewWorkerService(&workerScanGuardStub{}, &workerSettingsStoreStub{})

	_, err := service.GetProviderConfig(1, "  ")
	if !errors.Is(err, ErrWorkerToolRequired) {
		t.Fatalf("expected ErrWorkerToolRequired, got %v", err)
	}
}

func TestWorkerServiceGetProviderConfigScanGuardError(t *testing.T) {
	guard := &workerScanGuardStub{err: ErrWorkerScanNotFound}
	service := NewWorkerService(guard, &workerSettingsStoreStub{})

	_, err := service.GetProviderConfig(9, "subfinder")
	if !errors.Is(err, ErrWorkerScanNotFound) {
		t.Fatalf("expected ErrWorkerScanNotFound, got %v", err)
	}
	if !guard.called || guard.lastID != 9 {
		t.Fatalf("expected scan guard called with id=9, got called=%v id=%d", guard.called, guard.lastID)
	}
}

func TestWorkerServiceGetProviderConfigSettingsNotFound(t *testing.T) {
	service := NewWorkerService(&workerScanGuardStub{}, &workerSettingsStoreStub{err: ErrWorkerProviderSettingsNotFound})

	config, err := service.GetProviderConfig(1, "subfinder")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if config != "" {
		t.Fatalf("expected empty config, got %q", config)
	}
}

func TestWorkerServiceGetProviderConfigSubfinder(t *testing.T) {
	guard := &workerScanGuardStub{}
	settings := &catalogdomain.SubfinderProviderSettings{Providers: catalogdomain.ProviderConfigs{
		"fofa": {Enabled: true, Email: "test@example.com", APIKey: "secret"},
	}}
	service := NewWorkerService(guard, &workerSettingsStoreStub{settings: settings})

	config, err := service.GetProviderConfig(7, "subfinder")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !guard.called || guard.lastID != 7 {
		t.Fatalf("expected scan guard called with id=7, got called=%v id=%d", guard.called, guard.lastID)
	}
	if !strings.Contains(config, "fofa:") {
		t.Fatalf("expected config contains fofa section, got %q", config)
	}
	if !strings.Contains(config, "test@example.com:secret") {
		t.Fatalf("expected fofa credential in config, got %q", config)
	}
}
