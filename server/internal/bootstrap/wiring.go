package bootstrap

import (
	"github.com/yyhuni/lunafox/server/internal/config"
	agenthandler "github.com/yyhuni/lunafox/server/internal/modules/agent/handler"
	agentrepo "github.com/yyhuni/lunafox/server/internal/modules/agent/repository"
	agentservice "github.com/yyhuni/lunafox/server/internal/modules/agent/service"
	assethandler "github.com/yyhuni/lunafox/server/internal/modules/asset/handler"
	endpointhandler "github.com/yyhuni/lunafox/server/internal/modules/asset/handler/endpoint"
	websitehandler "github.com/yyhuni/lunafox/server/internal/modules/asset/handler/website"
	assetrepo "github.com/yyhuni/lunafox/server/internal/modules/asset/repository"
	assetservice "github.com/yyhuni/lunafox/server/internal/modules/asset/service"
	cataloghandler "github.com/yyhuni/lunafox/server/internal/modules/catalog/handler"
	catalogrepo "github.com/yyhuni/lunafox/server/internal/modules/catalog/repository"
	catalogservice "github.com/yyhuni/lunafox/server/internal/modules/catalog/service"
	identityhandler "github.com/yyhuni/lunafox/server/internal/modules/identity/handler"
	identityrepo "github.com/yyhuni/lunafox/server/internal/modules/identity/repository"
	identityservice "github.com/yyhuni/lunafox/server/internal/modules/identity/service"
	scanhandler "github.com/yyhuni/lunafox/server/internal/modules/scan/handler"
	scanrepo "github.com/yyhuni/lunafox/server/internal/modules/scan/repository"
	scanservice "github.com/yyhuni/lunafox/server/internal/modules/scan/service"
	securityhandler "github.com/yyhuni/lunafox/server/internal/modules/security/handler"
	securityrepo "github.com/yyhuni/lunafox/server/internal/modules/security/repository"
	securityservice "github.com/yyhuni/lunafox/server/internal/modules/security/service"
	snapshothandler "github.com/yyhuni/lunafox/server/internal/modules/snapshot/handler"
	snapshotrepo "github.com/yyhuni/lunafox/server/internal/modules/snapshot/repository"
	snapshotservice "github.com/yyhuni/lunafox/server/internal/modules/snapshot/service"
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

	userSvc := identityservice.NewUserService(userRepo)
	authSvc := identityservice.NewAuthService(userRepo, infra.jwtManager)
	orgSvc := identityservice.NewOrganizationService(orgRepo)
	targetSvc := catalogservice.NewTargetService(targetRepo, orgRepo)
	engineSvc := catalogservice.NewEngineService(engineRepo)
	wordlistSvc := catalogservice.NewWordlistService(wordlistRepo, cfg.Storage.WordlistsBasePath)
	websiteSvc := assetservice.NewWebsiteService(websiteRepo, targetRepo)
	subdomainSvc := assetservice.NewSubdomainService(subdomainRepo, targetRepo)
	endpointSvc := assetservice.NewEndpointService(endpointRepo, targetRepo)
	directorySvc := assetservice.NewDirectoryService(directoryRepo, targetRepo)
	hostPortSvc := assetservice.NewHostPortService(hostPortRepo, targetRepo)
	screenshotSvc := assetservice.NewScreenshotService(screenshotRepo, targetRepo)
	vulnerabilitySvc := securityservice.NewVulnerabilityService(vulnerabilityRepo, targetRepo)
	scanSvc := scanservice.NewScanService(scanRepo, scanTaskRepo, infra.wsHub, targetRepo, orgRepo)
	scanLogSvc := scanservice.NewScanLogService(scanLogRepo, scanRepo)
	workerSvc := catalogservice.NewWorkerService(scanRepo, subfinderProviderSettingsRepo)
	scanTaskSvc := scanservice.NewScanTaskService(scanTaskRepo, scanRepo)
	agentSvc := agentservice.NewAgentService(agentRepo, registrationTokenRepo)

	websiteSnapshotSvc := snapshotservice.NewWebsiteSnapshotService(websiteSnapshotRepo, scanRepo, websiteSvc)
	subdomainSnapshotSvc := snapshotservice.NewSubdomainSnapshotService(subdomainSnapshotRepo, scanRepo, subdomainSvc)
	endpointSnapshotSvc := snapshotservice.NewEndpointSnapshotService(endpointSnapshotRepo, scanRepo, endpointSvc)
	directorySnapshotSvc := snapshotservice.NewDirectorySnapshotService(directorySnapshotRepo, scanRepo, directorySvc)
	hostPortSnapshotSvc := snapshotservice.NewHostPortSnapshotService(hostPortSnapshotRepo, scanRepo, hostPortSvc)
	screenshotSnapshotSvc := snapshotservice.NewScreenshotSnapshotService(screenshotSnapshotRepo, scanRepo, screenshotSvc)
	vulnerabilitySnapshotSvc := snapshotservice.NewVulnerabilitySnapshotService(vulnerabilitySnapshotRepo, scanRepo, vulnerabilitySvc)
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
