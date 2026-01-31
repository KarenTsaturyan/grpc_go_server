package app

import (
	"log/slog"
	"time"

	grpcapp "github.com/go_grpc/auth/internal/app/grpc"
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
	//  TODO: INIT other services (storage, auth service, etc)

	grpcApp := grpcapp.NewApp(log, grpcPort)
	return &App{
		GRPCSrv: grpcApp,
	}
}
