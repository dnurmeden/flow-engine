-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- +goose StatementBegin
DO $do$
    BEGIN
        IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname='wf_instance_status') THEN
            CREATE TYPE wf_instance_status AS ENUM ('running','paused','canceled','completed');
        END IF;

        IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname='wf_task_type') THEN
            CREATE TYPE wf_task_type AS ENUM ('user','service');
        END IF;

        IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname='wf_task_status') THEN
            CREATE TYPE wf_task_status AS ENUM ('ready','claimed','in_progress','completed','canceled','error');
        END IF;
    END;
$do$;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION set_updated_at()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP FUNCTION IF EXISTS set_updated_at;
-- +goose StatementEnd

-- +goose StatementBegin
DO $do$
    BEGIN
        IF EXISTS (SELECT 1 FROM pg_type WHERE typname='wf_task_status') THEN
            DROP TYPE wf_task_status;
        END IF;
        IF EXISTS (SELECT 1 FROM pg_type WHERE typname='wf_task_type') THEN
            DROP TYPE wf_task_type;
        END IF;
        IF EXISTS (SELECT 1 FROM pg_type WHERE typname='wf_instance_status') THEN
            DROP TYPE wf_instance_status;
        END IF;
    END;
$do$;
-- +goose StatementEnd
