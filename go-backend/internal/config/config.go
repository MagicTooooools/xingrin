package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	Log      LogConfig
	JWT      JWTConfig
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port int    `mapstructure:"SERVER_PORT"`
	Mode string `mapstructure:"GIN_MODE"`
}

// DatabaseConfig holds database-related configuration
type DatabaseConfig struct {
	Host            string `mapstructure:"DB_HOST"`
	Port            int    `mapstructure:"DB_PORT"`
	User            string `mapstructure:"DB_USER"`
	Password        string `mapstructure:"DB_PASSWORD"`
	Name            string `mapstructure:"DB_NAME"`
	SSLMode         string `mapstructure:"DB_SSLMODE"`
	MaxOpenConns    int    `mapstructure:"DB_MAX_OPEN_CONNS"`
	MaxIdleConns    int    `mapstructure:"DB_MAX_IDLE_CONNS"`
	ConnMaxLifetime int    `mapstructure:"DB_CONN_MAX_LIFETIME"`
}

// RedisConfig holds Redis-related configuration
type RedisConfig struct {
	Host     string `mapstructure:"REDIS_HOST"`
	Port     int    `mapstructure:"REDIS_PORT"`
	Password string `mapstructure:"REDIS_PASSWORD"`
	DB       int    `mapstructure:"REDIS_DB"`
}

// LogConfig holds logging-related configuration
type LogConfig struct {
	Level  string `mapstructure:"LOG_LEVEL"`
	Format string `mapstructure:"LOG_FORMAT"`
}

// JWTConfig holds JWT-related configuration
type JWTConfig struct {
	Secret        string        `mapstructure:"JWT_SECRET"`
	AccessExpire  time.Duration `mapstructure:"JWT_ACCESS_EXPIRE"`
	RefreshExpire time.Duration `mapstructure:"JWT_REFRESH_EXPIRE"`
}

// Load reads configuration from .env file and environment variables
// Priority: environment variables > .env file > defaults
func Load() (*Config, error) {
	v := viper.New()

	// Set default values
	setDefaults(v)

	// Try to read .env file (optional, won't fail if not found)
	v.SetConfigName(".env")
	v.SetConfigType("env")
	v.AddConfigPath(".")            // Current directory (go-backend/)
	v.AddConfigPath("./go-backend") // When running from project root
	if err := v.ReadInConfig(); err != nil {
		// .env file not found is OK, we'll use env vars or defaults
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// Only return error if it's not a "file not found" error
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	// Environment variables override .env file
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	cfg := &Config{}

	// Server config
	cfg.Server.Port = v.GetInt("SERVER_PORT")
	cfg.Server.Mode = v.GetString("GIN_MODE")

	// Database config
	cfg.Database.Host = v.GetString("DB_HOST")
	cfg.Database.Port = v.GetInt("DB_PORT")
	cfg.Database.User = v.GetString("DB_USER")
	cfg.Database.Password = v.GetString("DB_PASSWORD")
	cfg.Database.Name = v.GetString("DB_NAME")
	cfg.Database.SSLMode = v.GetString("DB_SSLMODE")
	cfg.Database.MaxOpenConns = v.GetInt("DB_MAX_OPEN_CONNS")
	cfg.Database.MaxIdleConns = v.GetInt("DB_MAX_IDLE_CONNS")
	cfg.Database.ConnMaxLifetime = v.GetInt("DB_CONN_MAX_LIFETIME")

	// Redis config
	cfg.Redis.Host = v.GetString("REDIS_HOST")
	cfg.Redis.Port = v.GetInt("REDIS_PORT")
	cfg.Redis.Password = v.GetString("REDIS_PASSWORD")
	cfg.Redis.DB = v.GetInt("REDIS_DB")

	// Log config
	cfg.Log.Level = v.GetString("LOG_LEVEL")
	cfg.Log.Format = v.GetString("LOG_FORMAT")

	// JWT config
	cfg.JWT.Secret = v.GetString("JWT_SECRET")
	cfg.JWT.AccessExpire = v.GetDuration("JWT_ACCESS_EXPIRE")
	cfg.JWT.RefreshExpire = v.GetDuration("JWT_REFRESH_EXPIRE")

	return cfg, nil
}

// setDefaults sets default values for configuration
func setDefaults(v *viper.Viper) {
	// Server defaults
	v.SetDefault("SERVER_PORT", 8888)
	v.SetDefault("GIN_MODE", "release")

	// Database defaults
	v.SetDefault("DB_HOST", "localhost")
	v.SetDefault("DB_PORT", 5432)
	v.SetDefault("DB_USER", "postgres")
	v.SetDefault("DB_PASSWORD", "")
	v.SetDefault("DB_NAME", "xingrin")
	v.SetDefault("DB_SSLMODE", "disable")
	v.SetDefault("DB_MAX_OPEN_CONNS", 25)
	v.SetDefault("DB_MAX_IDLE_CONNS", 5)
	v.SetDefault("DB_CONN_MAX_LIFETIME", 300)

	// Redis defaults
	v.SetDefault("REDIS_HOST", "localhost")
	v.SetDefault("REDIS_PORT", 6379)
	v.SetDefault("REDIS_PASSWORD", "")
	v.SetDefault("REDIS_DB", 0)

	// Log defaults
	v.SetDefault("LOG_LEVEL", "info")
	v.SetDefault("LOG_FORMAT", "json")

	// JWT defaults
	v.SetDefault("JWT_SECRET", "change-me-in-production-use-a-long-random-string")
	v.SetDefault("JWT_ACCESS_EXPIRE", "15m")
	v.SetDefault("JWT_REFRESH_EXPIRE", "168h") // 7 days
}

// DSN returns the database connection string
func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode,
	)
}

// RedisAddr returns the Redis address
func (c *RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// GetDefaults returns a Config with all default values (for testing)
func GetDefaults() *Config {
	return &Config{
		Server: ServerConfig{
			Port: 8888,
			Mode: "release",
		},
		Database: DatabaseConfig{
			Host:            "localhost",
			Port:            5432,
			User:            "postgres",
			Password:        "",
			Name:            "xingrin",
			SSLMode:         "disable",
			MaxOpenConns:    25,
			MaxIdleConns:    5,
			ConnMaxLifetime: 300,
		},
		Redis: RedisConfig{
			Host:     "localhost",
			Port:     6379,
			Password: "",
			DB:       0,
		},
		Log: LogConfig{
			Level:  "info",
			Format: "json",
		},
		JWT: JWTConfig{
			Secret:        "change-me-in-production-use-a-long-random-string",
			AccessExpire:  15 * time.Minute,
			RefreshExpire: 168 * time.Hour,
		},
	}
}
