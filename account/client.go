package account

import (
	"context"

	pb "github.com/Asif-Faizal/Minimum-Viable-Shop/account/pb"
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
func (client *AccountClient) CreateOrUpdateAccount(ctx context.Context, id, name, email, password string) (*pb.CreateOrUpdateAccountResponse, error) {
	response, err := client.client.CreateOrUpdateAccount(ctx, &pb.CreateOrUpdateAccountRequest{
		Id:       id,
		Name:     name,
		Email:    email,
		Password: password,
	})
	if err != nil {
		return nil, err
	}
	return response, nil
}

// Get Account by ID
func (client *AccountClient) GetAccountByID(ctx context.Context, id string) (*pb.GetAccountByIDResponse, error) {
	response, err := client.client.GetAccountByID(ctx, &pb.GetAccountByIDRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}
	return response, nil
}

// List Accounts
func (client *AccountClient) ListAccounts(ctx context.Context, skip uint32, take uint32) (*pb.ListAccountsResponse, error) {
	response, err := client.client.ListAccounts(ctx, &pb.ListAccountsRequest{
		Skip: skip,
		Take: take,
	})
	if err != nil {
		return nil, err
	}
	return response, nil
}

// Check Email Exists
func (client *AccountClient) CheckEmailExists(ctx context.Context, email string) (bool, error) {
	response, err := client.client.CheckEmailExists(ctx, &pb.CheckEmailExistsRequest{
		Email: email,
	})
	if err != nil {
		return false, err
	}
	return response.Exists, nil
}

// Login
func (client *AccountClient) Login(ctx context.Context, email, password string) (*pb.LoginResponse, error) {
	response, err := client.client.Login(ctx, &pb.LoginRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return nil, err
	}
	return response, nil
}
