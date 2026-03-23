create table if not exists users (
    id bigserial primary key,
    login varchar(64) not null,
    password_hash varchar(512) not null
);

create table if not exists wallets (
    id bigserial primary key,
    user_id bigint references users (id) on delete cascade,
    balance bigint not null default 0 check ( balance >= 0 )
);

create table if not exists basket_content (
    id bigserial primary key,
    user_id bigint references users(id) on delete cascade,
    product_id bigint not null,
    product_quantity bigint not null,
    current_price bigint not null,
    CONSTRAINT user_id_product_id UNIQUE (user_id, product_id)
);
