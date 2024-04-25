-- +goose Up
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS bash_log_bash_id_fkey
ON scripts.bash_log (bash_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS scripts.bash_log_bash_id_fkey;
-- +goose StatementEnd
