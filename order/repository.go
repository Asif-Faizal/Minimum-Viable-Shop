package order

import (
	"context"
	"database/sql"
	"time"
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
	var err error
	rows, err := repository.db.QueryContext(
		ctx,
		`SELECT
			o.id,
			o.createdAt,
			o.accountId,
			o.totalPrice,
			op.productId,
			op.quantity,
			p.name,
			p.price
		FROM orders o
		JOIN order_products op ON (o.id = op.orderId)
		JOIN products p ON (op.productId = p.id)
		WHERE o.accountId = $1
		ORDER BY o.createdAt DESC`,
		accountId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ordersMap := make(map[string]*Order)
	var orders []*Order

	for rows.Next() {
		var orderID string
		var createdAt time.Time
		var accountID string
		var totalPrice float64
		var productID string
		var quantity int
		var productName string
		var productPrice float64

		if err = rows.Scan(
			&orderID,
			&createdAt,
			&accountID,
			&totalPrice,
			&productID,
			&quantity,
			&productName,
			&productPrice,
		); err != nil {
			return nil, err
		}

		order, ok := ordersMap[orderID]
		if !ok {
			order = &Order{
				ID:         orderID,
				CreatedAt:  createdAt,
				AccountID:  accountID,
				TotalPrice: totalPrice,
				Products:   []*OrderProduct{},
			}
			ordersMap[orderID] = order
			orders = append(orders, order)
		}

		order.Products = append(order.Products, &OrderProduct{
			OrderID:     orderID,
			ProductID:   productID,
			ProductName: productName,
			Price:       productPrice,
			Quantity:    quantity,
		})
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}
