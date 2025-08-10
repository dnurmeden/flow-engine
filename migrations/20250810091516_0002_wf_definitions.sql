-- +goose Up
CREATE TABLE IF NOT EXISTS wf_definitions
(
    id         BIGSERIAL PRIMARY KEY,
    key        TEXT        NOT NULL,
    version    INT         NOT NULL,
    definition JSONB       NOT NULL,
    is_active  BOOLEAN     NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (key, version)
);
CREATE INDEX IF NOT EXISTS idx_wf_definitions_key ON wf_definitions (key);

-- +goose Down
DROP INDEX IF EXISTS idx_wf_definitions_key;
DROP TABLE IF EXISTS wf_definitions;
