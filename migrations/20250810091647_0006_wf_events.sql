-- +goose Up
CREATE TABLE IF NOT EXISTS wf_events
(
    id          BIGSERIAL PRIMARY KEY,
    instance_id BIGINT      NOT NULL REFERENCES wf_instances (id) ON DELETE CASCADE,
    type        TEXT        NOT NULL,
    payload     JSONB       NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_wf_events_instance ON wf_events (instance_id);
CREATE INDEX IF NOT EXISTS idx_wf_events_type_time ON wf_events (type, created_at);

-- +goose Down
DROP INDEX IF EXISTS idx_wf_events_type_time;
DROP INDEX IF EXISTS idx_wf_events_instance;
DROP TABLE IF EXISTS wf_events;
