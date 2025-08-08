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

type StartProcessRequest struct {
	Key     string                 `json:"key"`
	Version *int                   `json:"version,omitempty"`
	Ctx     map[string]interface{} `json:"ctx"`
}

type StartProcessResponse struct {
	InstanceID int64  `json:"instance_id"`
	Status     string `json:"status"`
}
