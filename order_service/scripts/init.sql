create table if not exists orders (
    id bigserial primary key,
    user_id bigint not null,
    total_price bigint not null,
    status varchar(64) not null
);

create table if not exists order_contents (
    id bigserial primary key,
    order_id bigint references orders (id) on delete restrict,
    product_id bigint not null,
    product_quantity bigint not null,
    price_per_unit bigint not null
);