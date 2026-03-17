CREATE TABLE IF NOT EXISTS orders (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    status VARCHAR(32)
);

CREATE TABLE IF NOT EXISTS order_contents (
    order_id BIGINT REFERENCES orders (id) ON DELETE RESTRICT,
    product_id BIGINT NOT NULL,
    price_snapshot BIGINT NOT NULL
);
