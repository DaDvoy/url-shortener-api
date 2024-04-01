package app

import (
	grpcapp "github.com/DaDvoy/url-shortener-api.git/internal/app/grpc"
	"github.com/DaDvoy/url-shortener-api.git/internal/services"
	"github.com/DaDvoy/url-shortener-api.git/internal/storage"
	"golang.org/x/exp/slog"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(log *slog.Logger, storage storage.Storage, grpcPort int) *App {
	apiService := services.New(log, storage, storage)

	grpcApp := grpcapp.New(log, apiService, grpcPort)

	return &App{GRPCSrv: grpcApp}
}
