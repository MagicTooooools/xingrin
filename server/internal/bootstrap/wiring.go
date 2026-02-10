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

	agentRepo := agentrepo.NewAgentRepository(db)
	registrationTokenRepo := agentrepo.NewRegistrationTokenRepository(db)
	scanTaskRepo := scanrepo.NewScanTaskRepository(db)

	identityUserStore := identitywiring.NewIdentityUserStoreAdapter(userRepo)
	identityAuthUserStore := identitywiring.NewIdentityAuthUserStoreAdapter(userRepo)
	identityOrgStore := identitywiring.NewIdentityOrganizationStoreAdapter(orgRepo)
	userSvc := identityservice.NewUserFacade(identityUserStore)
	authSvc := identityservice.NewAuthFacade(identityAuthUserStore, infra.jwtManager)
	orgSvc := identityservice.NewOrganizationFacade(identityOrgStore)

	catalogTargetStore := catalogwiring.NewCatalogTargetStoreAdapter(targetRepo)
	catalogOrganizationStore := catalogwiring.NewCatalogOrganizationStoreAdapter(orgRepo)
	targetSvc := catalogservice.NewTargetFacade(catalogTargetStore, catalogOrganizationStore)
	catalogEngineStore := catalogwiring.NewCatalogEngineStoreAdapter(engineRepo)
	engineSvc := catalogservice.NewEngineFacade(catalogEngineStore)
	catalogWordlistStore := catalogwiring.NewCatalogWordlistStoreAdapter(wordlistRepo)
	wordlistSvc := catalogservice.NewWordlistFacade(catalogWordlistStore, cfg.Storage.WordlistsBasePath)

	assetTargetLookup := assetwiring.NewAssetTargetLookupAdapter(targetRepo)
	assetWebsiteStore := assetwiring.NewAssetWebsiteStoreAdapter(websiteRepo)
	websiteSvc := assetservice.NewWebsiteFacade(assetWebsiteStore, assetTargetLookup)
	assetSubdomainStore := assetwiring.NewAssetSubdomainStoreAdapter(subdomainRepo)
	subdomainSvc := assetservice.NewSubdomainFacade(assetSubdomainStore, assetTargetLookup)
	assetEndpointStore := assetwiring.NewAssetEndpointStoreAdapter(endpointRepo)
	endpointSvc := assetservice.NewEndpointFacade(assetEndpointStore, assetTargetLookup)
	assetDirectoryStore := assetwiring.NewAssetDirectoryStoreAdapter(directoryRepo)
	directorySvc := assetservice.NewDirectoryFacade(assetDirectoryStore, assetTargetLookup)
	assetHostPortStore := assetwiring.NewAssetHostPortStoreAdapter(hostPortRepo)
	hostPortSvc := assetservice.NewHostPortFacade(assetHostPortStore, assetTargetLookup)
	assetScreenshotStore := assetwiring.NewAssetScreenshotStoreAdapter(screenshotRepo)
	screenshotSvc := assetservice.NewScreenshotFacade(assetScreenshotStore, assetTargetLookup)

	securityVulnerabilityStore := securitywiring.NewSecurityVulnerabilityStoreAdapter(vulnerabilityRepo)
	securityTargetLookup := securitywiring.NewSecurityTargetLookupAdapter(targetRepo)
	vulnerabilitySvc := securityservice.NewVulnerabilityFacade(securityVulnerabilityStore, securityTargetLookup)
	scanQueryStore := scanwiring.NewScanQueryStoreAdapter(scanRepo)
	scanCommandStore := scanwiring.NewScanCommandStoreAdapter(scanRepo)
	scanDomainRepository := scanwiring.NewScanDomainRepositoryAdapter(scanRepo)
	scanTaskCanceller := scanwiring.NewScanTaskCancellerAdapter(scanTaskRepo)
	scanTargetLookup := scanwiring.NewScanTargetLookupAdapter(targetRepo)
	scanSvc := scanwiring.NewScanApplicationService(
		scanQueryStore,
		scanCommandStore,
		scanDomainRepository,
		scanTaskCanceller,
		infra.wsHub,
		scanTargetLookup,
	)
	scanLogStore := scanlogwiring.NewScanLogStoreAdapter(scanLogRepo)
	scanLogLookup := scanlogwiring.NewScanLogScanLookupAdapter(scanRepo)
	scanLogSvc := scanlogwiring.NewScanLogApplicationService(scanLogStore, scanLogLookup)
	workerScanGuard := workerwiring.NewWorkerScanGuardAdapter(scanRepo)
	workerSettingsStore := workerwiring.NewWorkerProviderSettingsStoreAdapter(subfinderProviderSettingsRepo)
	workerSvc := workerwiring.NewWorkerApplicationService(workerScanGuard, workerSettingsStore)
	scanTaskStore := scanwiring.NewScanTaskStoreAdapter(scanTaskRepo)
	scanTaskRuntimeStore := scanwiring.NewScanTaskRuntimeStoreAdapter(scanRepo)
	scanTaskSvc := scanwiring.NewScanTaskApplicationService(scanTaskStore, scanTaskRuntimeStore)
	agentSvc := agentservice.NewAgentFacade(agentRepo, registrationTokenRepo)
	agentRuntimeSvc := agentservice.NewAgentRuntimeService(agentRepo, infra.heartbeatCache, ws.NewAgentMessagePublisher(infra.wsHub), infra.serverVersion, infra.agentImage)
	agentTaskSvc := agentservice.NewAgentTaskService(scanTaskSvc)

	snapshotScanLookup := snapshotwiring.NewSnapshotScanRefLookupAdapter(scanRepo)

	websiteSnapshotStore := snapshotwiring.NewSnapshotWebsiteStoreAdapter(websiteSnapshotRepo)
	websiteAssetSync := snapshotwiring.NewSnapshotWebsiteAssetSyncAdapter(websiteSvc)
	websiteSnapshotSvc := snapshotwiring.NewSnapshotWebsiteApplicationService(websiteSnapshotStore, snapshotScanLookup, websiteAssetSync)

	subdomainSnapshotStore := snapshotwiring.NewSnapshotSubdomainStoreAdapter(subdomainSnapshotRepo)
	subdomainAssetSync := snapshotwiring.NewSnapshotSubdomainAssetSyncAdapter(subdomainSvc)
	subdomainSnapshotSvc := snapshotwiring.NewSnapshotSubdomainApplicationService(subdomainSnapshotStore, snapshotScanLookup, subdomainAssetSync)

	endpointSnapshotStore := snapshotwiring.NewSnapshotEndpointStoreAdapter(endpointSnapshotRepo)
	endpointAssetSync := snapshotwiring.NewSnapshotEndpointAssetSyncAdapter(endpointSvc)
	endpointSnapshotSvc := snapshotwiring.NewSnapshotEndpointApplicationService(endpointSnapshotStore, snapshotScanLookup, endpointAssetSync)

	directorySnapshotStore := snapshotwiring.NewSnapshotDirectoryStoreAdapter(directorySnapshotRepo)
	directoryAssetSync := snapshotwiring.NewSnapshotDirectoryAssetSyncAdapter(directorySvc)
	directorySnapshotSvc := snapshotwiring.NewSnapshotDirectoryApplicationService(directorySnapshotStore, snapshotScanLookup, directoryAssetSync)

	hostPortSnapshotStore := snapshotwiring.NewSnapshotHostPortStoreAdapter(hostPortSnapshotRepo)
	hostPortAssetSync := snapshotwiring.NewSnapshotHostPortAssetSyncAdapter(hostPortSvc)
	hostPortSnapshotSvc := snapshotwiring.NewSnapshotHostPortApplicationService(hostPortSnapshotStore, snapshotScanLookup, hostPortAssetSync)

	screenshotSnapshotStore := snapshotwiring.NewSnapshotScreenshotStoreAdapter(screenshotSnapshotRepo)
	screenshotAssetSync := snapshotwiring.NewSnapshotScreenshotAssetSyncAdapter(screenshotSvc)
	screenshotSnapshotSvc := snapshotwiring.NewSnapshotScreenshotApplicationService(screenshotSnapshotStore, snapshotScanLookup, screenshotAssetSync)

	vulnerabilitySnapshotStore := snapshotwiring.NewSnapshotVulnerabilityStoreAdapter(vulnerabilitySnapshotRepo)
	vulnerabilityAssetSync := snapshotwiring.NewSnapshotVulnerabilityAssetSyncAdapter(vulnerabilitySvc)
	vulnerabilityRawOutputCodec := snapshotwiring.NewSnapshotVulnerabilityRawOutputCodec()
	vulnerabilitySnapshotSvc := snapshotwiring.NewSnapshotVulnerabilityApplicationService(vulnerabilitySnapshotStore, snapshotScanLookup, vulnerabilityAssetSync, vulnerabilityRawOutputCodec)
	presetSvc := preset.NewService(infra.presetLoader)

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
