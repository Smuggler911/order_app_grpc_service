package redis

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"github.com/smugglerv1/internal/domain/models"
	"github.com/smugglerv1/pkg/orders"
	"time"
)

type Redis struct {
	redisClient *redis.Client
}
type JsonResponse struct {
	Data   []models.Order `json:"data"`
	Source string         `json:"source"`
}

func New() *Redis {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "http://213.226.127.170:6379",
		Password: "",
		DB:       0,
	})
	return &Redis{
		redisClient: rdb,
	}

}

func (r *Redis) GetOrderCached(ctx context.Context, dbOrders []*orders.Order) (orders []*orders.Order, err error) {
	cachedOrders, err := r.redisClient.Get(ctx, "orders").Bytes()
	//response := models.Order{}
	if err != nil {
		cachedData, err := json.Marshal(dbOrders)
		if err != nil {
			return nil, err
		}
		err = r.redisClient.Set(ctx, "orders", cachedData, 10*time.Second).Err()
		if err != nil {
			return nil, err
		}
	}
	err = json.Unmarshal(cachedOrders, &orders)
	if err != nil {
		return nil, err
	}

	return orders, nil
}
