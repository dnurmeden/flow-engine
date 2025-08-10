package models

import (
	"time"
)

type ProcessDefinition struct {
	ID         int64     `db:"id"`
	Key        string    `db:"key"`
	Version    int       `db:"version"`
	Definition []byte    `db:"definition"`
	IsActive   bool      `db:"is_active"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}

type ProcessInstance struct {
	ID           int64     `db:"id"`
	DefinitionID int64     `db:"definition_id"`
	Status       string    `db:"status"`
	Ctx          []byte    `db:"ctx"`
	Tokens       []byte    `db:"tokens"`
	TenantID     *string   `db:"tenant_id"`
	Rev          int       `db:"rev"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

type Task struct {
	ID         int64      `json:"id"`
	InstanceID int64      `json:"instance_id"`
	Type       string     `json:"type"`   // "user"
	Name       string     `json:"name"`   // "draft", "manager_approve", ...
	Status     string     `json:"status"` // ready|claimed|in_progress|completed|canceled|error
	Assignee   *string    `json:"assignee,omitempty"`
	Candidates []byte     `json:"candidates,omitempty"` // json
	DueAt      *time.Time `json:"due_at,omitempty"`
	Payload    []byte     `json:"payload,omitempty"` // json
	RetryCount int        `json:"retry_count"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

type ClaimTaskRequest struct {
	User string `json:"user"` // login/id исполняющего
}

type CompleteTaskRequest struct {
	User           string                 `json:"user"`
	Output         map[string]interface{} `json:"output"`
	IdempotencyKey string                 `json:"idempotencyKey,omitempty"`
}

type StartProcessRequest struct {
	Key     string                 `json:"key"`
	Version *int                   `json:"version,omitempty"`
	Ctx     map[string]interface{} `json:"ctx"`
}

type StartProcessResponse struct {
	InstanceID int64  `json:"instance_id"`
	Status     string `json:"status"`
}

type GetInstanceResponse struct {
	Instance ProcessInstance `json:"instance"`
	Tasks    []Task          `json:"tasks"`
}
