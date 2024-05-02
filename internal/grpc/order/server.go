package order

import (
	"context"
	"errors"
	"github.com/smugglerv1/internal/repository"
	"github.com/smugglerv1/pkg/orders"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Order interface {
	CreateOrder(ctx context.Context, customerId int32, productId int32, quantity int32) (string, error)
	GetOrderCache(ctx context.Context, customerId int32) ([]*orders.Order, error)
}
type serverApi struct {
	orders.UnimplementedOrderServiceServer
	order Order
}

func Register(gRPC *grpc.Server, order Order) {
	orders.RegisterOrderServiceServer(gRPC, &serverApi{order: order})
}

const (
	emptyValue = 0
)

func (s *serverApi) CreateOrder(ctx context.Context, req *orders.CreateOrderRequest) (*orders.CreateOrderResponse, error) {
	if err := ValidateCreate(req); err != nil {
		return nil, err
	}
	response, err := s.order.CreateOrder(ctx, req.GetCustomerId(), req.GetProductId(), req.GetQuantity())
	if err != nil {
		if errors.Is(err, repository.ErrUserDoesntExist) {
			return nil, status.Error(codes.InvalidArgument, "customer doesnt exist")
		}
	}
	return &orders.CreateOrderResponse{
		Status: response,
	}, nil
}

func (s *serverApi) GetOrderCache(ctx context.Context, req *orders.GetOrderRequest) (*orders.GetOrderResponse, error) {
	if err := ValidateGet(req); err != nil {
		return nil, err
	}
	response, err := s.order.GetOrderCache(ctx, req.GetCustomerId())
	if err != nil {
		if errors.Is(err, repository.ErrUserDoesntExist) {
			return nil, status.Error(codes.InvalidArgument, "customer doesnt exist")
		}
	}
	return &orders.GetOrderResponse{
		Orders: response,
	}, nil
}

func ValidateCreate(req *orders.CreateOrderRequest) error {
	if req.GetCustomerId() == emptyValue {
		return status.Errorf(codes.InvalidArgument, "customer id ")
	}
	if req.GetQuantity() == emptyValue {
		return status.Errorf(codes.InvalidArgument, "quantity")
	}
	if req.GetProductId() == emptyValue {
		return status.Errorf(codes.InvalidArgument, "product id")
	}
	return nil
}

func ValidateGet(req *orders.GetOrderRequest) error {
	if req.GetCustomerId() == emptyValue {
		return status.Errorf(codes.InvalidArgument, "customer id")
	}
	return nil
}
