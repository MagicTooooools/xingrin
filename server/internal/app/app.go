package app

import (
	"context"
	"embed"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/yyhuni/lunafox/server/internal/auth"
	"github.com/yyhuni/lunafox/server/internal/cache"
	"github.com/yyhuni/lunafox/server/internal/config"
	"github.com/yyhuni/lunafox/server/internal/database"
	"github.com/yyhuni/lunafox/server/internal/handler"
	"github.com/yyhuni/lunafox/server/internal/job"
	"github.com/yyhuni/lunafox/server/internal/middleware"
	"github.com/yyhuni/lunafox/server/internal/pkg"
	pkgvalidator "github.com/yyhuni/lunafox/server/internal/pkg/validator"
	"github.com/yyhuni/lunafox/server/internal/repository"
	"github.com/yyhuni/lunafox/server/internal/router"
	"github.com/yyhuni/lunafox/server/internal/service"
	ws "github.com/yyhuni/lunafox/server/internal/websocket"
	"go.uber.org/zap"
)

// Run wires dependencies and starts the HTTP server.
func Run(ctx context.Context, cfg *config.Config, migrationsFS embed.FS) {
	pkg.Info("Starting server",
		zap.Int("port", cfg.Server.Port),
		zap.String("mode", cfg.Server.Mode),
	)

	serverVersion := pkg.ReadVersion("VERSION")
	agentImage := "yyhuni/lunafox-agent"

	// Initialize custom validator with translations
	if err := pkgvalidator.Init(); err != nil {
		pkg.Fatal("Failed to initialize validator", zap.Error(err))
	}
	pkg.Info("Validator initialized with custom translations")

	// Initialize database
	db, err := database.NewDatabase(&cfg.Database)
	if err != nil {
		pkg.Fatal("Failed to connect to database", zap.Error(err))
	}
	pkg.Info("Database connected",
		zap.String("host", cfg.Database.Host),
		zap.Int("port", cfg.Database.Port),
		zap.String("name", cfg.Database.Name),
	)

	// Run SQL migrations
	database.MigrationsFS = migrationsFS
	database.MigrationsPath = "migrations"

	sqlDB, err := db.DB()
	if err != nil {
		pkg.Fatal("Failed to get underlying sql.DB", zap.Error(err))
	}

	if err := database.RunMigrations(sqlDB); err != nil {
		pkg.Fatal("Failed to run database migrations", zap.Error(err))
	}

	// Initialize Redis (optional)
	var redisClient *redis.Client
	if cfg.Redis.Host != "" {
		redisClient = redis.NewClient(&redis.Options{
			Addr:     cfg.Redis.Addr(),
			Password: cfg.Redis.Password,
			DB:       cfg.Redis.DB,
		})

		rcCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := redisClient.Ping(rcCtx).Err(); err != nil {
			pkg.Warn("Failed to connect to Redis, continuing without Redis",
				zap.Error(err),
			)
			redisClient = nil
		} else {
			pkg.Info("Redis connected",
				zap.String("addr", cfg.Redis.Addr()),
			)
		}
	}
	var heartbeatCache cache.HeartbeatCache
	if redisClient != nil {
		heartbeatCache = cache.NewHeartbeatCache(redisClient)
	}

	wsHub := ws.NewHub()
	go wsHub.Run()

	// Initialize JWT manager
	jwtManager := auth.NewJWTManager(
		cfg.JWT.Secret,
		cfg.JWT.AccessExpire,
		cfg.JWT.RefreshExpire,
	)

	// Set Gin mode
	gin.SetMode(cfg.Server.Mode)

	// Create Gin router
	engine := gin.New()

	// Disable automatic redirect for trailing slash
	// This prevents 301/307 redirects when URL has/doesn't have trailing slash
	engine.RedirectTrailingSlash = false
	engine.RedirectFixedPath = false

	// Add middleware
	engine.Use(middleware.Recovery())
	engine.Use(middleware.Logger())

	// Create repositories
	userRepo := repository.NewUserRepository(db)
	orgRepo := repository.NewOrganizationRepository(db)
	targetRepo := repository.NewTargetRepository(db)
	engineRepo := repository.NewEngineRepository(db)
	wordlistRepo := repository.NewWordlistRepository(db)
	websiteRepo := repository.NewWebsiteRepository(db)
	subdomainRepo := repository.NewSubdomainRepository(db)
	endpointRepo := repository.NewEndpointRepository(db)
	directoryRepo := repository.NewDirectoryRepository(db)
	hostPortRepo := repository.NewHostPortRepository(db)
	screenshotRepo := repository.NewScreenshotRepository(db)
	vulnerabilityRepo := repository.NewVulnerabilityRepository(db)
	scanRepo := repository.NewScanRepository(db)
	scanLogRepo := repository.NewScanLogRepository(db)
	subfinderProviderSettingsRepo := repository.NewSubfinderProviderSettingsRepository(db)
	websiteSnapshotRepo := repository.NewWebsiteSnapshotRepository(db)
	subdomainSnapshotRepo := repository.NewSubdomainSnapshotRepository(db)
	endpointSnapshotRepo := repository.NewEndpointSnapshotRepository(db)
	directorySnapshotRepo := repository.NewDirectorySnapshotRepository(db)
	hostPortSnapshotRepo := repository.NewHostPortSnapshotRepository(db)
	screenshotSnapshotRepo := repository.NewScreenshotSnapshotRepository(db)
	vulnerabilitySnapshotRepo := repository.NewVulnerabilitySnapshotRepository(db)
	agentRepo := repository.NewAgentRepository(db)
	registrationTokenRepo := repository.NewRegistrationTokenRepository(db)
	scanTaskRepo := repository.NewScanTaskRepository(db)

	// Create services
	userSvc := service.NewUserService(userRepo)
	orgSvc := service.NewOrganizationService(orgRepo)
	targetSvc := service.NewTargetService(targetRepo, orgRepo)
	engineSvc := service.NewEngineService(engineRepo)
	wordlistSvc := service.NewWordlistService(wordlistRepo, cfg.Storage.WordlistsBasePath)
	websiteSvc := service.NewWebsiteService(websiteRepo, targetRepo)
	subdomainSvc := service.NewSubdomainService(subdomainRepo, targetRepo)
	endpointSvc := service.NewEndpointService(endpointRepo, targetRepo)
	directorySvc := service.NewDirectoryService(directoryRepo, targetRepo)
	hostPortSvc := service.NewHostPortService(hostPortRepo, targetRepo)
	screenshotSvc := service.NewScreenshotService(screenshotRepo, targetRepo)
	vulnerabilitySvc := service.NewVulnerabilityService(vulnerabilityRepo, targetRepo)
	scanSvc := service.NewScanService(scanRepo, scanLogRepo, scanTaskRepo, wsHub, targetRepo, orgRepo)
	scanLogSvc := service.NewScanLogService(scanLogRepo, scanRepo)
	workerSvc := service.NewWorkerService(scanRepo, subfinderProviderSettingsRepo)
	agentSvc := service.NewAgentService(agentRepo, registrationTokenRepo)
	scanTaskSvc := service.NewScanTaskService(scanTaskRepo, scanRepo)
	websiteSnapshotSvc := service.NewWebsiteSnapshotService(websiteSnapshotRepo, scanRepo, websiteSvc)
	subdomainSnapshotSvc := service.NewSubdomainSnapshotService(subdomainSnapshotRepo, scanRepo, subdomainSvc)
	endpointSnapshotSvc := service.NewEndpointSnapshotService(endpointSnapshotRepo, scanRepo, endpointSvc)
	directorySnapshotSvc := service.NewDirectorySnapshotService(directorySnapshotRepo, scanRepo, directorySvc)
	hostPortSnapshotSvc := service.NewHostPortSnapshotService(hostPortSnapshotRepo, scanRepo, hostPortSvc)
	screenshotSnapshotSvc := service.NewScreenshotSnapshotService(screenshotSnapshotRepo, scanRepo, screenshotSvc)
	vulnerabilitySnapshotSvc := service.NewVulnerabilitySnapshotService(vulnerabilitySnapshotRepo, scanRepo, vulnerabilitySvc)

	// Create handlers
	healthHandler := handler.NewHealthHandler(db, redisClient)
	authHandler := handler.NewAuthHandler(db, jwtManager)
	userHandler := handler.NewUserHandler(userSvc)
	orgHandler := handler.NewOrganizationHandler(orgSvc)
	targetHandler := handler.NewTargetHandler(targetSvc)
	engineHandler := handler.NewEngineHandler(engineSvc)
	wordlistHandler := handler.NewWordlistHandler(wordlistSvc)
	websiteHandler := handler.NewWebsiteHandler(websiteSvc)
	subdomainHandler := handler.NewSubdomainHandler(subdomainSvc)
	endpointHandler := handler.NewEndpointHandler(endpointSvc)
	directoryHandler := handler.NewDirectoryHandler(directorySvc)
	hostPortHandler := handler.NewHostPortHandler(hostPortSvc)
	screenshotHandler := handler.NewScreenshotHandler(screenshotSvc)
	vulnerabilityHandler := handler.NewVulnerabilityHandler(vulnerabilitySvc)
	scanHandler := handler.NewScanHandler(scanSvc)
	scanLogHandler := handler.NewScanLogHandler(scanLogSvc)
	workerHandler := handler.NewWorkerHandler(workerSvc)
	agentHandler := handler.NewAgentHandler(agentSvc, cfg.PublicURL, serverVersion, agentImage, cfg.Worker.Token, heartbeatCache, wsHub)
	agentWSHandler := handler.NewAgentWebSocketHandler(wsHub, agentRepo, heartbeatCache, serverVersion, agentImage)
	agentTaskHandler := handler.NewAgentTaskHandler(scanTaskSvc)
	websiteSnapshotHandler := handler.NewWebsiteSnapshotHandler(websiteSnapshotSvc)
	subdomainSnapshotHandler := handler.NewSubdomainSnapshotHandler(subdomainSnapshotSvc)
	endpointSnapshotHandler := handler.NewEndpointSnapshotHandler(endpointSnapshotSvc)
	directorySnapshotHandler := handler.NewDirectorySnapshotHandler(directorySnapshotSvc)
	hostPortSnapshotHandler := handler.NewHostPortSnapshotHandler(hostPortSnapshotSvc)
	screenshotSnapshotHandler := handler.NewScreenshotSnapshotHandler(screenshotSnapshotSvc)
	vulnerabilitySnapshotHandler := handler.NewVulnerabilitySnapshotHandler(vulnerabilitySnapshotSvc)

	jobCtx, jobCancel := context.WithCancel(context.Background())
	defer jobCancel()

	agentMonitor := job.NewAgentMonitor(agentRepo, scanTaskRepo, time.Minute, 120*time.Second)
	go agentMonitor.Run(jobCtx)

	router.RegisterHealthRoutes(engine, healthHandler)

	// API routes
	api := engine.Group("/api")
	{
		router.RegisterAuthRoutes(api, authHandler)
		router.RegisterPublicRoutes(api, screenshotHandler, screenshotSnapshotHandler)
		router.RegisterWorkerRoutes(api, cfg.Worker.Token, workerHandler, wordlistHandler, subdomainSnapshotHandler, websiteSnapshotHandler, endpointSnapshotHandler)

		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(jwtManager))

		router.RegisterAgentRoutes(api, protected, agentHandler, agentWSHandler, agentTaskHandler, agentRepo)
		router.RegisterAuthProtectedRoutes(protected, authHandler)
		router.RegisterUserRoutes(protected, userHandler)
		router.RegisterOrganizationRoutes(protected, orgHandler)
		router.RegisterTargetRoutes(protected, targetHandler)
		router.RegisterWebsiteRoutes(protected, websiteHandler)
		router.RegisterSubdomainRoutes(protected, subdomainHandler)
		router.RegisterEndpointRoutes(protected, endpointHandler, endpointSnapshotHandler)
		router.RegisterDirectoryRoutes(protected, directoryHandler)
		router.RegisterHostPortRoutes(protected, hostPortHandler, hostPortSnapshotHandler)
		router.RegisterScreenshotRoutes(protected, screenshotHandler, screenshotSnapshotHandler)
		router.RegisterVulnerabilityRoutes(protected, vulnerabilityHandler)
		router.RegisterEngineRoutes(protected, engineHandler)
		router.RegisterWordlistRoutes(protected, wordlistHandler)
		router.RegisterScanRoutes(protected, scanHandler)
		router.RegisterScanLogRoutes(protected, scanLogHandler)
		router.RegisterScanSnapshotRoutes(
			protected,
			websiteSnapshotHandler,
			subdomainSnapshotHandler,
			endpointSnapshotHandler,
			directorySnapshotHandler,
			hostPortSnapshotHandler,
			screenshotSnapshotHandler,
			vulnerabilitySnapshotHandler,
		)
	}

	// Create HTTP server with trailing slash normalization
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      middleware.NormalizeTrailingSlash(engine),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		pkg.Info("Server listening", zap.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			pkg.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal
	<-ctx.Done()

	pkg.Info("Shutting down server...")

	jobCancel()

	// Graceful shutdown with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		pkg.Error("Server forced to shutdown", zap.Error(err))
	}

	// Close database connection
	if sqlDB, err := db.DB(); err == nil {
		if err := sqlDB.Close(); err != nil {
			pkg.Error("Failed to close database connection", zap.Error(err))
		}
	}

	// Close Redis connection
	if redisClient != nil {
		if err := redisClient.Close(); err != nil {
			pkg.Error("Failed to close Redis connection", zap.Error(err))
		}
	}

	pkg.Info("Server exited")
}
