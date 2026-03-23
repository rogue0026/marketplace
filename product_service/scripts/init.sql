create table if not exists products (
    id bigserial primary key,
    name varchar(128) not null,
    stock_remaining bigint not null check ( stock_remaining >= 0 ),
    current_price bigint not null
);

create table if not exists product_reservations (
    id bigserial primary key,
    product_id bigint references products (id) on delete restrict,
    amount bigint not null check ( amount >= 0 ),
    order_id bigint not null,
    expires_at timestamp not null default (now() + 15 * interval '1 minute')
);


