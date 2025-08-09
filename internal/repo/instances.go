package repo

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/dnurmeden/flow-engine/internal/models"
)

type InstanceRepo struct {
	db *sql.DB
}

func NewInstanceRepo(db *sql.DB) *InstanceRepo {
	return &InstanceRepo{db: db}
}

func (r *InstanceRepo) Create(ctx context.Context, defID int64, ctxVars map[string]interface{}) (int64, error) {
	ctxJSON, _ := json.Marshal(ctxVars)
	tokens := []byte("[]")
	var id int64

	err := r.db.QueryRowContext(ctx, `
		INSERT INTO wf_instances (definition_id, status, ctx, tokens, rev)
		VALUES ($1, 'running', $2, $3, 0)
		RETURNING id
	`, defID, ctxJSON, tokens).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *InstanceRepo) LogEvent(ctx context.Context, instanceID int64, eventType string, payload map[string]interface{}) error {
	payloadJSON, _ := json.Marshal(payload)
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO wf_events (instance_id, type, payload)
		VALUES ($1, $2, $3)
	`, instanceID, eventType, payloadJSON)
	return err
}

func (r *InstanceRepo) GetByID(ctx context.Context, id int64) (*models.ProcessInstance, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, definition_id, status, ctx, tokens, tenant_id, rev, created_at, updated_at
		FROM wf_instances WHERE id = $1
	`, id)

	var inst models.ProcessInstance
	if err := row.Scan(&inst.ID, &inst.DefinitionID, &inst.Status, &inst.Ctx, &inst.Tokens,
		&inst.TenantID, &inst.Rev, &inst.CreatedAt, &inst.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &inst, nil
}

func (r *InstanceRepo) ListOpenTasks(ctx context.Context, instanceID int64) ([]models.Task, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, instance_id, type, name, status, assignee, candidates, due_at, payload, retry_count, created_at, updated_at
		FROM wf_tasks
		WHERE instance_id = $1 AND status IN ('ready','claimed','in_progress')
		ORDER BY id
	`, instanceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.Task
	for rows.Next() {
		var t models.Task
		if err := rows.Scan(&t.ID, &t.InstanceID, &t.Type, &t.Name, &t.Status, &t.Assignee,
			&t.Candidates, &t.DueAt, &t.Payload, &t.RetryCount, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}
