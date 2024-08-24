-- +goose Up
-- +goose StatementBegin

create table if not exists orders
(
    id      serial primary key,
    user_id bigint not null,
    status  varchar not null
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS orders CASCADE;
-- +goose StatementEnd
