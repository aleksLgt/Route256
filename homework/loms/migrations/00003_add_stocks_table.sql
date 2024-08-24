-- +goose Up
-- +goose StatementBegin

create table if not exists stocks
(
    id          serial primary key,
    sku         int not null,
    total_count int not null,
    reserved    int not null
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS stocks CASCADE;
-- +goose StatementEnd
