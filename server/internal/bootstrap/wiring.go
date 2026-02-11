package bootstrap

import (
	assetwiring "github.com/yyhuni/lunafox/server/internal/bootstrap/wiring/asset"
	catalogwiring "github.com/yyhuni/lunafox/server/internal/bootstrap/wiring/catalog"
	identitywiring "github.com/yyhuni/lunafox/server/internal/bootstrap/wiring/identity"
	scanwiring "github.com/yyhuni/lunafox/server/internal/bootstrap/wiring/scan"
	scanlogwiring "github.com/yyhuni/lunafox/server/internal/bootstrap/wiring/scanlog"
	securitywiring "github.com/yyhuni/lunafox/server/internal/bootstrap/wiring/security"
	snapshotwiring "github.com/yyhuni/lunafox/server/internal/bootstrap/wiring/snapshot"
	workerwiring "github.com/yyhuni/lunafox/server/internal/bootstrap/wiring/worker"
	"github.com/yyhuni/lunafox/server/internal/config"
	agentservice "github.com/yyhuni/lunafox/server/internal/modules/agent/application"
	agentdomain "github.com/yyhuni/lunafox/server/internal/modules/agent/domain"
	agenthandler "github.com/yyhuni/lunafox/server/internal/modules/agent/handler"
	agentinfra "github.com/yyhuni/lunafox/server/internal/modules/agent/infrastructure"
	agentrepo "github.com/yyhuni/lunafox/server/internal/modules/agent/repository"
	assetservice "github.com/yyhuni/lunafox/server/internal/modules/asset/application"
	assethandler "github.com/yyhuni/lunafox/server/internal/modules/asset/handler"
	directoryhandler "github.com/yyhuni/lunafox/server/internal/modules/asset/handler/directory"
	endpointhandler "github.com/yyhuni/lunafox/server/internal/modules/asset/handler/endpoint"
	hostporthandler "github.com/yyhuni/lunafox/server/internal/modules/asset/handler/host_port"
	screenshothandler "github.com/yyhuni/lunafox/server/internal/modules/asset/handler/screenshot"
	subdomainhandler "github.com/yyhuni/lunafox/server/internal/modules/asset/handler/subdomain"
	websitehandler "github.com/yyhuni/lunafox/server/internal/modules/asset/handler/website"
	assetrepo "github.com/yyhuni/lunafox/server/internal/modules/asset/repository"
	catalogservice "github.com/yyhuni/lunafox/server/internal/modules/catalog/application"
	cataloghandler "github.com/yyhuni/lunafox/server/internal/modules/catalog/handler"
	catalogrepo "github.com/yyhuni/lunafox/server/internal/modules/catalog/repository"
	identityservice "github.com/yyhuni/lunafox/server/internal/modules/identity/application"
	identityhandler "github.com/yyhuni/lunafox/server/internal/modules/identity/handler"
	identityrepo "github.com/yyhuni/lunafox/server/internal/modules/identity/repository"
	scanhandler "github.com/yyhuni/lunafox/server/internal/modules/scan/handler"
	scanrepo "github.com/yyhuni/lunafox/server/internal/modules/scan/repository"
	securityservice "github.com/yyhuni/lunafox/server/internal/modules/security/application"
	securityhandler "github.com/yyhuni/lunafox/server/internal/modules/security/handler"
	securityrepo "github.com/yyhuni/lunafox/server/internal/modules/security/repository"
	snapshothandler "github.com/yyhuni/lunafox/server/internal/modules/snapshot/handler"
	snapshotrepo "github.com/yyhuni/lunafox/server/internal/modules/snapshot/repository"
	"github.com/yyhuni/lunafox/server/internal/preset"
	ws "github.com/yyhuni/lunafox/server/internal/websocket"
)

