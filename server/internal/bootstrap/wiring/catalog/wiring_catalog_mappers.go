package catalogwiring

import (
	catalogdomain "github.com/yyhuni/lunafox/server/internal/modules/catalog/domain"
	catalogmodel "github.com/yyhuni/lunafox/server/internal/modules/catalog/repository/persistence"
)

func catalogModelEngineToDomain(engine *catalogmodel.ScanEngine) *catalogdomain.ScanEngine {
	if engine == nil {
		return nil
	}
	return &catalogdomain.ScanEngine{
		ID:            engine.ID,
		Name:          engine.Name,
		Configuration: engine.Configuration,
		CreatedAt:     engine.CreatedAt,
		UpdatedAt:     engine.UpdatedAt,
	}
}

func catalogDomainEngineToModel(engine *catalogdomain.ScanEngine) *catalogmodel.ScanEngine {
	if engine == nil {
		return nil
	}
	return &catalogmodel.ScanEngine{
		ID:            engine.ID,
		Name:          engine.Name,
		Configuration: engine.Configuration,
		CreatedAt:     engine.CreatedAt,
		UpdatedAt:     engine.UpdatedAt,
	}
}

func catalogModelEngineListToDomain(engines []catalogmodel.ScanEngine) []catalogdomain.ScanEngine {
	results := make([]catalogdomain.ScanEngine, 0, len(engines))
	for index := range engines {
		results = append(results, *catalogModelEngineToDomain(&engines[index]))
	}
	return results
}

func catalogModelTargetToDomain(target *catalogmodel.Target) *catalogdomain.Target {
	if target == nil {
		return nil
	}

	organizations := make([]catalogdomain.TargetOrganizationRef, 0, len(target.Organizations))
	for _, organization := range target.Organizations {
		organizations = append(organizations, catalogdomain.TargetOrganizationRef{
			ID:          organization.ID,
			Name:        organization.Name,
			Description: organization.Description,
			CreatedAt:   organization.CreatedAt,
			DeletedAt:   organization.DeletedAt,
		})
	}

	return &catalogdomain.Target{
		ID:            target.ID,
		Name:          target.Name,
		Type:          target.Type,
		CreatedAt:     target.CreatedAt,
		LastScannedAt: target.LastScannedAt,
		DeletedAt:     target.DeletedAt,
		Organizations: organizations,
	}
}

func catalogDomainTargetToModel(target *catalogdomain.Target) *catalogmodel.Target {
	if target == nil {
		return nil
	}

	organizations := make([]catalogmodel.TargetOrganizationRef, 0, len(target.Organizations))
	for _, organization := range target.Organizations {
		organizations = append(organizations, catalogmodel.TargetOrganizationRef{
			ID:          organization.ID,
			Name:        organization.Name,
			Description: organization.Description,
			CreatedAt:   organization.CreatedAt,
			DeletedAt:   organization.DeletedAt,
		})
	}

	return &catalogmodel.Target{
		ID:            target.ID,
		Name:          target.Name,
		Type:          target.Type,
		CreatedAt:     target.CreatedAt,
		LastScannedAt: target.LastScannedAt,
		DeletedAt:     target.DeletedAt,
		Organizations: organizations,
	}
}

func catalogModelTargetListToDomain(targets []catalogmodel.Target) []catalogdomain.Target {
	results := make([]catalogdomain.Target, 0, len(targets))
	for index := range targets {
		results = append(results, *catalogModelTargetToDomain(&targets[index]))
	}
	return results
}

func catalogDomainTargetListToModel(targets []catalogdomain.Target) []catalogmodel.Target {
	results := make([]catalogmodel.Target, 0, len(targets))
	for index := range targets {
		results = append(results, *catalogDomainTargetToModel(&targets[index]))
	}
	return results
}

func catalogModelWordlistToDomain(wordlist *catalogmodel.Wordlist) *catalogdomain.Wordlist {
	if wordlist == nil {
		return nil
	}
	return &catalogdomain.Wordlist{
		ID:          wordlist.ID,
		Name:        wordlist.Name,
		Description: wordlist.Description,
		FilePath:    wordlist.FilePath,
		FileSize:    wordlist.FileSize,
		LineCount:   wordlist.LineCount,
		FileHash:    wordlist.FileHash,
		CreatedAt:   wordlist.CreatedAt,
		UpdatedAt:   wordlist.UpdatedAt,
	}
}

func catalogDomainWordlistToModel(wordlist *catalogdomain.Wordlist) *catalogmodel.Wordlist {
	if wordlist == nil {
		return nil
	}
	return &catalogmodel.Wordlist{
		ID:          wordlist.ID,
		Name:        wordlist.Name,
		Description: wordlist.Description,
		FilePath:    wordlist.FilePath,
		FileSize:    wordlist.FileSize,
		LineCount:   wordlist.LineCount,
		FileHash:    wordlist.FileHash,
		CreatedAt:   wordlist.CreatedAt,
		UpdatedAt:   wordlist.UpdatedAt,
	}
}

func catalogModelWordlistListToDomain(wordlists []catalogmodel.Wordlist) []catalogdomain.Wordlist {
	results := make([]catalogdomain.Wordlist, 0, len(wordlists))
	for index := range wordlists {
		results = append(results, *catalogModelWordlistToDomain(&wordlists[index]))
	}
	return results
}
