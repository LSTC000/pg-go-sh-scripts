-- +goose Up
-- +goose StatementBegin
DROP INDEX IF EXISTS scripts.bash_log_id;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
CREATE UNIQUE INDEX IF NOT EXISTS bash_log_id
ON scripts.bash_log (id);
-- +goose StatementEnd
