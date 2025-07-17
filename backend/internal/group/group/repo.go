package group

import (
	"context"
	"time"

	"github.com/kostinp/edu-platform-backend/pkg/db"

	"github.com/google/uuid"
)

type PostgresRepo struct{}

type Repository interface {
	Create(ctx context.Context, g *Group) error
	GetByID(ctx context.Context, id uuid.UUID) (*Group, error)
	ListByUser(ctx context.Context, userID uuid.UUID) ([]Group, error)
	AddMember(ctx context.Context, groupID, userID uuid.UUID, role string) error
}

func NewRepository() Repository {
	return &PostgresRepo{}
}

func (r *PostgresRepo) Create(ctx context.Context, g *Group) error {
	_, err := db.Pool.Exec(ctx, `
		INSERT INTO groups (id, name, description, owner_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, g.ID, g.Name, g.Description, g.OwnerID, time.Now(), time.Now())
	return err
}

func (r *PostgresRepo) GetByID(ctx context.Context, id uuid.UUID) (*Group, error) {
	row := db.Pool.QueryRow(ctx, `
		SELECT id, name, description, owner_id, created_at, updated_at
		FROM groups WHERE id = $1
	`, id)

	var g Group
	err := row.Scan(&g.ID, &g.Name, &g.Description, &g.OwnerID, &g.CreatedAt, &g.UpdatedAt)
	return &g, err
}

func (r *PostgresRepo) ListByUser(ctx context.Context, userID uuid.UUID) ([]Group, error) {
	rows, err := db.Pool.Query(ctx, `
		SELECT g.id, g.name, g.description, g.owner_id, g.created_at, g.updated_at
		FROM groups g
		JOIN group_members gm ON gm.group_id = g.id
		WHERE gm.user_id = $1
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []Group
	for rows.Next() {
		var g Group
		err := rows.Scan(&g.ID, &g.Name, &g.Description, &g.OwnerID, &g.CreatedAt, &g.UpdatedAt)
		if err != nil {
			return nil, err
		}
		result = append(result, g)
	}
	return result, nil
}

func (r *PostgresRepo) AddMember(ctx context.Context, groupID, userID uuid.UUID, role string) error {
	_, err := db.Pool.Exec(ctx, `
		INSERT INTO group_members (group_id, user_id, role)
		VALUES ($1, $2, $3)
		ON CONFLICT DO NOTHING
	`, groupID, userID, role)
	return err
}