type deps struct {
	healthHandler        *assethandler.HealthHandler
	authHandler          *identityhandler.AuthHandler
	userHandler          *identityhandler.UserHandler
	orgHandler           *identityhandler.OrganizationHandler
	targetHandler        *cataloghandler.TargetHandler
	engineHandler        *cataloghandler.EngineHandler
	wordlistHandler      *cataloghandler.WordlistHandler
	websiteHandler       *websitehandler.WebsiteHandler
	subdomainHandler     *subdomainhandler.SubdomainHandler
	endpointHandler      *endpointhandler.EndpointHandler
	directoryHandler     *directoryhandler.DirectoryHandler
	hostPortHandler      *hostporthandler.HostPortHandler
	screenshotHandler    *screenshothandler.ScreenshotHandler
	vulnerabilityHandler *securityhandler.VulnerabilityHandler
	scanHandler          *scanhandler.ScanHandler
	scanLogHandler       *scanhandler.ScanLogHandler
	workerHandler        *cataloghandler.WorkerHandler
	workerScanHandler    *scanhandler.WorkerScanHandler

	agentHandler     *agenthandler.AgentHandler
	agentWSHandler   *agenthandler.AgentWebSocketHandler
	agentTaskHandler *agenthandler.AgentTaskHandler

	websiteSnapshotHandler       *snapshothandler.WebsiteSnapshotHandler
	subdomainSnapshotHandler     *snapshothandler.SubdomainSnapshotHandler
	endpointSnapshotHandler      *snapshothandler.EndpointSnapshotHandler
	directorySnapshotHandler     *snapshothandler.DirectorySnapshotHandler
	hostPortSnapshotHandler      *snapshothandler.HostPortSnapshotHandler
	screenshotSnapshotHandler    *snapshothandler.ScreenshotSnapshotHandler
	vulnerabilitySnapshotHandler *snapshothandler.VulnerabilitySnapshotHandler
	presetHandler                *cataloghandler.PresetHandler

	agentRepo    agentdomain.AgentRepository
	scanTaskRepo scanrepo.ScanTaskRepository
}

