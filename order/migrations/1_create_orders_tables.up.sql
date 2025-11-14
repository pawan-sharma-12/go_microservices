-- Orders table
CREATE TABLE IF NOT EXISTS orders (
    id CHAR(27) PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    account_id CHAR(27) NOT NULL,
    total_price DOUBLE PRECISION NOT NULL
);

-- Order products table
CREATE TABLE IF NOT EXISTS order_products (
    order_id CHAR(27) NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id CHAR(27) NOT NULL,
    quantity BIGINT NOT NULL,
    PRIMARY KEY (order_id, product_id)
);
