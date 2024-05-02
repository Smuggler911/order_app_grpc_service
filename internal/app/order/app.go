package order

import (
	"fmt"
	health2 "github.com/smugglerv1/internal/grpc/health"
	ordergrpc "github.com/smugglerv1/internal/grpc/order"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"log/slog"
	"net"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(log *slog.Logger, orderService ordergrpc.Order, port int) *App {
	gRPCServer := grpc.NewServer()

	ordergrpc.Register(gRPCServer, orderService)
	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "grpcapp.Run"
	log := a.log.With(
		slog.String("op", op),
		slog.Int("port", a.port),
	)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s : %w ", op, err)
	}
	log.Info("grpc server is running", slog.String("addr", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s : %w ", op, err)
	}

	health := health2.NewHServer()
	grpc_health_v1.RegisterHealthServer(a.gRPCServer, health)

	return nil
}

func (a *App) Stop() {
	const op = "grpcapp.Stop"

	a.log.With(slog.String("op", op)).Info("stopping grpc server", slog.Int("port", a.port))

	a.gRPCServer.GracefulStop()

}
