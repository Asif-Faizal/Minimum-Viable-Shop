package catalog

import (
	"context"

	"github.com/segmentio/ksuid"
)

type Service struct {
	CreateOrUpdateProduct func(ctx context.Context, product *Product) (*Product, error)
	GetProductById        func(ctx context.Context, id string) (*Product, error)
	ListProducts          func(ctx context.Context, skip uint64, take uint64) ([]*Product, error)
	ListProductsWithIds   func(ctx context.Context, ids []string) ([]*Product, error)
	SearchProducts        func(ctx context.Context, query string, skip uint64, take uint64) ([]*Product, error)
}

type CatalogService struct {
	repository Repository
}

func NewCatalogService(repository Repository) *CatalogService {
	return &CatalogService{repository: repository}
}

func (service *CatalogService) CreateOrUpdateProduct(ctx context.Context, product *Product) (*Product, error) {
	newProduct := &Product{
		ID:          ksuid.New().String(),
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
	}
	if _, err := service.repository.CreateOrUpdateProduct(ctx, newProduct); err != nil {
		return nil, err
	}
	return newProduct, nil
}

func (service *CatalogService) GetProductById(ctx context.Context, id string) (*Product, error) {
	product, err := service.repository.GetProductById(ctx, id)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (service *CatalogService) ListProducts(ctx context.Context, skip uint64, take uint64) ([]*Product, error) {
	if take > 100 || (skip == 0 && take == 0) {
		take = 100
	}
	products, err := service.repository.ListProducts(ctx, skip, take)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (service *CatalogService) ListProductsWithIds(ctx context.Context, ids []string) ([]*Product, error) {
	products, err := service.repository.ListProductsWithIds(ctx, ids)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (service *CatalogService) SearchProducts(ctx context.Context, query string, skip uint64, take uint64) ([]*Product, error) {
	if take > 100 || (skip == 0 && take == 0) {
		take = 100
	}
	products, err := service.repository.SearchProducts(ctx, query, skip, take)
	if err != nil {
		return nil, err
	}
	return products, nil
}
