package order

import (
	"context"
	"time"

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
	createdAt := order.CreatedAt
	if id == "" {
		id = ksuid.New().String()
		createdAt = time.Now().UTC()
	}
	newOrder := &Order{
		ID:         id,
		CreatedAt:  createdAt,
		AccountID:  order.AccountID,
		TotalPrice: order.TotalPrice,
		Products:   order.Products,
	}
	newOrder.TotalPrice = 0.0
	for _, product := range newOrder.Products {
		newOrder.TotalPrice += product.Price * float64(product.Quantity)
	}
	if _, err := service.repository.CreateOrUpdateOrder(ctx, newOrder); err != nil {
		return nil, err
	}
	return newOrder, nil
}

func (service *OrderService) GetOrderById(ctx context.Context, id string) (*Order, error) {
	order, err := service.repository.GetOrderById(ctx, id)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (service *OrderService) GetOrdersForAccount(ctx context.Context, accountID string) ([]*Order, error) {
	orders, err := service.repository.GetOrdersForAccount(ctx, accountID)
	if err != nil {
		return nil, err
	}
	return orders, nil
}
