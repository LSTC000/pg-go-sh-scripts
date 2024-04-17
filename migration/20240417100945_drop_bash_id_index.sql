-- +goose Up
-- +goose StatementBegin
DROP INDEX IF EXISTS scripts.bash_id;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
CREATE UNIQUE INDEX IF NOT EXISTS bash_id
ON scripts.bash (id);
-- +goose StatementEnd
