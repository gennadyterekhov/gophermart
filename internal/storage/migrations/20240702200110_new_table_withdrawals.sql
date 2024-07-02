-- +goose Up
-- +goose StatementBegin
create table if not exists withdrawals
(
    id serial not null primary key,
    user_id int references users(id),
    order_number varchar(255) not null,
    total_sum int not null,
    processed_at timestamp with time zone not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS withdrawals;
-- +goose StatementEnd
