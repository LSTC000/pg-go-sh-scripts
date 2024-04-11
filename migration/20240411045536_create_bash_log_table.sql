-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS scripts.bash_log (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    bash_id uuid NOT NULL,
    body VARCHAR NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    FOREIGN KEY (bash_id) REFERENCES scripts.bash (id)
);

CREATE UNIQUE INDEX IF NOT EXISTS bash_log_id
ON scripts.bash_log (id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS bash_log_id;

DROP TABLE IF EXISTS scripts.bash_log;
-- +goose StatementEnd
