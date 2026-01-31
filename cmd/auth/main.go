package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/go_grpc/auth/internal/app"
	"github.com/go_grpc/auth/internal/config"
)

const (
	envLocal   = "local"
	envDev     = "dev"
	envStaging = "staging"
	envProd    = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("Config loaded", slog.String("env", cfg.Env), slog.Int("grpc_port", cfg.GRPC.Port))

	application := app.NewApp(log, cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTTL)

	go application.GRPCSrv.MustStart()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	s := <-stop

	log.Info("Server stopping")
	application.GRPCSrv.Stop()
	log.Info("Server stopped", slog.String("signal", s.String()))
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envStaging:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn}))
	default:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
