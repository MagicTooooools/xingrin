package snapshotwiring

import (
	securitydto "github.com/yyhuni/lunafox/server/internal/modules/security/dto"
	snapshotapp "github.com/yyhuni/lunafox/server/internal/modules/snapshot/application"
)

func snapshotVulnerabilityAssetCreateItemsToDTO(items []snapshotapp.VulnerabilityAssetCreateItem) []securitydto.VulnerabilityCreateItem {
	results := make([]securitydto.VulnerabilityCreateItem, 0, len(items))
	for index := range items {
		item := items[index]
		results = append(results, securitydto.VulnerabilityCreateItem{
			URL:         item.URL,
			VulnType:    item.VulnType,
			Severity:    item.Severity,
			Source:      item.Source,
			CVSSScore:   item.CVSSScore,
			Description: item.Description,
			RawOutput:   item.RawOutput,
		})
	}
	return results
}
