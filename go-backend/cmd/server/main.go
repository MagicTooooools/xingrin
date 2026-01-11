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
	"github.com/xingrin/go-backend/internal/config"
	"github.com/xingrin/go-backend/internal/database"
	"github.com/xingrin/go-backend/internal/handler"
	"github.com/xingrin/go-backend/internal/middleware"
	"github.com/xingrin/go-backend/internal/pkg"
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


	// Set Gin mode
	gin.SetMode(cfg.Server.Mode)

	// Create Gin router
	router := gin.New()

	// Add middleware
	router.Use(middleware.Recovery())
	router.Use(middleware.Logger())

	// Create handlers
	healthHandler := handler.NewHealthHandler(db, redisClient)

	// Register routes
	router.GET("/health", healthHandler.Check)
	router.GET("/health/live", healthHandler.Liveness)
	router.GET("/health/ready", healthHandler.Readiness)

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
