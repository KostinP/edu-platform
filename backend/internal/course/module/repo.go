package module

import (
	"context"

	"github.com/google/uuid"
	"github.com/kostinp/edu-platform-backend/pkg/db"
)

type Repository interface {
	Create(ctx context.Context, m *Module) error
	GetByID(ctx context.Context, id uuid.UUID) (*Module, error)
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
		INSERT INTO modules (id, course_id, title, description, ordinal, author_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
	`, m.ID, m.CourseID, m.Title, m.Description, m.Ordinal, m.AuthorID)
	return err
}

func (r *PostgresRepo) GetByID(ctx context.Context, id uuid.UUID) (*Module, error) {
	row := db.Pool.QueryRow(ctx, `
		SELECT id, course_id, title, description, ordinal, author_id, created_at, updated_at, deleted_at
		FROM modules
		WHERE id = $1
	`, id)

	var m Module
	err := row.Scan(&m.ID, &m.CourseID, &m.Title, &m.Description, &m.Ordinal, &m.AuthorID, &m.CreatedAt, &m.UpdatedAt, &m.DeletedAt)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *PostgresRepo) GetByCourseID(ctx context.Context, courseID uuid.UUID) ([]Module, error) {
	rows, err := db.Pool.Query(ctx, `
		SELECT
			m.id,
			m.course_id,
			m.title,
			m.description,
			m.ordinal,
			m.author_id,
			m.created_at,
			m.updated_at,
			COALESCE(SUM(l.duration), 0) AS duration
		FROM modules m
		LEFT JOIN lessons l ON l.module_id = m.id AND l.deleted_at IS NULL
		WHERE m.course_id = $1 AND m.deleted_at IS NULL
		GROUP BY m.id
		ORDER BY m.ordinal ASC
	`, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var modules []Module
	for rows.Next() {
		var m Module
		if err := rows.Scan(
			&m.ID,
			&m.CourseID,
			&m.Title,
			&m.Description,
			&m.Ordinal,
			&m.AuthorID,
			&m.CreatedAt,
			&m.UpdatedAt,
			&m.Duration,
		); err != nil {
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
