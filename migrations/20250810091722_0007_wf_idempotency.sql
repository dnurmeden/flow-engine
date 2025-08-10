-- +goose Up
CREATE TABLE IF NOT EXISTS wf_idempotency
(
    key        TEXT PRIMARY KEY,
    scope      TEXT        NOT NULL,
    result     JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE IF EXISTS wf_idempotency;
