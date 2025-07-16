package module

import (
	"context"

	"github.com/kostinp/edu-platform-backend/pkg/db"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, m *Module) error
	GetByCourseID(ctx context.Context, courseID uuid.UUID) ([]Module, error)
	Update(ctx context.Context, m *Module) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
}

type PostgresRepo struct{}

func NewPostgresRepo() Repository {
	return &PostgresRepo{}
}

func (r *PostgresRepo) Create(ctx context.Context, m *Module) error {
	_, err := db.Pool.Exec(ctx, `
		INSERT INTO modules (id, course_id, title, description, ordinal, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
	`, m.ID, m.CourseID, m.Title, m.Description, m.Ordinal)
	return err
}

func (r *PostgresRepo) GetByCourseID(ctx context.Context, courseID uuid.UUID) ([]Module, error) {
	rows, err := db.Pool.Query(ctx, `
		SELECT id, course_id, title, description, ordinal, created_at
		FROM modules
		WHERE course_id = $1
		ORDER BY ordinal ASC
	`, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var modules []Module
	for rows.Next() {
		var m Module
		if err := rows.Scan(&m.ID, &m.CourseID, &m.Title, &m.Description, &m.Ordinal, &m.CreatedAt); err != nil {
			return nil, err
		}
		modules = append(modules, m)
	}
	return modules, nil
}

func (r *PostgresRepo) Update(ctx context.Context, m *Module) error {
	_, err := db.Pool.Exec(ctx, `
		UPDATE modules SET
			title = $1,
			description = $2,
			ordinal = $3,
			updated_at = NOW()
		WHERE id = $4 AND deleted_at IS NULL
	`, m.Title, m.Description, m.Ordinal, m.ID)
	return err
}

func (r *PostgresRepo) SoftDelete(ctx context.Context, id uuid.UUID) error {
	_, err := db.Pool.Exec(ctx, `
		UPDATE modules SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`, id)
	return err
}
