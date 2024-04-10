-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS scripts.bash (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR NOT NULL,
    body VARCHAR NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX IF NOT EXISTS bash_id
ON scripts.bash (id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS bash_id;

DROP TABLE IF EXISTS scripts.bash;
-- +goose StatementEnd
