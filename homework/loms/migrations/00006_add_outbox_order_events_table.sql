-- +goose Up
-- +goose StatementBegin

create table if not exists outbox_order_events
(
    id          serial primary key,
    order_id    bigint not null,
    event_type  varchar not null,
    was_sent    bool not null default false,
    created_at  timestamp with time zone default now()
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS outbox_order_events CASCADE;
-- +goose StatementEnd
