package graphql

import (
	"context"
	"fmt"

	orderpb "github.com/Asif-Faizal/Minimum-Viable-Shop/order/pb/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type mutationResolver struct {
	server *Server
}

// CreateAccount creates or updates an account
func (r *mutationResolver) CreateAccount(ctx context.Context, input AccountInput) (*Account, error) {
	// Validate input
	if input.Name != nil && *input.Name == "" {
		return nil, fmt.Errorf("account name is required")
	}
	if input.Email == "" {
		return nil, fmt.Errorf("account email is required")
	}
	if input.Password == "" {
		return nil, fmt.Errorf("account password is required")
	}
	if input.UserType == "" {
		return nil, fmt.Errorf("account user_type is required")
	}

	// Determine ID: use provided ID or empty string for new account
	id := ""
	if input.ID != nil {
		id = *input.ID
	}

	// Call account service
	name := ""
	if input.Name != nil {
		name = *input.Name
	}
	resp, err := r.server.accountClient.CreateOrUpdateAccount(ctx, id, name, input.UserType, input.Email, input.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to create/update account: %w", err)
	}

	if resp == nil || resp.Account == nil {
		return nil, fmt.Errorf("unexpected response from account service")
	}

	return &Account{
		ID:       resp.Account.Id,
		Name:     resp.Account.Name,
		UserType: resp.Account.Usertype,
		Email:    resp.Account.Email,
	}, nil
}

// CreateProduct creates or updates a product in the catalog
func (r *mutationResolver) CreateProduct(ctx context.Context, input ProductInput) (*Product, error) {
	// Validate input
	if input.Name == "" {
		return nil, fmt.Errorf("product name is required")
	}
	if input.Price <= 0 {
		return nil, fmt.Errorf("product price must be positive")
	}

	// Determine ID: use provided ID or empty string for new product
	id := ""
	if input.ID != nil {
		id = *input.ID
	}

	// Call catalog service (convert float64 to float32)
	response, err := r.server.catalogClient.CreateOrUpdateProduct(ctx, id, input.Name, input.Description, float32(input.Price))
	if err != nil {
		return nil, fmt.Errorf("failed to create/update product: %w", err)
	}

	if response == nil || response.Product == nil {
		return nil, fmt.Errorf("unexpected response from catalog service")
	}

	return &Product{
		ID:          response.Product.Id,
		Name:        response.Product.Name,
		Description: response.Product.Description,
		Price:       float64(response.Product.Price),
	}, nil
}

// CreateOrder creates or updates an order
func (r *mutationResolver) CreateOrder(ctx context.Context, input OrderInput) (*Order, error) {
	// Validate input
	if input.AccountID == "" {
		return nil, fmt.Errorf("account_id is required")
	}
	if len(input.Products) == 0 {
		return nil, fmt.Errorf("order must contain at least one product")
	}

	// Convert input products to proto format
	protoProducts := make([]*orderpb.OrderProduct, 0, len(input.Products))
	for _, product := range input.Products {
		if product.ID == "" {
			return nil, fmt.Errorf("product id is required")
		}
		if product.Quantity <= 0 {
			return nil, fmt.Errorf("product quantity must be positive")
		}
		protoProducts = append(protoProducts, &orderpb.OrderProduct{
			ProductId: product.ID,
			Quantity:  int32(product.Quantity),
		})
	}

	// Determine ID: use provided ID or empty string for new order
	id := ""
	if input.ID != nil {
		id = *input.ID
	}

	// Call order service
	response, err := r.server.orderClient.CreateOrUpdateOrder(ctx, &orderpb.Order{
		Id:        id,
		AccountId: input.AccountID,
		Products:  protoProducts,
		CreatedAt: timestamppb.Now(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create/update order: %w", err)
	}

	if response == nil || response.Order == nil {
		return nil, fmt.Errorf("unexpected response from order service")
	}

	// Convert response to GraphQL type
	orderedProducts := make([]*OrderedProduct, 0, len(response.Order.Products))
	for _, product := range response.Order.Products {
		orderedProducts = append(orderedProducts, &OrderedProduct{
			ID:          product.ProductId,
			Name:        product.ProductName,
			Description: product.ProductDescription,
			Price:       product.Price,
			Quantity:    int(product.Quantity),
		})
	}

	createdAt := response.Order.CreatedAt.AsTime()
	return &Order{
		ID:         response.Order.Id,
		CreatedAt:  createdAt,
		TotalPrice: response.Order.TotalPrice,
		Products:   orderedProducts,
	}, nil
}
