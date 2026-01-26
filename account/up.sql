CREATE TABLE IF NOT EXISTS accounts (
    id CHAR(27) PRIMARY KEY,
    name VARCHAR(24),
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS sessions (
    id CHAR(27) PRIMARY KEY,
    account_id CHAR(27) NOT NULL,
    access_token VARCHAR(255) NOT NULL,
    refresh_token VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL,
    is_revoked BOOLEAN NOT NULL DEFAULT FALSE,
    FOREIGN KEY (account_id) REFERENCES accounts(id)
);

CREATE TABLE IF NOT EXISTS device_info (
    id CHAR(27) PRIMARY KEY,
    session_id CHAR(27) NOT NULL,
    device_type enum('mobile', 'desktop', 'other') NOT NULL,
    device_model VARCHAR(255) NOT NULL,
    device_os VARCHAR(255) NOT NULL,
    device_os_version VARCHAR(255),
    device_id VARCHAR(255) NOT NULL,
    ip_address VARCHAR(255) NOT NULL,
    user_agent TEXT,
    created_at TIMESTAMP NOT NULL,
    FOREIGN KEY (session_id) REFERENCES sessions(id)
);