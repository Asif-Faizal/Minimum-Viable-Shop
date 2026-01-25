package graphql

import (
	"context"
	"fmt"
)

type accountResolver struct {
	server *Server
}

// Orders retrieves all orders for a given account
func (resolver *accountResolver) Orders(ctx context.Context, account *Account) ([]*Order, error) {
	if account == nil || account.ID == "" {
		return nil, fmt.Errorf("account id is required")
	}

	// Call order service to get orders for account
	resp, err := resolver.server.orderClient.ListOrdersForAccount(ctx, account.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch orders for account %s: %w", account.ID, err)
	}

	if resp == nil {
		return []*Order{}, nil
	}

	// Convert proto orders to GraphQL orders
	orders := make([]*Order, 0, len(resp.Orders))
	for _, order := range resp.Orders {
		// Convert products
		orderedProducts := make([]*OrderedProduct, 0, len(order.Products))
		for _, p := range order.Products {
			orderedProducts = append(orderedProducts, &OrderedProduct{
				ID:          p.ProductId,
				Name:        p.ProductName,
				Description: p.ProductDescription,
				Price:       p.Price,
				Quantity:    int(p.Quantity),
			})
		}

		createdAt := order.CreatedAt.AsTime()
		orders = append(orders, &Order{
			ID:         order.Id,
			CreatedAt:  createdAt,
			TotalPrice: order.TotalPrice,
			Products:   orderedProducts,
		})
	}

	return orders, nil
}
