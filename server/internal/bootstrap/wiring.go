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
	agenthandler "github.com/yyhuni/lunafox/server/internal/modules/agent/handler"
	agentrepo "github.com/yyhuni/lunafox/server/internal/modules/agent/repository"
	assetservice "github.com/yyhuni/lunafox/server/internal/modules/asset/application"
	assethandler "github.com/yyhuni/lunafox/server/internal/modules/asset/handler"
	endpointhandler "github.com/yyhuni/lunafox/server/internal/modules/asset/handler/endpoint"
	websitehandler "github.com/yyhuni/lunafox/server/internal/modules/asset/handler/website"
	assetrepo "github.com/yyhuni/lunafox/server/internal/modules/asset/repository"
	catalogservice "github.com/yyhuni/lunafox/server/internal/modules/catalog/application"
	cataloghandler "github.com/yyhuni/lunafox/server/internal/modules/catalog/handler"
	catalogrepo "github.com/yyhuni/lunafox/server/internal/modules/catalog/repository"
	identityservice "github.com/yyhuni/lunafox/server/internal/modules/identity/application"
	identityhandler "github.com/yyhuni/lunafox/server/internal/modules/identity/handler"
	identityrepo "github.com/yyhuni/lunafox/server/internal/modules/identity/repository"
	scanservice "github.com/yyhuni/lunafox/server/internal/modules/scan/application"
	scanhandler "github.com/yyhuni/lunafox/server/internal/modules/scan/handler"
	scanrepo "github.com/yyhuni/lunafox/server/internal/modules/scan/repository"
	securityservice "github.com/yyhuni/lunafox/server/internal/modules/security/application"
	securityhandler "github.com/yyhuni/lunafox/server/internal/modules/security/handler"
	securityrepo "github.com/yyhuni/lunafox/server/internal/modules/security/repository"
	snapshotservice "github.com/yyhuni/lunafox/server/internal/modules/snapshot/application"
	snapshothandler "github.com/yyhuni/lunafox/server/internal/modules/snapshot/handler"
	snapshotrepo "github.com/yyhuni/lunafox/server/internal/modules/snapshot/repository"
	"github.com/yyhuni/lunafox/server/internal/preset"
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
	subdomainHandler     *assethandler.SubdomainHandler
	endpointHandler      *endpointhandler.EndpointHandler
	directoryHandler     *assethandler.DirectoryHandler
	hostPortHandler      *assethandler.HostPortHandler
	screenshotHandler    *assethandler.ScreenshotHandler
	vulnerabilityHandler *securityhandler.VulnerabilityHandler
	scanHandler          *scanhandler.ScanHandler
	scanLogHandler       *scanhandler.ScanLogHandler
	workerHandler        *cataloghandler.WorkerHandler

	agentHandler     *agenthandler.AgentHandler
	agentWSHandler   *agenthandler.AgentWebSocketHandler
	agentTaskHandler *scanhandler.AgentTaskHandler

	websiteSnapshotHandler       *snapshothandler.WebsiteSnapshotHandler
	subdomainSnapshotHandler     *snapshothandler.SubdomainSnapshotHandler
	endpointSnapshotHandler      *snapshothandler.EndpointSnapshotHandler
	directorySnapshotHandler     *snapshothandler.DirectorySnapshotHandler
	hostPortSnapshotHandler      *snapshothandler.HostPortSnapshotHandler
	screenshotSnapshotHandler    *snapshothandler.ScreenshotSnapshotHandler
	vulnerabilitySnapshotHandler *snapshothandler.VulnerabilitySnapshotHandler
	presetHandler                *cataloghandler.PresetHandler

	agentRepo    agentrepo.AgentRepository
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

	identityUserStore := identitywiring.NewUserStoreAdapter(userRepo)
	identityOrgStore := identitywiring.NewOrganizationStoreAdapter(orgRepo)
	userSvc := identityservice.NewUserService(identityUserStore)
	authSvc := identityservice.NewAuthService(identityUserStore, infra.jwtManager)
	orgSvc := identityservice.NewOrganizationService(identityOrgStore)
	targetSvc := catalogservice.NewTargetService(catalogwiring.NewTargetStoreAdapter(targetRepo), catalogwiring.NewOrganizationStoreAdapter(orgRepo))
	engineSvc := catalogservice.NewEngineService(catalogwiring.NewEngineStoreAdapter(engineRepo))
	wordlistSvc := catalogservice.NewWordlistService(catalogwiring.NewWordlistStoreAdapter(wordlistRepo), cfg.Storage.WordlistsBasePath)
	assetTargetLookup := assetwiring.NewTargetLookupAdapter(targetRepo)
	websiteSvc := assetservice.NewWebsiteService(assetwiring.NewWebsiteStoreAdapter(websiteRepo), assetTargetLookup)
	subdomainSvc := assetservice.NewSubdomainService(assetwiring.NewSubdomainStoreAdapter(subdomainRepo), assetTargetLookup)
	endpointSvc := assetservice.NewEndpointService(assetwiring.NewEndpointStoreAdapter(endpointRepo), assetTargetLookup)
	directorySvc := assetservice.NewDirectoryService(assetwiring.NewDirectoryStoreAdapter(directoryRepo), assetTargetLookup)
	hostPortSvc := assetservice.NewHostPortService(assetwiring.NewHostPortStoreAdapter(hostPortRepo), assetTargetLookup)
	screenshotSvc := assetservice.NewScreenshotService(assetwiring.NewScreenshotStoreAdapter(screenshotRepo), assetTargetLookup)
	vulnerabilitySvc := securityservice.NewVulnerabilityService(vulnerabilityRepo, securitywiring.NewTargetLookupAdapter(targetRepo))
	scanSvc := scanservice.NewScanService(
		scanwiring.NewStoreAdapter(scanRepo),
		scanwiring.NewCommandStore(scanRepo),
		scanwiring.NewTaskCancellerAdapter(scanTaskRepo),
		infra.wsHub,
		scanwiring.NewCreateTargetLookupAdapter(targetRepo),
	)
	scanLogSvc := scanlogwiring.NewApplicationService(scanLogRepo, scanRepo)
	workerSvc := workerwiring.NewApplicationService(scanRepo, subfinderProviderSettingsRepo)
	scanTaskSvc := scanservice.NewScanTaskService(scanwiring.NewTaskStoreAdapter(scanTaskRepo), scanwiring.NewTaskRuntimeScanStoreAdapter(scanRepo))
	agentSvc := agentservice.NewAgentService(agentRepo, registrationTokenRepo)

	snapshotScanLookup := snapshotwiring.NewScanLookupAdapter(scanRepo)

	websiteSnapshotStore := snapshotwiring.NewWebsiteStoreAdapter(websiteSnapshotRepo)
	websiteSnapshotSvc := snapshotservice.NewWebsiteSnapshotService(
		snapshotservice.NewWebsiteSnapshotQueryService(websiteSnapshotStore, snapshotScanLookup),
		snapshotservice.NewWebsiteSnapshotCommandService(websiteSnapshotStore, snapshotScanLookup, snapshotwiring.NewWebsiteAssetSyncAdapter(websiteSvc)),
	)

	subdomainSnapshotStore := snapshotwiring.NewSubdomainStoreAdapter(subdomainSnapshotRepo)
	subdomainSnapshotSvc := snapshotservice.NewSubdomainSnapshotService(
		snapshotservice.NewSubdomainSnapshotQueryService(subdomainSnapshotStore, snapshotScanLookup),
		snapshotservice.NewSubdomainSnapshotCommandService(subdomainSnapshotStore, snapshotScanLookup, snapshotwiring.NewSubdomainAssetSyncAdapter(subdomainSvc)),
	)

	endpointSnapshotStore := snapshotwiring.NewEndpointStoreAdapter(endpointSnapshotRepo)
	endpointSnapshotSvc := snapshotservice.NewEndpointSnapshotService(
		snapshotservice.NewEndpointSnapshotQueryService(endpointSnapshotStore, snapshotScanLookup),
		snapshotservice.NewEndpointSnapshotCommandService(endpointSnapshotStore, snapshotScanLookup, snapshotwiring.NewEndpointAssetSyncAdapter(endpointSvc)),
	)

	directorySnapshotStore := snapshotwiring.NewDirectoryStoreAdapter(directorySnapshotRepo)
	directorySnapshotSvc := snapshotservice.NewDirectorySnapshotService(
		snapshotservice.NewDirectorySnapshotQueryService(directorySnapshotStore, snapshotScanLookup),
		snapshotservice.NewDirectorySnapshotCommandService(directorySnapshotStore, snapshotScanLookup, snapshotwiring.NewDirectoryAssetSyncAdapter(directorySvc)),
	)

	hostPortSnapshotStore := snapshotwiring.NewHostPortStoreAdapter(hostPortSnapshotRepo)
	hostPortSnapshotSvc := snapshotservice.NewHostPortSnapshotService(
		snapshotservice.NewHostPortSnapshotQueryService(hostPortSnapshotStore, snapshotScanLookup),
		snapshotservice.NewHostPortSnapshotCommandService(hostPortSnapshotStore, snapshotScanLookup, snapshotwiring.NewHostPortAssetSyncAdapter(hostPortSvc)),
	)

	screenshotSnapshotStore := snapshotwiring.NewScreenshotStoreAdapter(screenshotSnapshotRepo)
	screenshotSnapshotSvc := snapshotservice.NewScreenshotSnapshotService(
		snapshotservice.NewScreenshotSnapshotQueryService(screenshotSnapshotStore, snapshotScanLookup),
		snapshotservice.NewScreenshotSnapshotCommandService(screenshotSnapshotStore, snapshotScanLookup, snapshotwiring.NewScreenshotAssetSyncAdapter(screenshotSvc)),
	)

	vulnerabilitySnapshotStore := snapshotwiring.NewVulnerabilityStoreAdapter(vulnerabilitySnapshotRepo)
	vulnerabilitySnapshotSvc := snapshotservice.NewVulnerabilitySnapshotService(
		snapshotservice.NewVulnerabilitySnapshotQueryService(vulnerabilitySnapshotStore, snapshotScanLookup),
		snapshotservice.NewVulnerabilitySnapshotCommandService(vulnerabilitySnapshotStore, snapshotScanLookup, snapshotwiring.NewVulnerabilityAssetSyncAdapter(vulnerabilitySvc), snapshotwiring.NewVulnerabilityRawOutputCodec()),
	)
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
		subdomainHandler:     assethandler.NewSubdomainHandler(subdomainSvc),
		endpointHandler:      endpointhandler.NewEndpointHandler(endpointSvc),
		directoryHandler:     assethandler.NewDirectoryHandler(directorySvc),
		hostPortHandler:      assethandler.NewHostPortHandler(hostPortSvc),
		screenshotHandler:    assethandler.NewScreenshotHandler(screenshotSvc),
		vulnerabilityHandler: securityhandler.NewVulnerabilityHandler(vulnerabilitySvc),
		scanHandler:          scanhandler.NewScanHandler(scanSvc),
		scanLogHandler:       scanhandler.NewScanLogHandler(scanLogSvc),
		workerHandler:        cataloghandler.NewWorkerHandler(workerSvc),

		agentHandler: agenthandler.NewAgentHandler(
			agentSvc,
			cfg.PublicURL,
			infra.serverVersion,
			infra.agentImage,
			cfg.Worker.Token,
			infra.heartbeatCache,
			infra.wsHub,
		),
		agentWSHandler: agenthandler.NewAgentWebSocketHandler(
			infra.wsHub,
			agentRepo,
			infra.heartbeatCache,
			infra.serverVersion,
			infra.agentImage,
		),
		agentTaskHandler: scanhandler.NewAgentTaskHandler(scanTaskSvc),

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
