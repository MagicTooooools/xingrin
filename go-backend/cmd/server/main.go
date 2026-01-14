package main

import (
	"context"
	"embed"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/xingrin/go-backend/internal/auth"
	"github.com/xingrin/go-backend/internal/config"
	"github.com/xingrin/go-backend/internal/database"
	"github.com/xingrin/go-backend/internal/handler"
	"github.com/xingrin/go-backend/internal/middleware"
	"github.com/xingrin/go-backend/internal/pkg"
	pkgvalidator "github.com/xingrin/go-backend/internal/pkg/validator"
	"github.com/xingrin/go-backend/internal/repository"
	"github.com/xingrin/go-backend/internal/service"
	"go.uber.org/zap"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	if err := pkg.InitLogger(&pkg.LogConfig{
		Level:  cfg.Log.Level,
		Format: cfg.Log.Format,
	}); err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer pkg.Sync()

	pkg.Info("Starting server",
		zap.Int("port", cfg.Server.Port),
		zap.String("mode", cfg.Server.Mode),
	)

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

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := redisClient.Ping(ctx).Err(); err != nil {
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

	// Initialize JWT manager
	jwtManager := auth.NewJWTManager(
		cfg.JWT.Secret,
		cfg.JWT.AccessExpire,
		cfg.JWT.RefreshExpire,
	)

	// Set Gin mode
	gin.SetMode(cfg.Server.Mode)

	// Create Gin router
	router := gin.New()

	// Disable automatic redirect for trailing slash
	// This prevents 301/307 redirects when URL has/doesn't have trailing slash
	router.RedirectTrailingSlash = false
	router.RedirectFixedPath = false

	// Add middleware
	router.Use(middleware.Recovery())
	router.Use(middleware.Logger())

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
	ipAddressRepo := repository.NewIPAddressRepository(db)
	screenshotRepo := repository.NewScreenshotRepository(db)
	vulnerabilityRepo := repository.NewVulnerabilityRepository(db)
	scanRepo := repository.NewScanRepository(db)
	scanLogRepo := repository.NewScanLogRepository(db)

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
	ipAddressSvc := service.NewIPAddressService(ipAddressRepo, targetRepo)
	screenshotSvc := service.NewScreenshotService(screenshotRepo, targetRepo)
	vulnerabilitySvc := service.NewVulnerabilityService(vulnerabilityRepo, targetRepo)
	scanSvc := service.NewScanService(scanRepo, scanLogRepo, targetRepo, orgRepo)
	scanLogSvc := service.NewScanLogService(scanLogRepo, scanRepo)

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
	ipAddressHandler := handler.NewIPAddressHandler(ipAddressSvc)
	screenshotHandler := handler.NewScreenshotHandler(screenshotSvc)
	vulnerabilityHandler := handler.NewVulnerabilityHandler(vulnerabilitySvc)
	scanHandler := handler.NewScanHandler(scanSvc)
	scanLogHandler := handler.NewScanLogHandler(scanLogSvc)

	// Register health routes
	router.GET("/health", healthHandler.Check)
	router.GET("/health/live", healthHandler.Liveness)
	router.GET("/health/ready", healthHandler.Readiness)

	// API routes
	api := router.Group("/api")
	{
		// Auth routes (public)
		authGroup := api.Group("/auth")
		{
			authGroup.POST("/login", authHandler.Login)
			authGroup.POST("/refresh", authHandler.RefreshToken)
		}

		// Protected routes
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(jwtManager))
		{
			// Auth
			protected.GET("/auth/me", authHandler.GetCurrentUser)

			// Users
			protected.POST("/users", userHandler.Create)
			protected.GET("/users", userHandler.List)
			protected.PUT("/users/password", userHandler.UpdatePassword)

			// Organizations
			protected.POST("/organizations", orgHandler.Create)
			protected.POST("/organizations/bulk-delete", orgHandler.BulkDelete)
			protected.GET("/organizations", orgHandler.List)
			protected.GET("/organizations/:id", orgHandler.GetByID)
			protected.GET("/organizations/:id/targets", orgHandler.ListTargets)
			protected.POST("/organizations/:id/link_targets", orgHandler.LinkTargets)
			protected.POST("/organizations/:id/unlink_targets", orgHandler.UnlinkTargets)
			protected.PUT("/organizations/:id", orgHandler.Update)
			protected.DELETE("/organizations/:id", orgHandler.Delete)

			// Targets
			protected.POST("/targets", targetHandler.Create)
			protected.POST("/targets/batch_create", targetHandler.BatchCreate)
			protected.POST("/targets/bulk-delete", targetHandler.BulkDelete)
			protected.GET("/targets", targetHandler.List)
			protected.GET("/targets/:id", targetHandler.GetByID)
			protected.PUT("/targets/:id", targetHandler.Update)
			protected.DELETE("/targets/:id", targetHandler.Delete)

			// Websites (nested under targets)
			protected.GET("/targets/:id/websites", websiteHandler.List)
			protected.GET("/targets/:id/websites/export", websiteHandler.Export)
			protected.POST("/targets/:id/websites/bulk-create", websiteHandler.BulkCreate)
			protected.POST("/targets/:id/websites/bulk-upsert", websiteHandler.BulkUpsert)

			// Websites (standalone)
			protected.GET("/websites/:id", websiteHandler.GetByID)
			protected.DELETE("/websites/:id", websiteHandler.Delete)
			protected.POST("/websites/bulk-delete", websiteHandler.BulkDelete)

			// Subdomains (nested under targets)
			protected.GET("/targets/:id/subdomains", subdomainHandler.List)
			protected.GET("/targets/:id/subdomains/export", subdomainHandler.Export)
			protected.POST("/targets/:id/subdomains/bulk-create", subdomainHandler.BulkCreate)

			// Subdomains (standalone)
			protected.POST("/subdomains/bulk-delete", subdomainHandler.BulkDelete)

			// Endpoints (nested under targets)
			protected.GET("/targets/:id/endpoints", endpointHandler.List)
			protected.GET("/targets/:id/endpoints/export", endpointHandler.Export)
			protected.POST("/targets/:id/endpoints/bulk-create", endpointHandler.BulkCreate)
			protected.POST("/targets/:id/endpoints/bulk-upsert", endpointHandler.BulkUpsert)

			// Endpoints (standalone)
			protected.GET("/endpoints/:id", endpointHandler.GetByID)
			protected.DELETE("/endpoints/:id", endpointHandler.Delete)
			protected.POST("/endpoints/bulk-delete", endpointHandler.BulkDelete)

			// Directories (nested under targets)
			protected.GET("/targets/:id/directories", directoryHandler.List)
			protected.GET("/targets/:id/directories/export", directoryHandler.Export)
			protected.POST("/targets/:id/directories/bulk-create", directoryHandler.BulkCreate)
			protected.POST("/targets/:id/directories/bulk-upsert", directoryHandler.BulkUpsert)

			// Directories (standalone)
			protected.POST("/directories/bulk-delete", directoryHandler.BulkDelete)

			// IP Addresses (nested under targets)
			protected.GET("/targets/:id/ip-addresses", ipAddressHandler.List)
			protected.GET("/targets/:id/ip-addresses/export", ipAddressHandler.Export)
			protected.POST("/targets/:id/ip-addresses/bulk-upsert", ipAddressHandler.BulkUpsert)

			// IP Addresses (standalone)
			protected.POST("/ip-addresses/bulk-delete", ipAddressHandler.BulkDelete)

			// Screenshots (nested under targets)
			protected.GET("/targets/:id/screenshots", screenshotHandler.ListByTargetID)
			protected.POST("/targets/:id/screenshots/bulk-upsert", screenshotHandler.BulkUpsert)

			// Screenshots (standalone)
			protected.GET("/screenshots/:id/image", screenshotHandler.GetImage)
			protected.POST("/screenshots/bulk-delete", screenshotHandler.BulkDelete)

			// Vulnerabilities (global)
			protected.GET("/assets/vulnerabilities", vulnerabilityHandler.ListAll)
			protected.GET("/assets/vulnerabilities/:id", vulnerabilityHandler.GetByID)

			// Vulnerabilities (nested under targets)
			protected.GET("/targets/:id/vulnerabilities", vulnerabilityHandler.ListByTarget)
			protected.POST("/targets/:id/vulnerabilities/bulk-create", vulnerabilityHandler.BulkCreate)

			// Vulnerabilities (standalone)
			protected.POST("/vulnerabilities/bulk-delete", vulnerabilityHandler.BulkDelete)
			protected.PATCH("/vulnerabilities/:id/review", vulnerabilityHandler.MarkAsReviewed)
			protected.PATCH("/vulnerabilities/:id/unreview", vulnerabilityHandler.MarkAsUnreviewed)
			protected.POST("/vulnerabilities/bulk-review", vulnerabilityHandler.BulkMarkAsReviewed)
			protected.POST("/vulnerabilities/bulk-unreview", vulnerabilityHandler.BulkMarkAsUnreviewed)

			// Engines
			protected.POST("/engines", engineHandler.Create)
			protected.GET("/engines", engineHandler.List)
			protected.GET("/engines/:id", engineHandler.GetByID)
			protected.PUT("/engines/:id", engineHandler.Update)
			protected.PATCH("/engines/:id", engineHandler.Patch)
			protected.DELETE("/engines/:id", engineHandler.Delete)

			// Wordlists
			protected.POST("/wordlists", wordlistHandler.Create)
			protected.GET("/wordlists", wordlistHandler.List)
			protected.DELETE("/wordlists/:id", wordlistHandler.Delete)
			protected.GET("/wordlists/download", wordlistHandler.Download)
			protected.GET("/wordlists/:id/content", wordlistHandler.GetContent)
			protected.PUT("/wordlists/:id/content", wordlistHandler.UpdateContent)

			// Scans
			protected.GET("/scans", scanHandler.List)
			protected.GET("/scans/statistics", scanHandler.Statistics)
			protected.GET("/scans/:id", scanHandler.GetByID)
			protected.DELETE("/scans/:id", scanHandler.Delete)
			protected.POST("/scans/:id/stop", scanHandler.Stop)
			protected.POST("/scans/initiate", scanHandler.Initiate)
			protected.POST("/scans/quick", scanHandler.Quick)
			protected.POST("/scans/bulk-delete", scanHandler.BulkDelete)

			// Scan Logs (nested under scans)
			protected.GET("/scans/:id/logs", scanLogHandler.List)
			protected.POST("/scans/:id/logs", scanLogHandler.BulkCreate)
		}
	}

	// Create HTTP server with trailing slash normalization
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      middleware.NormalizeTrailingSlash(router), // Wrap router to strip trailing slashes
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
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	pkg.Info("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
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
