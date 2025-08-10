-- +goose Up
CREATE TABLE IF NOT EXISTS wf_instances
(
    id            BIGSERIAL PRIMARY KEY,
    definition_id BIGINT             NOT NULL REFERENCES wf_definitions (id) ON DELETE RESTRICT,
    status        wf_instance_status NOT NULL DEFAULT 'running',
    ctx           JSONB              NOT NULL DEFAULT '{}'::jsonb,
    tokens        JSONB              NOT NULL DEFAULT '[]'::jsonb,
    tenant_id     TEXT,
    rev           INT                NOT NULL DEFAULT 0,
    created_at    TIMESTAMPTZ        NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ        NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_wf_instances_def ON wf_instances (definition_id);
CREATE INDEX IF NOT EXISTS idx_wf_instances_status ON wf_instances (status);
CREATE INDEX IF NOT EXISTS idx_wf_instances_tenant ON wf_instances (tenant_id);
CREATE INDEX IF NOT EXISTS idx_wf_instances_ctx_doc_id ON wf_instances ((ctx ->> 'doc_id'));

CREATE TRIGGER trg_wf_instances_updated
    BEFORE UPDATE
    ON wf_instances
    FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

-- +goose Down
DROP TRIGGER IF EXISTS trg_wf_instances_updated ON wf_instances;
DROP INDEX IF EXISTS idx_wf_instances_ctx_doc_id;
DROP INDEX IF EXISTS idx_wf_instances_tenant;
DROP INDEX IF EXISTS idx_wf_instances_status;
DROP INDEX IF EXISTS idx_wf_instances_def;
DROP TABLE IF EXISTS wf_instances;
