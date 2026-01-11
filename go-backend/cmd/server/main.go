package main

import (
	"context"
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
	"github.com/xingrin/go-backend/internal/model"
	"github.com/xingrin/go-backend/internal/pkg"
	"github.com/xingrin/go-backend/internal/repository"
	"github.com/xingrin/go-backend/internal/service"
	"go.uber.org/zap"
)

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

	// Auto migrate database schema
	pkg.Info("Running database migrations...")
	if err := db.AutoMigrate(
		// Core models (no dependencies)
		&model.Organization{},
		&model.User{},
		&model.Target{},
		&model.ScanEngine{},
		&model.WorkerNode{},
		&model.Wordlist{},
		&model.NucleiTemplateRepo{},
		&model.BlacklistRule{},
		&model.NotificationSettings{},

		// Scan related (depends on Target, ScanEngine)
		&model.Scan{},
		&model.ScanInputTarget{},
		&model.ScanLog{},
		&model.ScheduledScan{},

		// Asset models (depends on Target, Scan)
		&model.Subdomain{},
		&model.HostPortMapping{},
		&model.Website{},
		&model.Endpoint{},
		&model.Directory{},
		&model.Screenshot{},
		&model.Vulnerability{},

		// Snapshot models (depends on Scan)
		&model.SubdomainSnapshot{},
		&model.HostPortMappingSnapshot{},
		&model.WebsiteSnapshot{},
		&model.EndpointSnapshot{},
		&model.DirectorySnapshot{},
		&model.ScreenshotSnapshot{},
		&model.VulnerabilitySnapshot{},

		// Statistics
		&model.AssetStatistics{},
		&model.StatisticsHistory{},
		&model.Notification{},

		// Auth
		&model.Session{},
		&model.SubfinderProviderSettings{},
	); err != nil {
		pkg.Fatal("Failed to migrate database", zap.Error(err))
	}
	pkg.Info("Database migrations completed")

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

	// Add middleware
	router.Use(middleware.Recovery())
	router.Use(middleware.Logger())

	// Create repositories
	userRepo := repository.NewUserRepository(db)
	orgRepo := repository.NewOrganizationRepository(db)
	targetRepo := repository.NewTargetRepository(db)
	engineRepo := repository.NewEngineRepository(db)

	// Create services
	userSvc := service.NewUserService(userRepo)
	orgSvc := service.NewOrganizationService(orgRepo)
	targetSvc := service.NewTargetService(targetRepo)
	engineSvc := service.NewEngineService(engineRepo)

	// Create handlers
	healthHandler := handler.NewHealthHandler(db, redisClient)
	authHandler := handler.NewAuthHandler(db, jwtManager)
	userHandler := handler.NewUserHandler(userSvc)
	orgHandler := handler.NewOrganizationHandler(orgSvc)
	targetHandler := handler.NewTargetHandler(targetSvc)
	engineHandler := handler.NewEngineHandler(engineSvc)

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
			protected.GET("/organizations", orgHandler.List)
			protected.GET("/organizations/:id", orgHandler.GetByID)
			protected.PUT("/organizations/:id", orgHandler.Update)
			protected.DELETE("/organizations/:id", orgHandler.Delete)

			// Targets
			protected.POST("/targets", targetHandler.Create)
			protected.GET("/targets", targetHandler.List)
			protected.GET("/targets/:id", targetHandler.GetByID)
			protected.PUT("/targets/:id", targetHandler.Update)
			protected.DELETE("/targets/:id", targetHandler.Delete)

			// Engines
			protected.POST("/engines", engineHandler.Create)
			protected.GET("/engines", engineHandler.List)
			protected.GET("/engines/:id", engineHandler.GetByID)
			protected.PUT("/engines/:id", engineHandler.Update)
			protected.DELETE("/engines/:id", engineHandler.Delete)
		}
	}

	// Create HTTP server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      router,
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
		sqlDB.Close()
	}

	// Close Redis connection
	if redisClient != nil {
		redisClient.Close()
	}

	pkg.Info("Server exited")
}
