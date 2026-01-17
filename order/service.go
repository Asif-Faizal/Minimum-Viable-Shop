package order

import (
	"context"

	"github.com/segmentio/ksuid"
)

type Service interface {
	CreateOrUpdateOrder(ctx context.Context, order *Order) (*Order, error)
	GetOrderById(ctx context.Context, id string) (*Order, error)
	GetOrdersForAccount(ctx context.Context, accountID string) ([]*Order, error)
}

type OrderService struct {
	repository Repository
}

func NewOrderService(repository Repository) *OrderService {
	return &OrderService{repository: repository}
}

func (service *OrderService) CreateOrUpdateOrder(ctx context.Context, order *Order) (*Order, error) {
	id := order.ID
	if id == "" {
		id = ksuid.New().String()
	}
	newOrder := &Order{
		ID:         id,
		CreatedAt:  order.CreatedAt,
		AccountID:  order.AccountID,
		TotalPrice: order.TotalPrice,
		Products:   order.Products,
	}
	if _, err := service.repository.CreateOrUpdateOrder(ctx, newOrder); err != nil {
		return nil, err
	}
	return newOrder, nil
}
