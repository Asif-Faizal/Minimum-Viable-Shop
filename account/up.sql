CREATE TABLE IF NOT EXISTS accounts (
    id CHAR(27) PRIMARY KEY,
    name VARCHAR(24),
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS sessions (
    id CHAR(27) PRIMARY KEY,
    account_id CHAR(27) NOT NULL,
    device_id VARCHAR(255) NOT NULL,
    access_token VARCHAR(255) NOT NULL,
    refresh_token VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL,
    is_revoked BOOLEAN NOT NULL DEFAULT FALSE,
    UNIQUE (account_id, device_id),
    FOREIGN KEY (account_id) REFERENCES accounts(id)
);

DO $$ BEGIN
    CREATE TYPE device_type_enum AS ENUM ('mobile', 'desktop', 'other');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

CREATE TABLE IF NOT EXISTS device_info (
    id CHAR(27) PRIMARY KEY,
    session_id CHAR(27) NOT NULL UNIQUE,
    device_type device_type_enum NOT NULL,
    device_model VARCHAR(255) NOT NULL,
    device_os VARCHAR(255) NOT NULL,
    device_os_version VARCHAR(255),
    ip_address VARCHAR(255) NOT NULL,
    user_agent TEXT,
    created_at TIMESTAMP NOT NULL,
    FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE CASCADE
);

-- Active Refresh Token lookup (for Token Refresh)
CREATE INDEX IF NOT EXISTS idx_sessions_refresh_token_active ON sessions (refresh_token) WHERE is_revoked = FALSE;

-- Access Token lookup (for Logout)
CREATE INDEX IF NOT EXISTS idx_sessions_access_token ON sessions (access_token);