package app

import (
	grpcapp "github.com/zaketn/sso/internal/app/grpc"
	"log/slog"
	"time"
)

type App struct {
	GRPCServer grpcapp.App
}

func New(log *slog.Logger, grpcPort int, storagePath string, tokenTTL time.Duration) *App {
	server := grpcapp.New(log, grpcPort)

	return &App{GRPCServer: *server}
}
