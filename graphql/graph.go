package graphql

import (
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/Asif-Faizal/Minimum-Viable-Shop/account"
	"github.com/Asif-Faizal/Minimum-Viable-Shop/catalog"
	"github.com/Asif-Faizal/Minimum-Viable-Shop/order"
)

type Server struct {
	accountClient *account.AccountClient
	catalogClient *catalog.CatalogClient
	orderClient   *order.OrderClient
}

// NewGraphQLServer initializes and returns a new GraphQL server with all microservice clients
func NewGraphQLServer(accountURL, catalogURL, orderURL string) (*Server, error) {
	// Validate URLs
	if accountURL == "" || catalogURL == "" || orderURL == "" {
		return nil, fmt.Errorf("all service URLs must be provided")
	}

	// Initialize Account Client
	accountClient, err := account.NewAccountClient(accountURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to account service: %w", err)
	}

	// Initialize Catalog Client
	catalogClient, err := catalog.NewCatalogClient(catalogURL)
	if err != nil {
		accountClient.Close()
		return nil, fmt.Errorf("failed to connect to catalog service: %w", err)
	}

	// Initialize Order Client
	orderClient, err := order.NewOrderClient(orderURL)
	if err != nil {
		accountClient.Close()
		catalogClient.Close()
		return nil, fmt.Errorf("failed to connect to order service: %w", err)
	}

	return &Server{
		accountClient: accountClient,
		catalogClient: catalogClient,
		orderClient:   orderClient,
	}, nil
}

// Close gracefully closes all client connections
func (s *Server) Close() error {
	if s.accountClient != nil {
		s.accountClient.Close()
	}
	if s.catalogClient != nil {
		s.catalogClient.Close()
	}
	if s.orderClient != nil {
		s.orderClient.Close()
	}
	return nil
}
func (s *Server) Mutation() MutationResolver {
	return &mutationResolver{
		server: s,
	}
}

func (s *Server) Query() QueryResolver {
	return &queryResolver{
		server: s,
	}
}

func (s *Server) Account() AccountResolver {
	return &accountResolver{
		server: s,
	}
}

func (s *Server) ToExecutableSchema() graphql.ExecutableSchema {
	return NewExecutableSchema(Config{
		Resolvers: s,
	})
}
