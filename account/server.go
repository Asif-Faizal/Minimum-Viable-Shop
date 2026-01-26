package account

import (
	"context"
	"fmt"
	"net"
	"net/http"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/Asif-Faizal/Minimum-Viable-Shop/account/pb"
	"github.com/Asif-Faizal/Minimum-Viable-Shop/util"
)

type GrpcServer struct {
	accountService Service
	logger         util.Logger
	pb.UnimplementedAccountServiceServer
}

type restServer struct {
	service Service
	logger  util.Logger
}

func ListenGrpcServer(service Service, logger util.Logger, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			util.UnaryServerInterceptor(logger),
		)),
	)

	server := &GrpcServer{
		accountService: service,
		logger:         logger,
	}
	pb.RegisterAccountServiceServer(grpcServer, server)
	reflection.Register(grpcServer)

	logger.Transport().Info().Int("port", port).Msg("gRPC server listening")
	return grpcServer.Serve(lis)
}

func ListenRestServer(service Service, logger util.Logger, port int) error {
	addr := fmt.Sprintf(":%d", port)
	server := &restServer{
		service: service,
		logger:  logger,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", server.handleHealth)
	mux.HandleFunc("/accounts/check-email", server.handleCheckEmail)

	return http.ListenAndServe(addr, mux)
}

func (s *restServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	util.WriteJSONResponse(w, http.StatusOK, true, "", map[string]string{
		"service": "account",
	})
}

func (s *restServer) handleCheckEmail(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "email is required", http.StatusBadRequest)
		return
	}

	exists, err := s.service.CheckEmailExists(r.Context(), email)
	if err != nil {
		s.logger.Service().Error().Err(err).Msg("failed to check if email exists")
		util.WriteJSONResponse(w, http.StatusInternalServerError, false, err.Error(), nil)
		return
	}

	message := ""
	if exists {
		message = "Email already exists"
	} else {
		message = "Email is available"
	}

	util.WriteJSONResponse(w, http.StatusOK, true, message, map[string]bool{
		"exists": exists,
	})
}

func (server *GrpcServer) CreateOrUpdateAccount(ctx context.Context, request *pb.CreateOrUpdateAccountRequest) (*pb.CreateOrUpdateAccountResponse, error) {
	account, err := server.accountService.CreateOrUpdateAccount(ctx, &Account{
		ID:       request.Id,
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

func (server *GrpcServer) CheckEmailExists(ctx context.Context, request *pb.CheckEmailExistsRequest) (*pb.CheckEmailExistsResponse, error) {
	exists, err := server.accountService.CheckEmailExists(ctx, request.Email)
	if err != nil {
		return nil, err
	}
	return &pb.CheckEmailExistsResponse{Exists: exists}, nil
}
