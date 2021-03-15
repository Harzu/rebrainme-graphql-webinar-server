-- +goose Up
-- +goose StatementBegin
CREATE TYPE roles AS ENUM ('ADMIN', 'CUSTOMER');

CREATE TABLE IF NOT EXISTS users (
    id            bigserial primary key,
    email         text not null unique,
    password_hash text not null,
    role          roles not null,
    created_at    timestamp with time zone not null,
    updated_at    timestamp with time zone,
    deleted_at    timestamp with time zone
);

CREATE TABLE IF NOT EXISTS customers (
    id         bigserial primary key,
    user_id    integer not null references users(id),
    name       text not null,
    address    text not null,
    created_at timestamp with time zone not null,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);

CREATE TABLE IF NOT EXISTS products (
    id         bigserial primary key,
    name       text not null unique,
    price      integer default 0,
    created_at timestamp with time zone not null,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);

CREATE TYPE order_status AS ENUM ('CREATED', 'DONE', 'CANCELED');

CREATE TABLE IF NOT EXISTS orders (
    id          bigserial primary key,
    status      order_status not null default 'CREATED',
    total_price integer default 0,
    customer_id integer not null references customers(id),
    created_at  timestamp with time zone not null,
    updated_at  timestamp with time zone,
    deleted_at  timestamp with time zone
);

CREATE TABLE IF NOT EXISTS order_products (
    id         bigserial primary key,
    order_id   integer not null references orders(id),
    product_id integer not null references products(id),
    created_at timestamp with time zone not null,
    deleted_at timestamp with time zone
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS order_products;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS customers;
DROP TABLE IF EXISTS users;

DROP TYPE IF EXISTS order_status;
DROP TYPE IF EXISTS roles;
-- +goose StatementEnd
