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

	// Session Management
	CreateOrUpdateSession(ctx context.Context, session *Session) error
	GetSession(ctx context.Context, id string) (*Session, error)
	GetSessionByRefreshToken(ctx context.Context, refreshToken string) (*Session, error)
	GetSessionByAccessToken(ctx context.Context, accessToken string) (*Session, error)
	RevokeSessionByAccessToken(ctx context.Context, accessToken string) error

	// Device Info
	CreateOrUpdateDeviceInfo(ctx context.Context, info *DeviceInfo) error
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

func (repository *PostgresRepository) CreateOrUpdateSession(ctx context.Context, session *Session) error {
	start := time.Now()
	query := `
		INSERT INTO sessions (id, account_id, device_id, access_token, refresh_token, expires_at, created_at, is_revoked)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (account_id, device_id) 
		DO UPDATE SET 
			access_token = EXCLUDED.access_token,
			refresh_token = EXCLUDED.refresh_token,
			expires_at = EXCLUDED.expires_at,
			is_revoked = EXCLUDED.is_revoked
		RETURNING id
	`

	err := repository.db.QueryRowContext(ctx, query,
		session.ID,
		session.AccountID,
		session.DeviceID,
		session.AccessToken,
		session.RefreshToken,
		session.ExpiresAt,
		session.CreatedAt,
		session.IsRevoked,
	).Scan(&session.ID)

	repository.logger.Database().Debug().
		Str("query", query).
		Str("duration", time.Since(start).String()).
		Bool("success", err == nil).
		Msg("Execute Query")

	return err
}

func (repository *PostgresRepository) GetSession(ctx context.Context, id string) (*Session, error) {
	start := time.Now()
	query := "SELECT id, account_id, device_id, access_token, refresh_token, expires_at, created_at, is_revoked FROM sessions WHERE id = $1"

	row := repository.db.QueryRowContext(ctx, query, id)
	session := &Session{}
	err := row.Scan(&session.ID, &session.AccountID, &session.DeviceID, &session.AccessToken, &session.RefreshToken, &session.ExpiresAt, &session.CreatedAt, &session.IsRevoked)

	repository.logger.Database().Debug().
		Str("query", query).
		Str("duration", time.Since(start).String()).
		Bool("success", err == nil).
		Msg("Query Row")

	if err != nil {
		return nil, err
	}
	return session, nil
}

func (repository *PostgresRepository) GetSessionByRefreshToken(ctx context.Context, refreshToken string) (*Session, error) {
	start := time.Now()
	query := "SELECT id, account_id, device_id, access_token, refresh_token, expires_at, created_at, is_revoked FROM sessions WHERE refresh_token = $1 AND is_revoked = false"

	row := repository.db.QueryRowContext(ctx, query, refreshToken)
	session := &Session{}
	err := row.Scan(&session.ID, &session.AccountID, &session.DeviceID, &session.AccessToken, &session.RefreshToken, &session.ExpiresAt, &session.CreatedAt, &session.IsRevoked)

	repository.logger.Database().Debug().
		Str("query", query).
		Str("duration", time.Since(start).String()).
		Bool("success", err == nil).
		Msg("Query Row")

	if err != nil {
		return nil, err
	}
	return session, nil
}

func (repository *PostgresRepository) GetSessionByAccessToken(ctx context.Context, accessToken string) (*Session, error) {
	start := time.Now()
	query := "SELECT id, account_id, device_id, access_token, refresh_token, expires_at, created_at, is_revoked FROM sessions WHERE access_token = $1"

	row := repository.db.QueryRowContext(ctx, query, accessToken)
	session := &Session{}
	err := row.Scan(&session.ID, &session.AccountID, &session.DeviceID, &session.AccessToken, &session.RefreshToken, &session.ExpiresAt, &session.CreatedAt, &session.IsRevoked)

	repository.logger.Database().Debug().
		Str("query", query).
		Str("duration", time.Since(start).String()).
		Bool("success", err == nil).
		Msg("Query Row")

	if err != nil {
		return nil, err
	}
	return session, nil
}

func (repository *PostgresRepository) RevokeSessionByAccessToken(ctx context.Context, accessToken string) error {
	start := time.Now()
	query := "UPDATE sessions SET is_revoked = true WHERE access_token = $1"

	_, err := repository.db.ExecContext(ctx, query, accessToken)

	repository.logger.Database().Debug().
		Str("query", query).
		Str("duration", time.Since(start).String()).
		Bool("success", err == nil).
		Msg("Execute Query")

	return err
}

func (repository *PostgresRepository) CreateOrUpdateDeviceInfo(ctx context.Context, info *DeviceInfo) error {
	start := time.Now()
	query := `
		INSERT INTO device_info (id, session_id, device_type, device_model, device_os, device_os_version, ip_address, user_agent, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (session_id) 
		DO UPDATE SET 
			device_type = EXCLUDED.device_type,
			device_model = EXCLUDED.device_model,
			device_os = EXCLUDED.device_os,
			device_os_version = EXCLUDED.device_os_version,
			ip_address = EXCLUDED.ip_address,
			user_agent = EXCLUDED.user_agent
	`

	_, err := repository.db.ExecContext(ctx, query,
		info.ID,
		info.SessionID,
		info.DeviceType,
		info.DeviceModel,
		info.DeviceOS,
		info.DeviceOSVersion,
		info.IPAddress,
		info.UserAgent,
		info.CreatedAt,
	)

	repository.logger.Database().Debug().
		Str("query", query).
		Str("duration", time.Since(start).String()).
		Bool("success", err == nil).
		Msg("Execute Query")

	return err
}
