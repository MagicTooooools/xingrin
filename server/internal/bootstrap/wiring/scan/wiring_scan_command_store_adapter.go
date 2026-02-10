package scanwiring

import (
	scandomain "github.com/yyhuni/lunafox/server/internal/modules/scan/domain"
	scanrepo "github.com/yyhuni/lunafox/server/internal/modules/scan/repository"
)

func newScanDomainCommandStoreAdapter(repo *scanrepo.ScanRepository) scandomain.ScanRepository {
	return scanrepo.NewDomainScanRepository(repo)
}
