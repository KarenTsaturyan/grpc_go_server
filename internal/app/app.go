package app

import (
	"log/slog"
	"time"

	grpcapp "github.com/go_grpc/auth/internal/app/grpc"
	"github.com/go_grpc/auth/internal/services/auth"
	"github.com/go_grpc/auth/internal/storage/sqlite"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func NewApp(
	log *slog.Logger,
	grpcPort int,
	storagePath string,
	tokenTTL time.Duration,
) *App {
	storage, err := sqlite.New(storagePath)
	if err != nil {
		panic(err)
	}

	authService := auth.New(log, storage, storage, storage, tokenTTL)

	grpcApp := grpcapp.NewApp(log, grpcPort, authService)
	return &App{
		GRPCSrv: grpcApp,
	}
}
