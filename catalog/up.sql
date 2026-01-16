CREATE TABLE IF NOT EXISTS products (
    id CHAR(27) PRIMARY KEY,
    name VARCHAR(24) NOT NULL,
    description TEXT,
    price MONEY NOT NULL
);