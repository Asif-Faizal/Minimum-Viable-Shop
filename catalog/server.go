package catalog

import (
	"context"
	"fmt"
	"net"

	pb "github.com/Asif-Faizal/Minimum-Viable-Shop/catalog/pb"
	"github.com/Asif-Faizal/Minimum-Viable-Shop/util"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GrpcServer struct {
	catalogService Service
	logger         util.Logger
	pb.UnimplementedCatalogServiceServer
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
	server := &GrpcServer{catalogService: service, logger: logger}
	pb.RegisterCatalogServiceServer(grpcServer, server)
	reflection.Register(grpcServer)
	return grpcServer.Serve(lis)
}

func (server *GrpcServer) CreateOrUpdateProduct(ctx context.Context, request *pb.CreateOrUpdateProductRequest) (*pb.CreateOrUpdateProductResponse, error) {
	product, err := server.catalogService.CreateOrUpdateProduct(ctx, &Product{
		Name:        request.Name,
		Description: request.Description,
		Price:       request.Price,
	})
	if err != nil {
		return nil, err
	}
	return &pb.CreateOrUpdateProductResponse{
		Product: &pb.Product{
			Id:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
		},
	}, nil
}

func (server *GrpcServer) GetProductByID(ctx context.Context, request *pb.GetProductByIDRequest) (*pb.GetProductByIDResponse, error) {
	product, err := server.catalogService.GetProductById(ctx, request.Id)
	if err != nil {
		return nil, err
	}
	return &pb.GetProductByIDResponse{
		Product: &pb.Product{
			Id:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
		},
	}, nil
}

func (server *GrpcServer) ListProducts(ctx context.Context, request *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	domainProducts, err := server.catalogService.ListProducts(ctx, request.Skip, request.Take)
	if err != nil {
		return nil, err
	}
	products := []*pb.Product{}
	for _, product := range domainProducts {
		products = append(products, &pb.Product{
			Id:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
		})
	}
	return &pb.ListProductsResponse{Products: products}, nil
}

func (server *GrpcServer) ListProductsWithIds(ctx context.Context, request *pb.ListProductsWithIdsRequest) (*pb.ListProductsWithIdsResponse, error) {
	products, err := server.catalogService.ListProductsWithIds(ctx, request.Ids)
	if err != nil {
		return nil, err
	}
	grpcProducts := []*pb.Product{}
	for _, product := range products {
		grpcProducts = append(grpcProducts, &pb.Product{
			Id:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
		})
	}
	return &pb.ListProductsWithIdsResponse{Products: grpcProducts}, nil
}

func (server *GrpcServer) SearchProducts(ctx context.Context, request *pb.SearchProductsRequest) (*pb.SearchProductsResponse, error) {
	products, err := server.catalogService.SearchProducts(ctx, request.Query, request.Skip, request.Take)
	if err != nil {
		return nil, err
	}
	grpcProducts := []*pb.Product{}
	for _, product := range products {
		grpcProducts = append(grpcProducts, &pb.Product{
			Id:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
		})
	}
	return &pb.SearchProductsResponse{Products: grpcProducts}, nil
}
