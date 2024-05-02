package app

import (
	orderApp "github.com/smugglerv1/internal/app/order"
	"github.com/smugglerv1/internal/repository/postgres"
	"github.com/smugglerv1/internal/repository/redis"
	"github.com/smugglerv1/internal/service/order"
	"log/slog"
)

type App struct {
	GRPSrv *orderApp.App
}

func New(log *slog.Logger, grpcPort int, db string) *App {
	dbStorage, err := postgres.New(db)
	cacheStorage := redis.New()
	if err != nil {
		panic(err)
	}
	orderService := order.New(log, dbStorage, dbStorage, cacheStorage)
	oApp := orderApp.New(log, orderService, grpcPort)

	return &App{
		GRPSrv: oApp,
	}
}
