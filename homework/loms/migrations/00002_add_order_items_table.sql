-- +goose Up
-- +goose StatementBegin

create table if not exists order_items
(
    id        serial primary key,
    order_id  bigint not null,
    sku       int not null,
    count     int not null
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS order_items CASCADE;
-- +goose StatementEnd
