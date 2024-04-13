-- +goose Up
-- +goose StatementBegin
ALTER TABLE IF EXISTS
    scripts.bash_log
DROP CONSTRAINT IF EXISTS
    bash_log_bash_id_fkey;

ALTER TABLE IF EXISTS
    scripts.bash_log
ADD CONSTRAINT
    bash_log_bash_id_fkey
FOREIGN KEY (bash_id) REFERENCES scripts.bash (id) ON DELETE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE IF EXISTS
    scripts.bash_log
DROP CONSTRAINT IF EXISTS
    bash_log_bash_id_fkey;

ALTER TABLE IF EXISTS
    scripts.bash_log
ADD CONSTRAINT
    bash_log_bash_id_fkey
FOREIGN KEY (bash_id) REFERENCES scripts.bash (id);
-- +goose StatementEnd
