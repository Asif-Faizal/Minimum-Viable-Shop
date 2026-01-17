package order

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
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

func (repository *PostgresRepository) CreateOrUpdateOrder(ctx context.Context, order *Order) (*Order, error) {
	tx, err := repository.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()
	_, err = tx.ExecContext(ctx, "INSERT INTO orders (id, accountId,createdAt, totalPrice) VALUES ($1, $2, $3, $4) ON CONFLICT (id) DO UPDATE SET accountId = $2, totalPrice = $3", order.ID, order.AccountID, order.CreatedAt, order.TotalPrice)
	if err != nil {
		return nil, err
	}
	for _, product := range order.Products {
		_, err = tx.ExecContext(ctx, "INSERT INTO order_products (orderId, productId, quantity) VALUES ($1, $2, $3) ON CONFLICT (orderId, productId) DO UPDATE SET quantity = $3", order.ID, product.ProductID, product.Quantity)
		if err != nil {
			return nil, err
		}
	}
	return order, nil
}

func (repository *PostgresRepository) GetOrderById(ctx context.Context, id string) (*Order, error) {
	row := repository.db.QueryRowContext(ctx, "SELECT id, createdAt, accountId, totalPrice FROM orders WHERE id = $1", id)
	order := &Order{}
	if err := row.Scan(&order.ID, &order.CreatedAt, &order.AccountID, &order.TotalPrice); err != nil {
		return nil, err
	}
	return order, nil
}

func (repository *PostgresRepository) GetOrdersForAccount(ctx context.Context, accountId string) ([]*Order, error) {
	rows, err := repository.db.QueryContext(ctx, "SELECT id, createdAt, accountId, totalPrice FROM orders WHERE accountId = $1", accountId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := []*Order{}
	orderIDs := []string{}

	for rows.Next() {
		order := &Order{
			Products: []*OrderProduct{},
		}
		if err = rows.Scan(&order.ID, &order.CreatedAt, &order.AccountID, &order.TotalPrice); err != nil {
			return nil, err
		}
		orders = append(orders, order)
		orderIDs = append(orderIDs, order.ID)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	if len(orderIDs) == 0 {
		return orders, nil
	}

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	productRows, err := repository.db.QueryContext(ctx, "SELECT orderId, productId, quantity FROM order_products WHERE orderId = ANY($1)", pq.Array(orderIDs))
	if err != nil {
		return nil, err
	}
	defer productRows.Close()

	productsByOrder := make(map[string][]*OrderProduct, len(orderIDs))
	for productRows.Next() {
		product := &OrderProduct{}
		if err = productRows.Scan(&product.OrderID, &product.ProductID, &product.Quantity); err != nil {
			return nil, err
		}
		productsByOrder[product.OrderID] = append(productsByOrder[product.OrderID], product)
	}
	if err := productRows.Err(); err != nil {
		return nil, err
	}

	for _, order := range orders {
		if products, ok := productsByOrder[order.ID]; ok {
			order.Products = products
		}
	}

	return orders, nil
}
