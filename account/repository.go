package account

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
)

type Repository interface {
	Close()
	CreateOrUpdateAccount(ctx context.Context, account *Account) error
	GetAccountByID(ctx context.Context, id string) (*Account, error)
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

func (repository *PostgresRepository) Close() {
	repository.db.Close()
}

func (repository *PostgresRepository) Ping() error {
	return repository.db.Ping()
}

func (repository *PostgresRepository) CreateOrUpdateAccount(ctx context.Context, account *Account) error {
	_, err := repository.db.ExecContext(ctx, "INSERT INTO accounts (id, name, email, password) VALUES ($1, $2, $3, $4) ON CONFLICT (id) DO UPDATE SET name = $2, email = $3, password = COALESCE(NULLIF($4, ''), accounts.password)", account.ID, account.Name, account.Email, account.Password)
	return err
}

func (repository *PostgresRepository) GetAccountByID(ctx context.Context, id string) (*Account, error) {
	row := repository.db.QueryRowContext(ctx, "SELECT id, name, email FROM accounts WHERE id = $1", id)
	account := &Account{}
	if err := row.Scan(&account.ID, &account.Name, &account.Email); err != nil {
		return nil, err
	}
	return account, nil
}

func (repository *PostgresRepository) ListAccounts(ctx context.Context, skip uint, take uint) ([]*Account, error) {
	rows, err := repository.db.QueryContext(ctx, "SELECT id, name, email FROM accounts ORDER by id DESC  LIMIT $1 OFFSET $2", take, skip)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	accounts := []*Account{}
	for rows.Next() {
		account := &Account{}
		if err := rows.Scan(&account.ID, &account.Name, &account.Email); err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return accounts, nil
}
