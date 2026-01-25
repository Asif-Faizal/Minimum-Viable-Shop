package catalog

import (
	"context"

	pb "github.com/Asif-Faizal/Minimum-Viable-Shop/catalog/pb"
	"google.golang.org/grpc"
)

type CatalogClient struct {
	connection *grpc.ClientConn
	client     pb.CatalogServiceClient
}

func NewCatalogClient(url string) (*CatalogClient, error) {
	connection, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return &CatalogClient{
		connection: connection,
		client:     pb.NewCatalogServiceClient(connection),
	}, nil
}

func (client *CatalogClient) Close() {
	client.connection.Close()
}

// CreateOrUpdate Product
func (client *CatalogClient) CreateOrUpdateProduct(ctx context.Context, id, name, description string, price float32) (*pb.CreateOrUpdateProductResponse, error) {
	response, err := client.client.CreateOrUpdateProduct(ctx, &pb.CreateOrUpdateProductRequest{
		Id:          id,
		Name:        name,
		Description: description,
		Price:       price,
	})
	if err != nil {
		return nil, err
	}
	return response, nil
}

// Get Product by ID
func (client *CatalogClient) GetProductByID(ctx context.Context, id string) (*pb.GetProductByIDResponse, error) {
	response, err := client.client.GetProductByID(ctx, &pb.GetProductByIDRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}
	return response, nil
}

// List Products
func (client *CatalogClient) ListProducts(ctx context.Context, skip uint64, take uint64) (*pb.ListProductsResponse, error) {
	response, err := client.client.ListProducts(ctx, &pb.ListProductsRequest{
		Skip: skip,
		Take: take,
	})
	if err != nil {
		return nil, err
	}
	return response, nil
}

// List Products With Ids
func (client *CatalogClient) ListProductsWithIDs(ctx context.Context, ids []string) (*pb.ListProductsWithIdsResponse, error) {
	response, err := client.client.ListProductsWithIds(ctx, &pb.ListProductsWithIdsRequest{
		Ids: ids,
	})
	if err != nil {
		return nil, err
	}
	return response, nil
}

// Search products
func (client *CatalogClient) SearchProducts(ctx context.Context, query string) (*pb.SearchProductsResponse, error) {
	response, err := client.client.SearchProducts(ctx, &pb.SearchProductsRequest{
		Query: query,
	})
	if err != nil {
		return nil, err
	}
	return response, nil
}
