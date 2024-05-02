package order

import (
	"context"
	"errors"
	"fmt"
	"github.com/smugglerv1/internal/repository"
	"github.com/smugglerv1/pkg/orders"
	"log/slog"
)

type Order struct {
	log           *slog.Logger
	orderSaver    OrderSaver
	orderProvider OrderProvider
	cachedOrder   CachedOrder
}
type OrderSaver interface {
	SaveOrder(ctx context.Context, customerId int32, productId int32, quantity int32) (string, error)
}
type OrderProvider interface {
	GetOrder(ctx context.Context, customerId int32) ([]*orders.Order, error)
}
type CachedOrder interface {
	GetOrderCached(ctx context.Context, dbOrders []*orders.Order) ([]*orders.Order, error)
}

var (
	ErrInvalidUserCredentials = errors.New("invalid credentials")
)

func New(log *slog.Logger, orderSaver OrderSaver, orderProvider OrderProvider, cachedOrder CachedOrder) *Order {
	return &Order{log: log, orderSaver: orderSaver, orderProvider: orderProvider, cachedOrder: cachedOrder}
}

func (o *Order) CreateOrder(ctx context.Context, customerId int32, productId int32, quantity int32) (string, error) {
	const op = "order.createOrder"
	log := o.log.With(
		slog.String("op", op),
	)
	log.Info("creating order")

	response, err := o.orderSaver.SaveOrder(ctx, customerId, productId, quantity)
	if err != nil {
		fmt.Print(err)
		log.Error("error creating order", err)
		return "", fmt.Errorf("%s:%w", op, err)
	}
	log.Info("order created")
	return response, nil
}
func (o *Order) GetOrderCache(ctx context.Context, customerId int32) ([]*orders.Order, error) {
	const op = "order.getOrderCache"
	log := o.log.With(
		slog.String("op", op),
	)
	dbOrders, err := o.GetOrderFromDB(ctx, customerId)
	if err != nil {
		fmt.Printf("error getting order from database")
	}
	response, err := o.cachedOrder.GetOrderCached(ctx, dbOrders)
	if err != nil {
		fmt.Print(err)
		log.Error("error getting order from REDIS cache")
		return nil, fmt.Errorf("%s:%w", op, err)
	}
	log.Info("getting order cache")
	return response, nil
}
func (o *Order) GetOrderFromDB(ctx context.Context, customerId int32) ([]*orders.Order, error) {
	const op = "order.getOrderFromDB"
	log := o.log.With(
		slog.String("op", op),
	)
	log.Info("getting order from db")
	orders, err := o.orderProvider.GetOrder(ctx, customerId)
	if err != nil {
		if errors.Is(err, repository.ErrUserDoesntExist) {
			o.log.Warn("customer does not exist", err)
			return nil, fmt.Errorf("%s:%w", op, ErrInvalidUserCredentials)
		}
		o.log.Error("failed to get order", err)
		return nil, fmt.Errorf("%s:%w", op, err)
	}
	return orders, nil
}
