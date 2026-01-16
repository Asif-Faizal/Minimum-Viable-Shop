package catalog

import (
	"fmt"
	"net"

	pb "github.com/Asif-Faizal/Minimum-Viable-Shop/catalog/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GrpcServer struct {
	catalogService Service
	pb.UnimplementedCatalogServiceServer
}

func ListenGrpcServer(service Service, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer()
	server := &GrpcServer{catalogService: service}
	pb.RegisterCatalogServiceServer(grpcServer, server)
	reflection.Register(grpcServer)
	return grpcServer.Serve(lis)
}
