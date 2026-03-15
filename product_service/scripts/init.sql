CREATE TABLE IF NOT EXISTS products (
                                        id BIGSERIAL PRIMARY KEY,
                                        name TEXT NOT NULL,
                                        price_per_unit BIGINT NOT NULL,
                                        total_quantity BIGINT NOT NULL
);