func buildDependencies(infra *infra, cfg *config.Config) *deps {
	db := infra.db

	// Base repositories
	userRepo := identityrepo.NewUserRepository(db)
	orgRepo := identityrepo.NewOrganizationRepository(db)
	targetRepo := catalogrepo.NewTargetRepository(db)
	engineRepo := catalogrepo.NewEngineRepository(db)
	wordlistRepo := catalogrepo.NewWordlistRepository(db)
	websiteRepo := assetrepo.NewWebsiteRepository(db)
	subdomainRepo := assetrepo.NewSubdomainRepository(db)
	endpointRepo := assetrepo.NewEndpointRepository(db)
	directoryRepo := assetrepo.NewDirectoryRepository(db)
	hostPortRepo := assetrepo.NewHostPortRepository(db)
	screenshotRepo := assetrepo.NewScreenshotRepository(db)
	vulnerabilityRepo := securityrepo.NewVulnerabilityRepository(db)
	scanRepo := scanrepo.NewScanRepository(db)
	scanLogRepo := scanrepo.NewScanLogRepository(db)
	subfinderProviderSettingsRepo := catalogrepo.NewSubfinderProviderSettingsRepository(db)
	websiteSnapshotRepo := snapshotrepo.NewWebsiteSnapshotRepository(db)
	subdomainSnapshotRepo := snapshotrepo.NewSubdomainSnapshotRepository(db)
	endpointSnapshotRepo := snapshotrepo.NewEndpointSnapshotRepository(db)
	directorySnapshotRepo := snapshotrepo.NewDirectorySnapshotRepository(db)
	hostPortSnapshotRepo := snapshotrepo.NewHostPortSnapshotRepository(db)
	screenshotSnapshotRepo := snapshotrepo.NewScreenshotSnapshotRepository(db)
	vulnerabilitySnapshotRepo := snapshotrepo.NewVulnerabilitySnapshotRepository(db)

	// Agent-specific repositories
	agentRepo := agentrepo.NewAgentRepository(db)
	registrationTokenRepo := agentrepo.NewRegistrationTokenRepository(db)
	scanTaskRepo := scanrepo.NewScanTaskRepository(db)

	// Identity module: stores and services
	identityUserQueryStore := identitywiring.NewIdentityUserQueryStoreAdapter(userRepo)
	identityUserCommandStore := identitywiring.NewIdentityUserCommandStoreAdapter(userRepo)
	identityOrgQueryStore := identitywiring.NewIdentityOrganizationQueryStoreAdapter(orgRepo)
	identityOrgCommandStore := identitywiring.NewIdentityOrganizationCommandStoreAdapter(orgRepo)
	identityAuthUserStore := identitywiring.NewIdentityAuthUserStoreAdapter(userRepo)

	userSvc := identityservice.NewUserFacade(identityUserQueryStore, identityUserCommandStore)
	orgSvc := identityservice.NewOrganizationFacade(identityOrgQueryStore, identityOrgCommandStore)
	authSvc := identityservice.NewAuthFacade(identityAuthUserStore, infra.jwtManager)

	// Catalog module: stores and services
	catalogTargetQueryStore := catalogwiring.NewCatalogTargetQueryStoreAdapter(targetRepo)
	catalogTargetCommandStore := catalogwiring.NewCatalogTargetCommandStoreAdapter(targetRepo)
	catalogEngineQueryStore := catalogwiring.NewCatalogEngineQueryStoreAdapter(engineRepo)
	catalogEngineCommandStore := catalogwiring.NewCatalogEngineCommandStoreAdapter(engineRepo)
	catalogWordlistQueryStore := catalogwiring.NewCatalogWordlistQueryStoreAdapter(wordlistRepo)
	catalogWordlistCommandStore := catalogwiring.NewCatalogWordlistCommandStoreAdapter(wordlistRepo)
	catalogOrganizationTargetBindingStore := catalogwiring.NewCatalogOrganizationTargetBindingStoreAdapter(orgRepo)

	targetSvc := catalogservice.NewTargetFacade(catalogTargetQueryStore, catalogTargetCommandStore, catalogOrganizationTargetBindingStore)
	engineSvc := catalogservice.NewEngineFacade(catalogEngineQueryStore, catalogEngineCommandStore)
	wordlistSvc := catalogservice.NewWordlistFacade(catalogWordlistQueryStore, catalogWordlistCommandStore, cfg.Storage.WordlistsBasePath)

	// Asset module: stores and services
	assetTargetLookup := assetwiring.NewAssetTargetLookupAdapter(targetRepo)
	assetWebsiteStore := assetwiring.NewAssetWebsiteStoreAdapter(websiteRepo)
	assetSubdomainStore := assetwiring.NewAssetSubdomainStoreAdapter(subdomainRepo)
	assetEndpointStore := assetwiring.NewAssetEndpointStoreAdapter(endpointRepo)
	assetDirectoryStore := assetwiring.NewAssetDirectoryStoreAdapter(directoryRepo)
	assetHostPortStore := assetwiring.NewAssetHostPortStoreAdapter(hostPortRepo)
	assetScreenshotStore := assetwiring.NewAssetScreenshotStoreAdapter(screenshotRepo)

	websiteSvc := assetservice.NewWebsiteFacade(assetWebsiteStore, assetTargetLookup)
	subdomainSvc := assetservice.NewSubdomainFacade(assetSubdomainStore, assetTargetLookup)
	endpointSvc := assetservice.NewEndpointFacade(assetEndpointStore, assetTargetLookup)
	directorySvc := assetservice.NewDirectoryFacade(assetDirectoryStore, assetTargetLookup)
	hostPortSvc := assetservice.NewHostPortFacade(assetHostPortStore, assetTargetLookup)
	screenshotSvc := assetservice.NewScreenshotFacade(assetScreenshotStore, assetTargetLookup)

	// Security module
	securityVulnerabilityStore := securitywiring.NewSecurityVulnerabilityStoreAdapter(vulnerabilityRepo)
	securityTargetLookup := securitywiring.NewSecurityTargetLookupAdapter(targetRepo)
	vulnerabilitySvc := securityservice.NewVulnerabilityFacade(securityVulnerabilityStore, securityTargetLookup)

	// Scan and scan-log modules
	scanQueryStore := scanwiring.NewScanQueryStoreAdapter(scanRepo)
	scanCommandStore := scanwiring.NewScanCommandStoreAdapter(scanRepo)
	scanDomainRepository := scanwiring.NewScanDomainRepositoryAdapter(scanRepo)
	scanTaskStore := scanwiring.NewScanTaskStoreAdapter(scanTaskRepo)
	scanTaskRuntimeStore := scanwiring.NewScanTaskRuntimeStoreAdapter(scanRepo)
	scanLogQueryStore := scanlogwiring.NewScanLogQueryStoreAdapter(scanLogRepo)
	scanLogCommandStore := scanlogwiring.NewScanLogCommandStoreAdapter(scanLogRepo)

	scanTaskCanceller := scanwiring.NewScanTaskCancellerAdapter(scanTaskRepo)
	scanTargetLookup := scanwiring.NewScanTargetLookupAdapter(targetRepo)
	scanLogLookup := scanlogwiring.NewScanLogScanLookupAdapter(scanRepo)

	scanSvc := scanwiring.NewScanApplicationService(
		scanQueryStore,
		scanCommandStore,
		scanDomainRepository,
		scanTaskCanceller,
		infra.wsHub,
		scanTargetLookup,
	)
	scanTaskSvc := scanwiring.NewScanTaskApplicationService(scanTaskStore, scanTaskRuntimeStore)
	scanLogSvc := scanlogwiring.NewScanLogApplicationService(scanLogQueryStore, scanLogCommandStore, scanLogLookup)

	// Worker module
	workerScanGuard := workerwiring.NewWorkerScanGuardAdapter(scanRepo)
	workerSettingsStore := workerwiring.NewWorkerProviderSettingsStoreAdapter(subfinderProviderSettingsRepo)
	workerSvc := workerwiring.NewWorkerApplicationService(workerScanGuard, workerSettingsStore)
	// Agent module services
	agentClock := agentinfra.NewSystemClock()
	agentTokenGenerator := agentinfra.NewCryptoTokenGenerator()
	agentSvc := agentservice.NewAgentFacade(agentRepo, registrationTokenRepo, agentClock, agentTokenGenerator)
	agentRuntimeSvc := agentservice.NewAgentRuntimeService(agentRepo, infra.heartbeatCache, ws.NewAgentMessagePublisher(infra.wsHub), agentClock, infra.serverVersion, infra.agentImage)
	agentTaskSvc := agentservice.NewAgentTaskService(scanTaskSvc)

	// Snapshot module: lookup, stores, sync adapters, and services
	snapshotScanLookup := snapshotwiring.NewSnapshotScanRefLookupAdapter(scanRepo)

	websiteSnapshotQueryStore := snapshotwiring.NewSnapshotWebsiteQueryStoreAdapter(websiteSnapshotRepo)
	subdomainSnapshotQueryStore := snapshotwiring.NewSnapshotSubdomainQueryStoreAdapter(subdomainSnapshotRepo)
	endpointSnapshotQueryStore := snapshotwiring.NewSnapshotEndpointQueryStoreAdapter(endpointSnapshotRepo)
	directorySnapshotQueryStore := snapshotwiring.NewSnapshotDirectoryQueryStoreAdapter(directorySnapshotRepo)
	hostPortSnapshotQueryStore := snapshotwiring.NewSnapshotHostPortQueryStoreAdapter(hostPortSnapshotRepo)
	screenshotSnapshotQueryStore := snapshotwiring.NewSnapshotScreenshotQueryStoreAdapter(screenshotSnapshotRepo)
	vulnerabilitySnapshotQueryStore := snapshotwiring.NewSnapshotVulnerabilityQueryStoreAdapter(vulnerabilitySnapshotRepo)

	websiteSnapshotCommandStore := snapshotwiring.NewSnapshotWebsiteCommandStoreAdapter(websiteSnapshotRepo)
	subdomainSnapshotCommandStore := snapshotwiring.NewSnapshotSubdomainCommandStoreAdapter(subdomainSnapshotRepo)
	endpointSnapshotCommandStore := snapshotwiring.NewSnapshotEndpointCommandStoreAdapter(endpointSnapshotRepo)
	directorySnapshotCommandStore := snapshotwiring.NewSnapshotDirectoryCommandStoreAdapter(directorySnapshotRepo)
	hostPortSnapshotCommandStore := snapshotwiring.NewSnapshotHostPortCommandStoreAdapter(hostPortSnapshotRepo)
	screenshotSnapshotCommandStore := snapshotwiring.NewSnapshotScreenshotCommandStoreAdapter(screenshotSnapshotRepo)
	vulnerabilitySnapshotCommandStore := snapshotwiring.NewSnapshotVulnerabilityCommandStoreAdapter(vulnerabilitySnapshotRepo)

	websiteAssetSync := snapshotwiring.NewSnapshotWebsiteAssetSyncAdapter(websiteSvc)
	subdomainAssetSync := snapshotwiring.NewSnapshotSubdomainAssetSyncAdapter(subdomainSvc)
	endpointAssetSync := snapshotwiring.NewSnapshotEndpointAssetSyncAdapter(endpointSvc)
	directoryAssetSync := snapshotwiring.NewSnapshotDirectoryAssetSyncAdapter(directorySvc)
	hostPortAssetSync := snapshotwiring.NewSnapshotHostPortAssetSyncAdapter(hostPortSvc)
	screenshotAssetSync := snapshotwiring.NewSnapshotScreenshotAssetSyncAdapter(screenshotSvc)
	vulnerabilityAssetSync := snapshotwiring.NewSnapshotVulnerabilityAssetSyncAdapter(vulnerabilitySvc)
	vulnerabilityRawOutputCodec := snapshotwiring.NewSnapshotVulnerabilityRawOutputCodec()

	websiteSnapshotSvc := snapshotwiring.NewSnapshotWebsiteApplicationService(websiteSnapshotQueryStore, websiteSnapshotCommandStore, snapshotScanLookup, websiteAssetSync)
	subdomainSnapshotSvc := snapshotwiring.NewSnapshotSubdomainApplicationService(subdomainSnapshotQueryStore, subdomainSnapshotCommandStore, snapshotScanLookup, subdomainAssetSync)
	endpointSnapshotSvc := snapshotwiring.NewSnapshotEndpointApplicationService(endpointSnapshotQueryStore, endpointSnapshotCommandStore, snapshotScanLookup, endpointAssetSync)
	directorySnapshotSvc := snapshotwiring.NewSnapshotDirectoryApplicationService(directorySnapshotQueryStore, directorySnapshotCommandStore, snapshotScanLookup, directoryAssetSync)
	hostPortSnapshotSvc := snapshotwiring.NewSnapshotHostPortApplicationService(hostPortSnapshotQueryStore, hostPortSnapshotCommandStore, snapshotScanLookup, hostPortAssetSync)
	screenshotSnapshotSvc := snapshotwiring.NewSnapshotScreenshotApplicationService(screenshotSnapshotQueryStore, screenshotSnapshotCommandStore, snapshotScanLookup, screenshotAssetSync)
	vulnerabilitySnapshotSvc := snapshotwiring.NewSnapshotVulnerabilityApplicationService(vulnerabilitySnapshotQueryStore, vulnerabilitySnapshotCommandStore, snapshotScanLookup, vulnerabilityAssetSync, vulnerabilityRawOutputCodec)
	// Preset module
	presetSvc := preset.NewService(infra.presetLoader)

	// HTTP handlers and exposed dependencies
	return &deps{
		healthHandler:        assethandler.NewHealthHandler(db, infra.redisClient),
		authHandler:          identityhandler.NewAuthHandler(authSvc),
		userHandler:          identityhandler.NewUserHandler(userSvc),
		orgHandler:           identityhandler.NewOrganizationHandler(orgSvc),
		targetHandler:        cataloghandler.NewTargetHandler(targetSvc),
		engineHandler:        cataloghandler.NewEngineHandler(engineSvc),
		wordlistHandler:      cataloghandler.NewWordlistHandler(wordlistSvc),
		websiteHandler:       websitehandler.NewWebsiteHandler(websiteSvc),
		subdomainHandler:     subdomainhandler.NewSubdomainHandler(subdomainSvc),
		endpointHandler:      endpointhandler.NewEndpointHandler(endpointSvc),
		directoryHandler:     directoryhandler.NewDirectoryHandler(directorySvc),
		hostPortHandler:      hostporthandler.NewHostPortHandler(hostPortSvc),
		screenshotHandler:    screenshothandler.NewScreenshotHandler(screenshotSvc),
		vulnerabilityHandler: securityhandler.NewVulnerabilityHandler(vulnerabilitySvc),
		scanHandler:          scanhandler.NewScanHandler(scanSvc),
		scanLogHandler:       scanhandler.NewScanLogHandler(scanLogSvc),
		workerHandler:        cataloghandler.NewWorkerHandler(workerSvc),
		workerScanHandler:    scanhandler.NewWorkerScanHandler(scanSvc),

		agentHandler: agenthandler.NewAgentHandler(
			agentSvc,
			agentRuntimeSvc,
			cfg.PublicURL,
			infra.serverVersion,
			infra.agentImage,
			cfg.Worker.Token,
			infra.heartbeatCache,
		),
		agentWSHandler: agenthandler.NewAgentWebSocketHandler(
			infra.wsHub,
			agentRuntimeSvc,
		),
		agentTaskHandler: agenthandler.NewAgentTaskHandler(agentTaskSvc),

		websiteSnapshotHandler:       snapshothandler.NewWebsiteSnapshotHandler(websiteSnapshotSvc),
		subdomainSnapshotHandler:     snapshothandler.NewSubdomainSnapshotHandler(subdomainSnapshotSvc),
		endpointSnapshotHandler:      snapshothandler.NewEndpointSnapshotHandler(endpointSnapshotSvc),
		directorySnapshotHandler:     snapshothandler.NewDirectorySnapshotHandler(directorySnapshotSvc),
		hostPortSnapshotHandler:      snapshothandler.NewHostPortSnapshotHandler(hostPortSnapshotSvc),
		screenshotSnapshotHandler:    snapshothandler.NewScreenshotSnapshotHandler(screenshotSnapshotSvc),
		vulnerabilitySnapshotHandler: snapshothandler.NewVulnerabilitySnapshotHandler(vulnerabilitySnapshotSvc),
		presetHandler:                cataloghandler.NewPresetHandler(presetSvc),

		agentRepo:    agentRepo,
		scanTaskRepo: scanTaskRepo,
	}
}
