package account

import (
	"context"

	"github.com/Asif-Faizal/Minimum-Viable-Shop/account/pb"
	"google.golang.org/grpc"
)

type AccountClient struct {
	connection *grpc.ClientConn
	client     pb.AccountServiceClient
}

func NewAccountClient(url string) (*AccountClient, error) {
	connection, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return &AccountClient{
		connection: connection,
		client:     pb.NewAccountServiceClient(connection),
	}, nil
}

func (client *AccountClient) Close() {
	client.connection.Close()
}

// CreateOrUpdate Account
func (client *AccountClient) CreateOrUpdateAccount(ctx context.Context, name, email, password string) (*pb.CreateOrUpdateAccountResponse, error) {
	response, err := client.client.CreateOrUpdateAccount(ctx, &pb.CreateOrUpdateAccountRequest{
		Name:     name,
		Email:    email,
		Password: password,
	})
	if err != nil {
		return nil, err
	}
	return response, nil
}
