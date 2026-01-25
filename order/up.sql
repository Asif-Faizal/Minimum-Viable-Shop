CREATE TABLE IF NOT EXISTS orders (
    id CHAR(27) PRIMARY KEY,
    createdAt TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    accountId CHAR(27) NOT NULL,
    totalPrice NUMERIC(19,4) NOT NULL
);

CREATE TABLE IF NOT EXISTS order_products (
  orderId CHAR(27) REFERENCES orders (id) ON DELETE CASCADE,
  productId CHAR(27),
  quantity INT NOT NULL,
  name VARCHAR(255) NOT NULL,
  description TEXT,
  price NUMERIC(19,4) NOT NULL,
  PRIMARY KEY (productId, orderId)
);