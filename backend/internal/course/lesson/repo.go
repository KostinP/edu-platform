package lesson

import (
	"context"

	"github.com/kostinp/edu-platform-backend/pkg/db"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, l *Lesson) error
	GetByID(ctx context.Context, id uuid.UUID) (*Lesson, error)
	GetByModuleID(ctx context.Context, moduleID uuid.UUID) ([]Lesson, error)
	Update(ctx context.Context, l *Lesson) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
}

type PostgresRepo struct{}

func NewPostgresRepo() Repository {
	return &PostgresRepo{}
}

func (r *PostgresRepo) Create(ctx context.Context, l *Lesson) error {
	_, err := db.DB().Exec(ctx, `
		INSERT INTO lessons (id, module_id, title, content, ordinal, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
	`, l.ID, l.ModuleID, l.Title, l.Content, l.Ordinal)
	return err
}

func (r *PostgresRepo) GetByModuleID(ctx context.Context, moduleID uuid.UUID) ([]Lesson, error) {
	rows, err := db.DB().Query(ctx, `
		SELECT id, module_id, title, content, ordinal, created_at
		FROM lessons
		WHERE module_id = $1
		ORDER BY ordinal ASC
	`, moduleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lessons []Lesson
	for rows.Next() {
		var l Lesson
		if err := rows.Scan(&l.ID, &l.ModuleID, &l.Title, &l.Content, &l.Ordinal, &l.CreatedAt); err != nil {
			return nil, err
		}
		lessons = append(lessons, l)
	}
	return lessons, nil
}

func (r *PostgresRepo) GetByID(ctx context.Context, id uuid.UUID) (*Lesson, error) {
	row := db.DB().QueryRow(ctx, `
		SELECT id, module_id, title, content, ordinal, author_id, created_at, updated_at, deleted_at
		FROM lessons
		WHERE id = $1
	`, id)

	var l Lesson
	err := row.Scan(&l.ID, &l.ModuleID, &l.Title, &l.Content, &l.Ordinal, &l.AuthorID, &l.CreatedAt, &l.UpdatedAt, &l.DeletedAt)
	if err != nil {
		return nil, err
	}
	return &l, nil
}

func (r *PostgresRepo) Update(ctx context.Context, l *Lesson) error {
	_, err := db.DB().Exec(ctx, `
		UPDATE lessons SET
			title = $1,
			content = $2,
			ordinal = $3,
			updated_at = NOW()
		WHERE id = $4 AND deleted_at IS NULL
	`, l.Title, l.Content, l.Ordinal, l.ID)
	return err
}

func (r *PostgresRepo) SoftDelete(ctx context.Context, id uuid.UUID) error {
	_, err := db.DB().Exec(ctx, `
		UPDATE lessons SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`, id)
	return err
}
