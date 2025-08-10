package repo

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dnurmeden/flow-engine/internal/models"
)

type TaskRepo struct{ db *sql.DB }

func NewTaskRepo(db *sql.DB) *TaskRepo { return &TaskRepo{db: db} }

// Создать user_task
func (r *TaskRepo) CreateUserTask(ctx context.Context, instanceID int64, name string, candidates map[string]any, dueAt *string) (int64, error) {
	var id int64
	candJSON, _ := json.Marshal(candidates)
	q := `INSERT INTO wf_tasks (instance_id,type,name,status,candidates,due_at)
	      VALUES ($1,'user',$2,'ready',$3, NULL) RETURNING id`
	if err := r.db.QueryRowContext(ctx, q, instanceID, name, candJSON).Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

// Claim: только из ready → claimed, атомарно
func (r *TaskRepo) Claim(ctx context.Context, taskID int64, user string) error {
	res, err := r.db.ExecContext(ctx, `
		UPDATE wf_tasks
		   SET status='claimed', assignee=$2, updated_at=now()
		 WHERE id=$1 AND status='ready'
	`, taskID, user)
	if err != nil {
		return err
	}
	aff, _ := res.RowsAffected()
	if aff == 0 {
		return errors.New("task not in claimable state")
	}
	return nil
}

// Complete: только исполнитель и только из claimed|in_progress → completed
func (r *TaskRepo) Complete(ctx context.Context, taskID int64, user string, output map[string]any) error {
	outJSON, _ := json.Marshal(map[string]any{"output": output})
	res, err := r.db.ExecContext(ctx, `
		UPDATE wf_tasks
		   SET status='completed', payload=$3, updated_at=now()
		 WHERE id=$1 AND assignee=$2 AND status IN ('claimed','in_progress')
	`, taskID, user, outJSON)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return errors.New("task not in completable state or assignee mismatch")
	}
	return nil
}

func (r *TaskRepo) GetByID(ctx context.Context, id int64) (*models.Task, error) {
	var t models.Task
	row := r.db.QueryRowContext(ctx, `
		SELECT id, instance_id, type, name, status, assignee, candidates, due_at, payload, retry_count, created_at, updated_at
		FROM wf_tasks WHERE id=$1
	`, id)
	if err := row.Scan(&t.ID, &t.InstanceID, &t.Type, &t.Name, &t.Status, &t.Assignee, &t.Candidates,
		&t.DueAt, &t.Payload, &t.RetryCount, &t.CreatedAt, &t.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &t, nil
}

func (r *TaskRepo) ListInbox(ctx context.Context, user string) ([]models.Task, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, instance_id, type, name, status, assignee, candidates, due_at, payload, retry_count, created_at, updated_at
		  FROM wf_tasks
		 WHERE assignee=$1 AND status IN ('claimed','in_progress','ready') -- на будущее можно разделить
		 ORDER BY id DESC
	`, user)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.Task
	for rows.Next() {
		var t models.Task
		if err := rows.Scan(&t.ID, &t.InstanceID, &t.Type, &t.Name, &t.Status, &t.Assignee, &t.Candidates,
			&t.DueAt, &t.Payload, &t.RetryCount, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

// Утилита: запись события
func (r *TaskRepo) LogEvent(ctx context.Context, instanceID int64, eventType string, payload map[string]any) error {
	b, _ := json.Marshal(payload)
	_, err := r.db.ExecContext(ctx, `INSERT INTO wf_events(instance_id,type,payload) VALUES ($1,$2,$3)`, instanceID, eventType, b)
	return err
}

// На будущее: продвижение процесса дальше после complete (тут пока заглушка)
func (r *TaskRepo) AdvanceAfterComplete(ctx context.Context, instanceID int64) error {
	// MVP: ничего не делаем. Позже тут будет оркестратор переходов.
	fmt.Println("advance placeholder for instance", instanceID)
	return nil
}
