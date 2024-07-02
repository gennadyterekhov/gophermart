-- +goose Up
-- +goose StatementBegin
create table if not exists users
(
    id serial not null primary key,
    login varchar(255) unique not null,
    password varchar(255) not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
