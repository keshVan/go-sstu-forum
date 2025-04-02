package app

import (
	"github.com/keshvan/go-sstu-forum/config"
	grpcapp "github.com/keshvan/go-sstu-forum/pkg/grpc_app"
	"github.com/keshvan/go-sstu-forum/pkg/logger"
)

type App struct {
	GRPC *grpcapp.GRPC
}

func New(cfg *config.Config) *App {
	l := logger.New(cfg.Log)
	grpc := grpcapp.New(l, cfg.GRPC.Port)
	return &App{GRPC: grpc}
}
