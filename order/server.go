package order

import (
	"fmt"
	"net"

	"github.com/Asif-Faizal/Minimum-Viable-Shop/account"
	"github.com/Asif-Faizal/Minimum-Viable-Shop/catalog"
	"github.com/Asif-Faizal/Minimum-Viable-Shop/order/pb/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GrpcServer struct {
	orderService  Service
	accountClient *account.AccountClient
	catalogClient *catalog.CatalogClient
	pb.UnimplementedOrderServiceServer
}

func ListenGrpcServer(service Service, accountUrl string, catalogUrl string, port int) error {
	accountClient, err := account.NewAccountClient(accountUrl)
	if err != nil {
		return err
	}
	catalogClient, err := catalog.NewCatalogClient(catalogUrl)
	if err != nil {
		accountClient.Close()
		return err
	}
	grpcServer := grpc.NewServer()
	server := &GrpcServer{orderService: service, accountClient: accountClient, catalogClient: catalogClient}
	pb.RegisterOrderServiceServer(grpcServer, server)
	reflection.Register(grpcServer)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		accountClient.Close()
		catalogClient.Close()
		return err
	}
	return grpcServer.Serve(lis)
}
