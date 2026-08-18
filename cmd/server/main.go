package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/example/inspection-platform/api"
	anomalyhttp "github.com/example/inspection-platform/internal/anomaly/adapter/http"
	anomalyapp "github.com/example/inspection-platform/internal/anomaly/application"
	anomalyinfra "github.com/example/inspection-platform/internal/anomaly/infrastructure"
	"github.com/example/inspection-platform/internal/config"
	devicehttp "github.com/example/inspection-platform/internal/device/adapter/http"
	deviceapp "github.com/example/inspection-platform/internal/device/application"
	deviceinfra "github.com/example/inspection-platform/internal/device/infrastructure"
	planhttp "github.com/example/inspection-platform/internal/plan/adapter/http"
	planapp "github.com/example/inspection-platform/internal/plan/application"
	planinfra "github.com/example/inspection-platform/internal/plan/infrastructure"
	"github.com/example/inspection-platform/internal/platform/database"
	"github.com/example/inspection-platform/internal/platform/middleware"
	recordhttp "github.com/example/inspection-platform/internal/record/adapter/http"
	recordapp "github.com/example/inspection-platform/internal/record/application"
	recordinfra "github.com/example/inspection-platform/internal/record/infrastructure"
	repairhttp "github.com/example/inspection-platform/internal/repair/adapter/http"
	repairapp "github.com/example/inspection-platform/internal/repair/application"
	repairinfra "github.com/example/inspection-platform/internal/repair/infrastructure"
	"github.com/example/inspection-platform/internal/scheduler"
	taskhttp "github.com/example/inspection-platform/internal/task/adapter/http"
	taskapp "github.com/example/inspection-platform/internal/task/application"
	taskinfra "github.com/example/inspection-platform/internal/task/infrastructure"
	userhttp "github.com/example/inspection-platform/internal/user/adapter/http"
	userapp "github.com/example/inspection-platform/internal/user/application"
	userinfra "github.com/example/inspection-platform/internal/user/infrastructure"
)

func main() {
	configPath := flag.String("config", "configs/config.yaml", "path to config file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "load config:", err)
		os.Exit(1)
	}
	setupLogger(cfg.Log.Level, cfg.Log.Format)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	db, err := database.Open(database.Options{
		Driver:       cfg.Database.Driver,
		SQLitePath:   cfg.Database.SQLitePath,
		PostgresDSN:  cfg.Database.PostgresDSN,
		MaxOpenConns: cfg.Database.MaxOpenConns,
		MaxIdleConns: cfg.Database.MaxIdleConns,
	})
	if err != nil {
		slog.Error("open database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	if cfg.Database.Driver == "sqlite" {
		if err := database.EnsureSQLiteSchema(db); err != nil {
			slog.Error("ensure sqlite schema", "error", err)
			os.Exit(1)
		}
	}

	userRepo := userinfra.NewSQLiteUserRepository(db)
	deviceRepo := deviceinfra.NewSQLiteDeviceRepository(db)
	planRepo := planinfra.NewSQLitePlanRepository(db)
	taskRepo := taskinfra.NewSQLiteTaskRepository(db)
	recordRepo := recordinfra.NewSQLiteRecordRepository(db)
	anomalyRepo := anomalyinfra.NewSQLiteAnomalyRepository(db)
	repairRepo := repairinfra.NewSQLiteRepairRepository(db)

	userService := userapp.NewService(userRepo, cfg.Auth.JWTSecret, cfg.Auth.TokenTTL)
	deviceService := deviceapp.NewService(deviceRepo, recordRepo, repairRepo)
	planService := planapp.NewService(planRepo, deviceRepo, taskRepo)
	taskService := taskapp.NewService(taskRepo, planRepo, recordRepo, anomalyRepo)
	recordService := recordapp.NewService(recordRepo)
	anomalyService := anomalyapp.NewService(anomalyRepo)
	repairService := repairapp.NewService(repairRepo, anomalyRepo)

	if err := userService.EnsureDefault(ctx, cfg.Auth.DefaultUsername, cfg.Auth.DefaultPassword, cfg.Auth.DefaultDisplayName); err != nil {
		slog.Error("ensure default user", "error", err)
		os.Exit(1)
	}

	metrics := middleware.NewMetrics()
	handler := api.NewRouter(api.Dependencies{
		DB:        db,
		Metrics:   metrics,
		JWTSecret: cfg.Auth.JWTSecret,
		Users:     userhttp.NewHandler(userService),
		Devices:   devicehttp.NewHandler(deviceService),
		Plans:     planhttp.NewHandler(planService),
		Tasks:     taskhttp.NewHandler(taskService),
		Records:   recordhttp.NewHandler(recordService),
		Anomalies: anomalyhttp.NewHandler(anomalyService),
		Repairs:   repairhttp.NewHandler(repairService),
	})

	sched := scheduler.New(planService, cfg.Scheduler.Interval)
	if cfg.Scheduler.Enabled {
		sched.Start(ctx)
		if count, err := planService.GenerateDue(ctx); err != nil {
			slog.Warn("initial task generation failed", "error", err)
		} else if count > 0 {
			slog.Info("initial task generation completed", "count", count)
		}
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      http.TimeoutHandler(handler, cfg.Server.WriteTimeout+5*time.Second, `{"error":"request timeout"}`),
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	go func() {
		slog.Info("server listening", "addr", server.Addr, "driver", cfg.Database.Driver)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server error", "error", err)
			stop()
		}
	}()

	<-ctx.Done()
	sched.Stop()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("server shutdown", "error", err)
	}
	slog.Info("server stopped")
}

func setupLogger(level, format string) {
	var slogLevel slog.Level
	switch level {
	case "debug":
		slogLevel = slog.LevelDebug
	case "warn":
		slogLevel = slog.LevelWarn
	case "error":
		slogLevel = slog.LevelError
	default:
		slogLevel = slog.LevelInfo
	}
	opts := &slog.HandlerOptions{Level: slogLevel}
	if format == "text" {
		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, opts)))
		return
	}
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, opts)))
}
