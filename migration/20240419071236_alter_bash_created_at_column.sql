-- +goose Up
-- +goose StatementBegin
ALTER TABLE IF EXISTS
    scripts.bash
ALTER COLUMN
    created_at
TYPE TIMESTAMP WITHOUT TIME ZONE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE IF EXISTS
    scripts.bash
ALTER COLUMN
    created_at
TYPE TIMESTAMP WITH TIME ZONE;
-- +goose StatementEnd
