package account

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GrpcServer struct {
	accountService Service
}

func ListenGrpcServer(service Service, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer()
	pb.RegisterAccountServiceServer(grpcServer, server)
	reflection.Register(grpcServer)
	return grpcServer.Serve(lis)
}

// Create and Update account
func (server *GrpcServer) CreateOrUpdateAccount(ctx context.Context, request *pb.CreateOrUpdateAccountRequest) (*pb.Account, error) {
	account, err := server.accountService.CreateOrUpdateAccount(ctx, request.Account)
	if err != nil {
		return nil, err
	}
	return &pb.Account{
		Id:    account.ID,
		Name:  account.Name,
		Email: account.Email,
	}, nil
}

// Get account by ID
func (server *GrpcServer) GetAccountByID(ctx context.Context, request *pb.GetAccountByIDRequest) (*pb.Account, error) {
	account, err := server.accountService.GetAccountByID(ctx, request.Id)
	if err != nil {
		return nil, err
	}
	return &pb.Account{
		Id:    account.ID,
		Name:  account.Name,
		Email: account.Email,
	}, nil
}

// List accounts
func (server *GrpcServer) ListAccounts(ctx context.Context, request *pb.ListAccountsRequest) (*pb.ListAccountsResponse, error) {
	accounts, err := server.accountService.ListAccounts(ctx, request.Skip, request.Take)
	if err != nil {
		return nil, err
	}
	return &pb.ListAccountsResponse{Accounts: accounts}, nil
}
