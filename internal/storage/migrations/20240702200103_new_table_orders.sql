-- +goose Up
-- +goose StatementBegin
create table if not exists orders
(
    number varchar(64) unique not null primary key,
    user_id int references users(id),
    status varchar(255) not null,
    accrual int default null,
    uploaded_at timestamp with time zone not null default now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS orders;
-- +goose StatementEnd
