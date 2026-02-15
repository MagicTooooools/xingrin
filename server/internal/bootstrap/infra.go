package bootstrap

import (
	"context"
	"embed"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/yyhuni/lunafox/server/internal/auth"
	"github.com/yyhuni/lunafox/server/internal/cache"
	"github.com/yyhuni/lunafox/server/internal/config"
	"github.com/yyhuni/lunafox/server/internal/database"
	"github.com/yyhuni/lunafox/server/internal/pkg"
	"github.com/yyhuni/lunafox/server/internal/pkg/imagecfg"
	pkgvalidator "github.com/yyhuni/lunafox/server/internal/pkg/validator"
	"github.com/yyhuni/lunafox/server/internal/preset"
	ws "github.com/yyhuni/lunafox/server/internal/websocket"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type infra struct {
	db             *gorm.DB
	redisClient    *redis.Client
	heartbeatCache cache.HeartbeatCache
	wsHub          *ws.Hub
	jwtManager     *auth.JWTManager
	presetLoader   *preset.Loader
	serverVersion  string
	agentImage     string
	workerImage    string
}

func initInfra(cfg *config.Config, migrationsFS embed.FS) *infra {
	serverVersion := strings.TrimSpace(os.Getenv("IMAGE_TAG"))
	if serverVersion == "" {
		pkg.Fatal("IMAGE_TAG environment variable is required")
	}
	agentImage := resolveAgentImage()
	workerImage := resolveWorkerImage(agentImage)

	if err := pkgvalidator.Init(); err != nil {
		pkg.Fatal("Failed to initialize validator", zap.Error(err))
	}
	pkg.Info("Validator initialized with custom translations")

	db, err := database.NewDatabase(&cfg.Database)
	if err != nil {
		pkg.Fatal("Failed to connect to database", zap.Error(err))
	}
	pkg.Info("Database connected",
		zap.String("host", cfg.Database.Host),
		zap.Int("port", cfg.Database.Port),
		zap.String("name", cfg.Database.Name),
	)

	database.MigrationsFS = migrationsFS
	database.MigrationsPath = "migrations"

	sqlDB, err := db.DB()
	if err != nil {
		pkg.Fatal("Failed to get underlying sql.DB", zap.Error(err))
	}
	if err := database.RunMigrations(sqlDB); err != nil {
		pkg.Fatal("Failed to run database migrations", zap.Error(err))
	}

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
			pkg.Warn("Failed to connect to Redis, continuing without Redis", zap.Error(err))
			if closeErr := redisClient.Close(); closeErr != nil {
				pkg.Warn("Failed to close Redis client after ping failure", zap.Error(closeErr))
			}
			redisClient = nil
		} else {
			pkg.Info("Redis connected", zap.String("addr", cfg.Redis.Addr()))
		}
	}

	var heartbeatCache cache.HeartbeatCache
	if redisClient != nil {
		heartbeatCache = cache.NewHeartbeatCache(redisClient)
	}

	wsHub := ws.NewHub()
	go wsHub.Run()

	presetLoader, err := preset.NewLoader()
	if err != nil {
		pkg.Fatal("Failed to load preset engines", zap.Error(err))
	}
	pkg.Info("Preset engines loaded", zap.Int("count", len(presetLoader.List())))

	jwtManager := auth.NewJWTManager(cfg.JWT.Secret, cfg.JWT.AccessExpire, cfg.JWT.RefreshExpire)
	gin.SetMode(cfg.Server.Mode)

	return &infra{
		db:             db,
		redisClient:    redisClient,
		heartbeatCache: heartbeatCache,
		wsHub:          wsHub,
		jwtManager:     jwtManager,
		presetLoader:   presetLoader,
		serverVersion:  serverVersion,
		agentImage:     agentImage,
		workerImage:    workerImage,
	}
}

func resolveAgentImage() string {
	if image := strings.TrimSpace(os.Getenv("AGENT_IMAGE")); image != "" {
		return image
	}
	return imagecfg.BuildAgentImage(os.Getenv("IMAGE_REGISTRY"), os.Getenv("IMAGE_NAMESPACE"))
}

func resolveWorkerImage(agentImage string) string {
	if image := strings.TrimSpace(os.Getenv("WORKER_IMAGE")); image != "" {
		return image
	}
	return imagecfg.FallbackWorkerImage(agentImage)
}
