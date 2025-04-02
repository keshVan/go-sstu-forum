package grpcapp

import (
	"fmt"
	"net"

	authgrpc "github.com/keshvan/go-sstu-forum/internal/grpc_server/auth"
	"github.com/keshvan/go-sstu-forum/pkg/logger"
	grpc "google.golang.org/grpc"
)

type GRPC struct {
	log        *logger.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(log *logger.Logger, port int) *GRPC {
	gRPCServer := grpc.NewServer()
	authgrpc.Register(gRPCServer)

	return &GRPC{log: log, gRPCServer: gRPCServer, port: port}
}

func (g *GRPC) Run() error {
	const op = "grpcapp.Run"
	g.log.Info(fmt.Sprintf(op+" Starting gRPC server port=%d", g.port))

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", g.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	g.log.Info(fmt.Sprintf("%s gRPC server is running port=%d address=%s", op, g.port, l.Addr().String()))

	err = g.gRPCServer.Serve(l)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (g *GRPC) Stop() {
	const op = "grpcapp.Stop"

	g.gRPCServer.GracefulStop()

	g.log.Info(fmt.Sprintf("%s stopping gRPC server", op))
}
