package repository

import (
	catalogdomain "github.com/yyhuni/lunafox/server/internal/modules/catalog/domain"
	"github.com/yyhuni/lunafox/server/internal/modules/catalog/repository/persistence"
)

func subfinderProviderSettingsModelToDomain(settings *model.SubfinderProviderSettings) *catalogdomain.SubfinderProviderSettings {
	if settings == nil {
		return nil
	}

	providers := make(catalogdomain.ProviderConfigs, len(settings.Providers))
	for providerName, providerConfig := range settings.Providers {
		providers[providerName] = catalogdomain.ProviderConfig{
			Enabled:   providerConfig.Enabled,
			Email:     providerConfig.Email,
			APIKey:    providerConfig.APIKey,
			APIID:     providerConfig.APIId,
			APISecret: providerConfig.APISecret,
		}
	}

	return &catalogdomain.SubfinderProviderSettings{
		ID:        settings.ID,
		Providers: providers,
	}
}

func subfinderProviderSettingsDomainToModel(settings *catalogdomain.SubfinderProviderSettings) *model.SubfinderProviderSettings {
	if settings == nil {
		return nil
	}

	providers := make(model.ProviderConfigs, len(settings.Providers))
	for providerName, providerConfig := range settings.Providers {
		providers[providerName] = model.ProviderConfig{
			Enabled:   providerConfig.Enabled,
			Email:     providerConfig.Email,
			APIKey:    providerConfig.APIKey,
			APIId:     providerConfig.APIID,
			APISecret: providerConfig.APISecret,
		}
	}

	return &model.SubfinderProviderSettings{
		ID:        settings.ID,
		Providers: providers,
	}
}
