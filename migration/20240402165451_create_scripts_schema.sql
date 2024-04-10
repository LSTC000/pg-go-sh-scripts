-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA IF NOT EXISTS scripts;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP SCHEMA IF EXISTS scripts;
-- +goose StatementEnd
