package repo

import (
	"context"
	"errors"

	"database/sql"
	"github.com/dnurmeden/flow-engine/internal/models"
)

type DefinitionRepo struct {
	db *sql.DB
}

func NewDefinitionRepo(db *sql.DB) *DefinitionRepo {
	return &DefinitionRepo{db: db}
}

func (r *DefinitionRepo) GetByKeyAndVersion(ctx context.Context, key string, version *int) (*models.ProcessDefinition, error) {
	def := &models.ProcessDefinition{}
	if version != nil {
		err := r.db.QueryRowContext(ctx, `
			SELECT id, key, version, definition, is_active, created_at, updated_at
			FROM wf_definitions
			WHERE key = $1 AND version = $2
		`, key, *version).Scan(&def.ID, &def.Key, &def.Version, &def.Definition, &def.IsActive, &def.CreatedAt, &def.UpdatedAt)
		if err != nil {
			return nil, err
		}
		return def, nil
	}

	// версия не указана → берём активную
	err := r.db.QueryRowContext(ctx, `
		SELECT id, key, version, definition, is_active, created_at, updated_at
		FROM wf_definitions
		WHERE key = $1 AND is_active = TRUE
		ORDER BY version DESC
		LIMIT 1
	`, key).Scan(&def.ID, &def.Key, &def.Version, &def.Definition, &def.IsActive, &def.CreatedAt, &def.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return def, nil
}
