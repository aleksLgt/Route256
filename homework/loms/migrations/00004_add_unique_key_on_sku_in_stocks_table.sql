-- +goose Up
-- +goose StatementBegin
ALTER TABLE stocks
ADD CONSTRAINT sku_unique UNIQUE (sku);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE stocks DROP CONSTRAINT sku_unique;
-- +goose StatementEnd
