CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(64) NOT NULL,
    password_hash VARCHAR(512) NOT NULL
);

CREATE TABLE IF NOT EXISTS wallets (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users (id) ON DELETE RESTRICT,
    balance BIGINT NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS basket_products (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users (id) ON DELETE RESTRICT,
    product_id BIGINT,
    quantity INT,
    price_per_unit BIGINT
);