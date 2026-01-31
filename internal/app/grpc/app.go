package grpcapp

import (
	"fmt"
	"log/slog"
	"net"

	authgrpc "github.com/go_grpc/auth/internal/grpc/auth"
	"google.golang.org/grpc"
)

type App struct {
	log       *slog.Logger
	grpcSever *grpc.Server
	port      int
}

func NewApp(
	log *slog.Logger,
	port int,
	authService authgrpc.Auth,
) *App {
	grpcServer := grpc.NewServer()

	authgrpc.RegisterServerAPI(grpcServer, authService)

	return &App{
		log:       log,
		grpcSever: grpcServer,
		port:      port,
	}
}

func (a *App) MustStart() {
	if err := a.Start(); err != nil {
		panic(err)
	}
}

func (a *App) Start() error {
	const op = "grpcapp.App.Start"
	log := a.log.With(
		slog.String("op", op),
		slog.Int("port", a.port),
	)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	log.Info(
		"starting grpc server",
		slog.String("addr", l.Addr().String()),
	)

	if err := a.grpcSever.Serve(l); err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	return nil
}

func (a *App) Stop() {
	const op = "grpcapp.App.Stop"

	a.log.
		With(slog.String("op", op)).
		Info("stopping GRPC server", slog.Int("port", a.port))

	a.grpcSever.GracefulStop()
}
