-- +goose Up
-- +goose StatementBegin
INSERT INTO stocks (sku, total_count, reserved)
SELECT 1076963, 30, 0
WHERE NOT EXISTS (SELECT 1 FROM stocks WHERE sku = 1076963);

INSERT INTO stocks (sku, total_count, reserved)
SELECT 1148162, 30, 0
WHERE NOT EXISTS (SELECT 1 FROM stocks WHERE sku = 1148162);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE stocks;
-- +goose StatementEnd
