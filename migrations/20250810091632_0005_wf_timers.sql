-- +goose Up
CREATE TABLE IF NOT EXISTS wf_timers
(
    id          BIGSERIAL PRIMARY KEY,
    instance_id BIGINT      NOT NULL REFERENCES wf_instances (id) ON DELETE CASCADE,
    fire_at     TIMESTAMPTZ NOT NULL,
    action      JSONB       NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_wf_timers_fire_at ON wf_timers (fire_at);

-- +goose Down
DROP INDEX IF EXISTS idx_wf_timers_fire_at;
DROP TABLE IF EXISTS wf_timers;
