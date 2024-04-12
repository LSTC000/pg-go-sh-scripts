-- +goose Up
-- +goose StatementBegin
ALTER TABLE IF EXISTS
    scripts.bash_log
ADD COLUMN IF NOT EXISTS
    is_error BOOLEAN NOT NULL DEFAULT false;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE IF EXISTS
    scripts.bash_log
DROP COLUMN IF EXISTS
    is_error;
-- +goose StatementEnd
