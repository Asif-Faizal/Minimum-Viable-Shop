package graphql

import (
	"context"
	"fmt"
)

type queryResolver struct {
	server *Server
}

// Accounts retrieves accounts with optional pagination and filtering
func (r *queryResolver) Accounts(ctx context.Context, pagination *PaginationInput, id *string) ([]*Account, error) {
	// If specific ID requested, get single account
	if id != nil {
		accountResp, err := r.server.accountClient.GetAccountByID(ctx, *id)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch account: %w", err)
		}
		if accountResp == nil || accountResp.Account == nil {
			return nil, fmt.Errorf("account not found: %s", *id)
		}
		return []*Account{
			{
				ID:       accountResp.Account.Id,
				Name:     accountResp.Account.Name,
				UserType: accountResp.Account.Usertype,
				Email:    accountResp.Account.Email,
			},
		}, nil
	}

	// List accounts with pagination
	skip := uint32(0)
	take := uint32(10)

	if pagination != nil {
		if pagination.Skip != nil {
			skip = uint32(*pagination.Skip)
		}
		if pagination.Take != nil {
			take = uint32(*pagination.Take)
		}
	}

	accountsResp, err := r.server.accountClient.ListAccounts(ctx, skip, take)
	if err != nil {
		return nil, fmt.Errorf("failed to list accounts: %w", err)
	}

	accounts := make([]*Account, 0, len(accountsResp.Accounts))
	for _, account := range accountsResp.Accounts {
		accounts = append(accounts, &Account{
			ID:       account.Id,
			Name:     account.Name,
			UserType: account.Usertype,
			Email:    account.Email,
		})
	}
	return accounts, nil
}

// Products retrieves products with optional pagination, filtering by ID, and search query
func (r *queryResolver) Products(ctx context.Context, pagination *PaginationInput, id *string, query *string) ([]*Product, error) {
	// If specific ID requested, get single product
	if id != nil {
		productResp, err := r.server.catalogClient.GetProductByID(ctx, *id)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch product: %w", err)
		}
		if productResp == nil || productResp.Product == nil {
			return nil, fmt.Errorf("product not found: %s", *id)
		}
		return []*Product{
			{
				ID:          productResp.Product.Id,
				Name:        productResp.Product.Name,
				Description: productResp.Product.Description,
				Price:       float64(productResp.Product.Price),
			},
		}, nil
	}

	// If search query provided, use search
	if query != nil {
		searchResp, err := r.server.catalogClient.SearchProducts(ctx, *query)
		if err != nil {
			return nil, fmt.Errorf("failed to search products: %w", err)
		}
		products := make([]*Product, 0, len(searchResp.Products))
		for _, p := range searchResp.Products {
			products = append(products, &Product{
				ID:          p.Id,
				Name:        p.Name,
				Description: p.Description,
				Price:       float64(p.Price),
			})
		}
		return products, nil
	}

	// List products with pagination
	skip := uint64(0)
	take := uint64(10)

	if pagination != nil {
		if pagination.Skip != nil {
			skip = uint64(*pagination.Skip)
		}
		if pagination.Take != nil {
			take = uint64(*pagination.Take)
		}
	}

	productsResp, err := r.server.catalogClient.ListProducts(ctx, skip, take)
	if err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}

	products := make([]*Product, 0, len(productsResp.Products))
	for _, p := range productsResp.Products {
		products = append(products, &Product{
			ID:          p.Id,
			Name:        p.Name,
			Description: p.Description,
			Price:       float64(p.Price),
		})
	}
	return products, nil
}

// Order retrieves a single order by ID
func (r *queryResolver) Order(ctx context.Context, id string) (*Order, error) {
	orderResp, err := r.server.orderClient.GetOrderByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch order: %w", err)
	}
	if orderResp == nil || orderResp.Order == nil {
		return nil, fmt.Errorf("order not found: %s", id)
	}

	orderedProducts := make([]*OrderedProduct, 0, len(orderResp.Order.Products))
	for _, p := range orderResp.Order.Products {
		orderedProducts = append(orderedProducts, &OrderedProduct{
			ID:          p.ProductId,
			Name:        p.ProductName,
			Description: p.ProductDescription,
			Price:       p.Price,
			Quantity:    int(p.Quantity),
		})
	}
	createdAt := orderResp.Order.CreatedAt.AsTime()
	return &Order{
		ID:         orderResp.Order.Id,
		CreatedAt:  createdAt,
		AccountID:  orderResp.Order.AccountId,
		TotalPrice: orderResp.Order.TotalPrice,
		Products:   orderedProducts,
	}, nil
}

// OrdersForAccount retrieves all orders for a given account
func (r *queryResolver) OrdersForAccount(ctx context.Context, accountID string) ([]*Order, error) {
	ordersResp, err := r.server.orderClient.ListOrdersForAccount(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to list orders for account: %w", err)
	}

	orders := make([]*Order, 0, len(ordersResp.Orders))
	for _, o := range ordersResp.Orders {
		orderedProducts := make([]*OrderedProduct, 0, len(o.Products))
		for _, p := range o.Products {
			orderedProducts = append(orderedProducts, &OrderedProduct{
				ID:          p.ProductId,
				Name:        p.ProductName,
				Description: p.ProductDescription,
				Price:       p.Price,
				Quantity:    int(p.Quantity),
			})
		}
		createdAt := o.CreatedAt.AsTime()
		orders = append(orders, &Order{
			ID:         o.Id,
			CreatedAt:  createdAt,
			AccountID:  o.AccountId,
			TotalPrice: o.TotalPrice,
			Products:   orderedProducts,
		})
	}
	return orders, nil
}
