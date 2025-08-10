-- +goose Up
CREATE TABLE IF NOT EXISTS wf_tasks
(
    id          BIGSERIAL PRIMARY KEY,
    instance_id BIGINT         NOT NULL REFERENCES wf_instances (id) ON DELETE CASCADE,
    type        wf_task_type   NOT NULL,
    name        TEXT           NOT NULL,
    status      wf_task_status NOT NULL DEFAULT 'ready',
    assignee    TEXT,
    candidates  JSONB,
    due_at      TIMESTAMPTZ,
    payload     JSONB                   DEFAULT '{}'::jsonb,
    retry_count INT            NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ    NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ    NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_wf_tasks_instance ON wf_tasks (instance_id);
CREATE INDEX IF NOT EXISTS idx_wf_tasks_status_due ON wf_tasks (status, due_at);
CREATE INDEX IF NOT EXISTS idx_wf_tasks_assignee ON wf_tasks (assignee);

CREATE TRIGGER trg_wf_tasks_updated
    BEFORE UPDATE
    ON wf_tasks
    FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

-- +goose Down
DROP TRIGGER IF EXISTS trg_wf_tasks_updated ON wf_tasks;
DROP INDEX IF EXISTS idx_wf_tasks_assignee;
DROP INDEX IF EXISTS idx_wf_tasks_status_due;
DROP INDEX IF EXISTS idx_wf_tasks_instance;
DROP TABLE IF EXISTS wf_tasks;
