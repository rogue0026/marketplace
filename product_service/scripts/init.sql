create table if not exists products (
    id bigserial primary key,
    name varchar(128) not null,
    stock_remaining bigint not null check ( stock_remaining >= 0 ),
    current_price bigint not null
);

create table if not exists product_reservations (
    id bigserial primary key,
    order_id bigint not null,
    product_id bigint references products (id) on delete restrict,
    quantity bigint not null check ( quantity >= 0 )
);


