package order

import (
	"context"
	"fmt"
	"net"

	"github.com/Asif-Faizal/Minimum-Viable-Shop/account"
	"github.com/Asif-Faizal/Minimum-Viable-Shop/catalog"
	"github.com/Asif-Faizal/Minimum-Viable-Shop/order/pb/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"
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

func (server *GrpcServer) CreateOrUpdateOrder(ctx context.Context, request *pb.CreateOrUpdateOrderRequest) (*pb.CreateOrUpdateOrderResponse, error) {
	// 1. Check if account exists
	accountResponse, err := server.accountClient.GetAccountByID(ctx, request.Order.AccountId)
	if err != nil {
		return nil, err
	}
	if accountResponse == nil || accountResponse.Account == nil {
		return nil, err
	}

	// 2. Get ordered products from catalog
	productIDs := []string{}
	for _, p := range request.Order.Products {
		productIDs = append(productIDs, p.ProductId)
	}
	catalogResp, err := server.catalogClient.ListProductsWithIDs(ctx, productIDs)
	if err != nil {
		return nil, err
	}

	// 3. Construct domain products
	products := []*OrderProduct{}
	for _, p := range catalogResp.Products {
		product := &OrderProduct{
			ProductID:          p.Id,
			ProductName:        p.Name,
			ProductDescription: p.Description,
			Price:              float64(p.Price),
			Quantity:           0,
		}
		for _, rp := range request.Order.Products {
			if rp.ProductId == p.Id {
				product.Quantity = rp.Quantity
				break
			}
		}

		if product.Quantity != 0 {
			products = append(products, product)
		}
	}

	if len(products) != len(request.Order.Products) {
		return nil, fmt.Errorf("one or more products not found in catalog")
	}

	// 4. Call service implementation
	domainOrder := &Order{
		ID:        request.Order.Id,
		AccountID: request.Order.AccountId,
		Products:  products,
	}
	if request.Order.CreatedAt != nil {
		domainOrder.CreatedAt = request.Order.CreatedAt.AsTime()
	}

	order, err := server.orderService.CreateOrUpdateOrder(ctx, domainOrder)
	if err != nil {
		return nil, err
	}

	// 5. Make response order
	pbProducts := []*pb.OrderProduct{}
	for _, p := range order.Products {
		pbProducts = append(pbProducts, &pb.OrderProduct{
			OrderId:            p.OrderID,
			ProductId:          p.ProductID,
			ProductName:        p.ProductName,
			ProductDescription: p.ProductDescription,
			Price:              p.Price,
			Quantity:           p.Quantity,
		})
	}

	return &pb.CreateOrUpdateOrderResponse{
		Order: &pb.Order{
			Id:         order.ID,
			AccountId:  order.AccountID,
			TotalPrice: order.TotalPrice,
			Products:   pbProducts,
			CreatedAt:  timestamppb.New(order.CreatedAt),
		},
	}, nil
}

func (server *GrpcServer) GetOrderByID(ctx context.Context, request *pb.GetOrderByIDRequest) (*pb.GetOrderByIDResponse, error) {
	order, err := server.orderService.GetOrderById(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	pbProducts := []*pb.OrderProduct{}
	for _, p := range order.Products {
		pbProducts = append(pbProducts, &pb.OrderProduct{
			OrderId:            p.OrderID,
			ProductId:          p.ProductID,
			ProductName:        p.ProductName,
			ProductDescription: p.ProductDescription,
			Price:              p.Price,
			Quantity:           p.Quantity,
		})
	}

	return &pb.GetOrderByIDResponse{
		Order: &pb.Order{
			Id:         order.ID,
			AccountId:  order.AccountID,
			TotalPrice: order.TotalPrice,
			Products:   pbProducts,
			CreatedAt:  timestamppb.New(order.CreatedAt),
		},
	}, nil
}

func (server *GrpcServer) GetOrdersForAccount(ctx context.Context, request *pb.GetOrdersForAccountRequest) (*pb.GetOrdersForAccountResponse, error) {
	orders, err := server.orderService.GetOrdersForAccount(ctx, request.AccountId)
	if err != nil {
		return nil, err
	}

	pbOrders := []*pb.Order{}
	for _, o := range orders {
		pbProducts := []*pb.OrderProduct{}
		for _, p := range o.Products {
			pbProducts = append(pbProducts, &pb.OrderProduct{
				OrderId:            p.OrderID,
				ProductId:          p.ProductID,
				ProductName:        p.ProductName,
				ProductDescription: p.ProductDescription,
				Price:              p.Price,
				Quantity:           p.Quantity,
			})
		}
		pbOrders = append(pbOrders, &pb.Order{
			Id:         o.ID,
			AccountId:  o.AccountID,
			TotalPrice: o.TotalPrice,
			Products:   pbProducts,
			CreatedAt:  timestamppb.New(o.CreatedAt),
		})
	}

	return &pb.GetOrdersForAccountResponse{
		Orders: pbOrders,
	}, nil
}
