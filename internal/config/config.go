package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"gopkg.in/yaml.v3"
)

// Config is the top-level runtime configuration for the platform.
type Config struct {
	Server    ServerConfig    `yaml:"server"`
	Database  DatabaseConfig  `yaml:"database"`
	Auth      AuthConfig      `yaml:"auth"`
	Scheduler SchedulerConfig `yaml:"scheduler"`
	Log       LogConfig       `yaml:"log"`
}

type ServerConfig struct {
	Port         int           `yaml:"port"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	IdleTimeout  time.Duration `yaml:"idle_timeout"`
}

type DatabaseConfig struct {
	Driver       string `yaml:"driver"`
	SQLitePath   string `yaml:"sqlite_path"`
	PostgresDSN  string `yaml:"postgres_dsn"`
	MaxOpenConns int    `yaml:"max_open_conns"`
	MaxIdleConns int    `yaml:"max_idle_conns"`
}

type AuthConfig struct {
	JWTSecret          string        `yaml:"jwt_secret"`
	TokenTTL           time.Duration `yaml:"token_ttl"`
	DefaultUsername    string        `yaml:"default_username"`
	DefaultPassword    string        `yaml:"default_password"`
	DefaultDisplayName string        `yaml:"default_display_name"`
}

type SchedulerConfig struct {
	Enabled  bool          `yaml:"enabled"`
	Interval time.Duration `yaml:"interval"`
}

type LogConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

func Default() Config {
	return Config{
		Server: ServerConfig{
			Port:         8080,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
		Database: DatabaseConfig{
			Driver:       "sqlite",
			SQLitePath:   "data/inspection.db",
			MaxOpenConns: 5,
			MaxIdleConns: 2,
		},
		Auth: AuthConfig{
			JWTSecret:          "local-development-secret-change-me",
			TokenTTL:           12 * time.Hour,
			DefaultUsername:    "admin",
			DefaultPassword:    "admin123",
			DefaultDisplayName: "系统管理员",
		},
		Scheduler: SchedulerConfig{Enabled: true, Interval: time.Minute},
		Log:       LogConfig{Level: "info", Format: "json"},
	}
}

func Load(path string) (Config, error) {
	cfg := Default()
	if path != "" {
		raw, err := os.ReadFile(path)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return cfg, nil
			}
			return cfg, fmt.Errorf("read config: %w", err)
		}
		if err := yaml.Unmarshal(raw, &cfg); err != nil {
			return cfg, fmt.Errorf("decode config: %w", err)
		}
	}
	applyEnv(&cfg)
	return cfg, nil
}

func applyEnv(cfg *Config) {
	if v := os.Getenv("APP_PORT"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			cfg.Server.Port = n
		}
	}
	if v := os.Getenv("DB_DRIVER"); v != "" {
		cfg.Database.Driver = v
	}
	if v := os.Getenv("SQLITE_PATH"); v != "" {
		cfg.Database.SQLitePath = v
	}
	if v := os.Getenv("POSTGRES_DSN"); v != "" {
		cfg.Database.PostgresDSN = v
	}
	if v := os.Getenv("JWT_SECRET"); v != "" {
		cfg.Auth.JWTSecret = v
	}
	if v := os.Getenv("LOG_LEVEL"); v != "" {
		cfg.Log.Level = v
	}
}
