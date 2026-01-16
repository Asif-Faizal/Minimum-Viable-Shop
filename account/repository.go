package account

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
)

type Repository interface {
	Close()
	CreateAccount(ctx context.Context, account *Account) error
	GetAccountbyID(ctx context.Context, id string) (*Account, error)
	ListAccounts(ctx context.Context, skip uint, take uint) ([]*Account, error)
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

func (r *PostgresRepository) Close() {
	r.db.Close()
}

func (r *PostgresRepository) Ping() error {
	return r.db.Ping()
}

func (r *PostgresRepository) CreateAccount(ctx context.Context, account *Account) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO accounts (id, name, email, password) VALUES ($1, $2, $3, $4)", account.ID, account.Name, account.Email, account.Password)
	return err
}

func (r *PostgresRepository) GetAccountbyID(ctx context.Context, id string) (*Account, error) {
	r.db.QueryRowContext(ctx, "SELECT id, name, email, password FROM accounts WHERE id = $1", id)
	account := &Account{}
	if err := row.Scan(&account.ID, &account.Name, &account.Email, &account.Password); err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *PostgresRepository) ListAccounts(ctx context.Context, skip uint, take uint) ([]*Account, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, name, email, password FROM accounts ORDER by id DESC  LIMIT $1 OFFSET $2", take, skip)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	accounts := []*Account{}
	for rows.Next() {
		account := &Account{}
		if err := rows.Scan(&account.ID, &account.Name, &account.Email, &account.Password); err != nil {
			return nil, err
		}
		accounts = append(accounts, *account)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return accounts, nil
}
