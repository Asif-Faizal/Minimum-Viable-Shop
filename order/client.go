package order

import (
	"context"

	"github.com/Asif-Faizal/Minimum-Viable-Shop/order/pb/pb"
	"google.golang.org/grpc"
)

type OrderClient struct {
	connection *grpc.ClientConn
	client     pb.OrderServiceClient
}

func NewOrderClient(url string) (*OrderClient, error) {
	connection, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return &OrderClient{
		connection: connection,
		client:     pb.NewOrderServiceClient(connection),
	}, nil
}

func (client *OrderClient) Close() {
	client.connection.Close()
}

// CreateOrUpdate Order
func (client *OrderClient) CreateOrUpdateOrder(ctx context.Context, order *pb.Order) (*pb.CreateOrUpdateOrderResponse, error) {
	response, err := client.client.CreateOrUpdateOrder(ctx, &pb.CreateOrUpdateOrderRequest{
		Order: order,
	})
	if err != nil {
		return nil, err
	}
	return response, nil
}

// Get Order by ID
func (client *OrderClient) GetOrderByID(ctx context.Context, id string) (*pb.GetOrderByIDResponse, error) {
	response, err := client.client.GetOrderByID(ctx, &pb.GetOrderByIDRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}
	return response, nil
}

// List Orders for Account
func (client *OrderClient) ListOrdersForAccount(ctx context.Context, accountId string) (*pb.GetOrdersForAccountResponse, error) {
	response, err := client.client.GetOrdersForAccount(ctx, &pb.GetOrdersForAccountRequest{
		AccountId: accountId,
	})
	if err != nil {
		return nil, err
	}
	return response, nil
}
