package account

import (
	"context"
	"database/sql"
	"time"

	"github.com/Asif-Faizal/Minimum-Viable-Shop/util"
	_ "github.com/lib/pq"
)

type Repository interface {
	Close()
	CreateOrUpdateAccount(ctx context.Context, account *Account) (*Account, error)
	GetAccountById(ctx context.Context, id string) (*Account, error)
	ListAccounts(ctx context.Context, skip uint, take uint) ([]*Account, error)
	CheckEmailExists(ctx context.Context, email string) (bool, error)
	GetAccountByEmail(ctx context.Context, email string) (*Account, error)
	CreateSession(ctx context.Context, session *Session) error
}

type PostgresRepository struct {
	db     *sql.DB
	logger util.Logger
}

func NewPostgresRepository(url string, logger util.Logger) (Repository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &PostgresRepository{
		db:     db,
		logger: logger,
	}, nil
}

func (repository *PostgresRepository) Close() {
	repository.db.Close()
}

func (repository *PostgresRepository) CreateOrUpdateAccount(ctx context.Context, account *Account) (*Account, error) {
	start := time.Now()
	query := "INSERT INTO accounts (id, name, email, password) VALUES ($1, NULLIF($2, ''), $3, $4) ON CONFLICT (id) DO UPDATE SET name = NULLIF($2, ''), email = $3, password = $4"

	_, err := repository.db.ExecContext(ctx, query, account.ID, account.Name, account.Email, account.Password)

	repository.logger.Database().Debug().
		Str("query", query).
		Str("duration", time.Since(start).String()).
		Bool("success", err == nil).
		Msg("Execute Query")

	if err != nil {
		return nil, err
	}
	return account, nil
}

func (repository *PostgresRepository) GetAccountById(ctx context.Context, id string) (*Account, error) {
	start := time.Now()
	query := "SELECT id, name, email FROM accounts WHERE id = $1"

	row := repository.db.QueryRowContext(ctx, query, id)
	account := &Account{}
	var name sql.NullString
	err := row.Scan(&account.ID, &name, &account.Email)

	repository.logger.Database().Debug().
		Str("query", query).
		Str("duration", time.Since(start).String()).
		Bool("success", err == nil).
		Msg("Query Row")

	if err != nil {
		return nil, err
	}
	account.Name = name.String
	return account, nil
}

func (repository *PostgresRepository) ListAccounts(ctx context.Context, skip uint, take uint) ([]*Account, error) {
	start := time.Now()
	query := "SELECT id, name, email FROM accounts ORDER by id DESC LIMIT $1 OFFSET $2"

	rows, err := repository.db.QueryContext(ctx, query, take, skip)

	repository.logger.Database().Debug().
		Str("query", query).
		Str("duration", time.Since(start).String()).
		Bool("success", err == nil).
		Msg("Query Context")

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	accounts := []*Account{}
	for rows.Next() {
		account := &Account{}
		var name sql.NullString
		if err := rows.Scan(&account.ID, &name, &account.Email); err != nil {
			return nil, err
		}
		account.Name = name.String
		accounts = append(accounts, account)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return accounts, nil
}

func (repository *PostgresRepository) CheckEmailExists(ctx context.Context, email string) (bool, error) {
	start := time.Now()
	query := "SELECT EXISTS(SELECT 1 FROM accounts WHERE email = $1)"

	row := repository.db.QueryRowContext(ctx, query, email)
	var exists bool
	err := row.Scan(&exists)

	repository.logger.Database().Debug().
		Str("query", query).
		Str("duration", time.Since(start).String()).
		Bool("success", err == nil).
		Msg("Query Row")

	if err != nil {
		return false, err
	}
	return exists, nil
}

func (repository *PostgresRepository) GetAccountByEmail(ctx context.Context, email string) (*Account, error) {
	start := time.Now()
	query := "SELECT id, name, email, password FROM accounts WHERE email = $1"

	row := repository.db.QueryRowContext(ctx, query, email)
	account := &Account{}
	var name sql.NullString
	err := row.Scan(&account.ID, &name, &account.Email, &account.Password)

	repository.logger.Database().Debug().
		Str("query", query).
		Str("duration", time.Since(start).String()).
		Bool("success", err == nil).
		Msg("Query Row")

	if err != nil {
		return nil, err
	}
	account.Name = name.String
	return account, nil
}

func (repository *PostgresRepository) CreateSession(ctx context.Context, session *Session) error {
	start := time.Now()
	query := "INSERT INTO sessions (id, account_id, access_token, refresh_token, expires_at, created_at, is_revoked) VALUES ($1, $2, $3, $4, $5, $6, $7)"

	_, err := repository.db.ExecContext(ctx, query,
		session.ID,
		session.AccountID,
		session.AccessToken,
		session.RefreshToken,
		session.ExpiresAt,
		session.CreatedAt,
		session.IsRevoked,
	)

	repository.logger.Database().Debug().
		Str("query", query).
		Str("duration", time.Since(start).String()).
		Bool("success", err == nil).
		Msg("Execute Query")

	return err
}
