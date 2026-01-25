package catalog

import (
	"context"
	"encoding/json"

	"github.com/olivere/elastic/v7"
)

type Repository interface {
	Close()
	CreateOrUpdateProduct(ctx context.Context, product *Product) (*Product, error)
	GetProductById(ctx context.Context, id string) (*Product, error)
	ListProducts(ctx context.Context, skip uint64, take uint64) ([]*Product, error)
	ListProductsWithIds(ctx context.Context, ids []string) ([]*Product, error)
	SearchProducts(ctx context.Context, query string, skip uint64, take uint64) ([]*Product, error)
}

type ElasticRepository struct {
	client *elastic.Client
}

func NewElasticRepository(url string) (Repository, error) {
	client, err := elastic.NewClient(elastic.SetURL(url), elastic.SetSniff(false))
	if err != nil {
		return nil, err
	}
	return &ElasticRepository{client: client}, nil
}

func (repository *ElasticRepository) Close() {
	repository.client.Stop()
}

func (repository *ElasticRepository) CreateOrUpdateProduct(ctx context.Context, product *Product) (*Product, error) {
	_, err := repository.client.Index().
		Index("catalog").
		Id(product.ID).
		BodyJson(ProductDocument{
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
		}).
		Do(ctx)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (repository *ElasticRepository) GetProductById(ctx context.Context, id string) (*Product, error) {
	res, err := repository.client.Get().
		Index("catalog").
		Id(id).
		Do(ctx)
	if err != nil {
		return nil, err
	}
	product := ProductDocument{}
	if err := json.Unmarshal(res.Source, &product); err != nil {
		return nil, err
	}
	return &Product{
		ID:          id,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
	}, nil
}

func (repository *ElasticRepository) ListProducts(ctx context.Context, skip uint64, take uint64) ([]*Product, error) {
	res, err := repository.client.Search().
		Index("catalog").
		Query(elastic.NewMatchAllQuery()).
		From(int(skip)).
		Size(int(take)).
		Do(ctx)
	if err != nil {
		return nil, err
	}
	products := []*Product{}
	for _, hit := range res.Hits.Hits {
		product := ProductDocument{}
		if err := json.Unmarshal(hit.Source, &product); err != nil {
			return nil, err
		}
		products = append(products, &Product{
			ID:          hit.Id,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
		})
	}
	return products, nil
}

func (repository *ElasticRepository) ListProductsWithIds(ctx context.Context, ids []string) ([]*Product, error) {
	items := []*elastic.MultiGetItem{}
	for _, id := range ids {
		items = append(
			items,
			elastic.NewMultiGetItem().
				Index("catalog").
				Id(id),
		)
	}
	res, err := repository.client.MultiGet().
		Add(items...).
		Do(ctx)
	if err != nil {
		return nil, err
	}
	products := []*Product{}
	for _, doc := range res.Docs {
		product := ProductDocument{}
		if err = json.Unmarshal(doc.Source, &product); err == nil {
			products = append(products, &Product{
				ID:          doc.Id,
				Name:        product.Name,
				Description: product.Description,
				Price:       product.Price,
			})
		}
	}
	return products, nil
}

func (repository *ElasticRepository) SearchProducts(ctx context.Context, query string, skip uint64, take uint64) ([]*Product, error) {
	res, err := repository.client.Search().
		Index("catalog").
		Query(elastic.NewMultiMatchQuery(query, "name", "description")).
		From(int(skip)).Size(int(take)).
		Do(ctx)
	if err != nil {
		return nil, err
	}
	products := []*Product{}
	for _, hit := range res.Hits.Hits {
		product := ProductDocument{}
		if err = json.Unmarshal(hit.Source, &product); err == nil {
			products = append(products, &Product{
				ID:          hit.Id,
				Name:        product.Name,
				Description: product.Description,
				Price:       product.Price,
			})
		}
	}
	return products, err
}
