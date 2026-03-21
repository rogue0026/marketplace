CREATE TABLE IF NOT EXISTS products (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    price_per_unit BIGINT NOT NULL,
    total_quantity BIGINT NOT NULL CHECK (total_quantity >= 0)
);

drop table products cascade;

CREATE TABLE IF NOT EXISTS product_reservations (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL,
    product_id BIGINT REFERENCES products (id) ON DELETE RESTRICT,
    quantity BIGINT NOT NULL,
    status INT NOT NULL,
    expires_at TIMESTAMP NOT NULL
);

DROP TABLE product_reservations;

select (current_timestamp < current_timestamp + 10 * interval '1 minutes') as result;
