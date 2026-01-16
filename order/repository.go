package order

import (
	"context"
	"database/sql"
)

type Repository interface {
	Close()
	CreateOrUpdateOrder(ctx context.Context, order *Order) (*Order, error)
	GetOrderById(ctx context.Context, id string) (*Order, error)
	GetOrdersForAccount(ctx context.Context, accountId string) ([]*Order, error)
}

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(url string) (Repository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &PostgresRepository{db}, nil
}

func (repository *PostgresRepository) Close() {
	repository.db.Close()
}
