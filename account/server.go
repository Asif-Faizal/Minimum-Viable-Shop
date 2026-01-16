package account

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/Asif-Faizal/Minimum-Viable-Shop/account/pb"
)

type GrpcServer struct {
	accountService Service
	pb.UnimplementedAccountServiceServer
}

func ListenGrpcServer(service Service, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer()
	server := &GrpcServer{accountService: service}
	pb.RegisterAccountServiceServer(grpcServer, server)
	reflection.Register(grpcServer)
	return grpcServer.Serve(lis)
}

func (server *GrpcServer) CreateOrUpdateAccount(ctx context.Context, request *pb.CreateOrUpdateAccountRequest) (*pb.CreateOrUpdateAccountResponse, error) {
	account, err := server.accountService.CreateOrUpdateAccount(ctx, &Account{
		Name:     request.Name,
		Email:    request.Email,
		Password: request.Password,
	})
	if err != nil {
		return nil, err
	}
	return &pb.CreateOrUpdateAccountResponse{
		Account: &pb.Account{
			Id:    account.ID,
			Name:  account.Name,
			Email: account.Email,
		},
	}, nil
}

func (server *GrpcServer) GetAccountByID(ctx context.Context, request *pb.GetAccountByIDRequest) (*pb.GetAccountByIDResponse, error) {
	account, err := server.accountService.GetAccountByID(ctx, request.Id)
	if err != nil {
		return nil, err
	}
	return &pb.GetAccountByIDResponse{
		Account: &pb.Account{
			Id:    account.ID,
			Name:  account.Name,
			Email: account.Email,
		},
	}, nil
}

func (server *GrpcServer) ListAccounts(ctx context.Context, request *pb.ListAccountsRequest) (*pb.ListAccountsResponse, error) {
	domainAccounts, err := server.accountService.ListAccounts(ctx, uint(request.Skip), uint(request.Take))
	if err != nil {
		return nil, err
	}
	accounts := []*pb.Account{}
	for _, a := range domainAccounts {
		accounts = append(accounts, &pb.Account{
			Id:    a.ID,
			Name:  a.Name,
			Email: a.Email,
		})
	}
	return &pb.ListAccountsResponse{Accounts: accounts}, nil
}
