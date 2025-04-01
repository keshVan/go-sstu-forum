// Package app configures and runs application.
package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/keshvan/go-sstu-forum/config"
	amqprpc "github.com/keshvan/go-sstu-forum/internal/controller/amqp_rpc"
	v1 "github.com/keshvan/go-sstu-forum/internal/controller/http"
	"github.com/keshvan/go-sstu-forum/internal/repo/persistent"
	"github.com/keshvan/go-sstu-forum/internal/repo/webapi"
	"github.com/keshvan/go-sstu-forum/internal/usecase/translation"
	"github.com/keshvan/go-sstu-forum/pkg/httpserver"
	"github.com/keshvan/go-sstu-forum/pkg/logger"
	"github.com/keshvan/go-sstu-forum/pkg/postgres"
	"github.com/keshvan/go-sstu-forum/pkg/rabbitmq/rmq_rpc/server"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	// Repository
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	// Use case
	translationUseCase := translation.New(
		persistent.New(pg),
		webapi.New(),
	)

	// RabbitMQ RPC Server
	rmqRouter := amqprpc.NewRouter(translationUseCase)

	rmqServer, err := server.New(cfg.RMQ.URL, cfg.RMQ.ServerExchange, rmqRouter, l)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - rmqServer - server.New: %w", err))
	}

	// HTTP Server
	httpServer := httpserver.New(httpserver.Port(cfg.HTTP.Port), httpserver.Prefork(cfg.HTTP.UsePreforkMode))
	v1.NewRouter(httpServer.App, cfg, l, translationUseCase)

	// Start servers
	rmqServer.Start()
	httpServer.Start()

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	case err = <-rmqServer.Notify():
		l.Error(fmt.Errorf("app - Run - rmqServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}

	err = rmqServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - rmqServer.Shutdown: %w", err))
	}
}
