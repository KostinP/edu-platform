package test

import (
	"context"
	"time"

	"github.com/kostinp/edu-platform-backend/pkg/db"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, t *Test) error
	Update(ctx context.Context, t *Test) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*Test, error)
	GetByLessonID(ctx context.Context, lessonID uuid.UUID) ([]Test, error)
}

type PostgresRepo struct{}

func NewPostgresRepo() Repository {
	return &PostgresRepo{}
}

func (r *PostgresRepo) Create(ctx context.Context, t *Test) error {
	_, err := db.Pool.Exec(ctx, `
		INSERT INTO tests (id, title, time_limit, shuffle, attempts, show_score, show_answer, access_from, access_to, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
	`, t.ID, t.Title, t.TimeLimit, t.Shuffle, t.Attempts, t.ShowScore, t.ShowAnswer, t.AccessFrom, t.AccessTo, t.CreatedAt, t.UpdatedAt)
	return err
}

func (r *PostgresRepo) Update(ctx context.Context, t *Test) error {
	_, err := db.Pool.Exec(ctx, `
		UPDATE tests
		SET title = $1, time_limit = $2, shuffle = $3, attempts = $4,
			show_score = $5, show_answer = $6, access_from = $7, access_to = $8, updated_at = $9
		WHERE id = $10 AND deleted_at IS NULL
	`, t.Title, t.TimeLimit, t.Shuffle, t.Attempts, t.ShowScore, t.ShowAnswer, t.AccessFrom, t.AccessTo, time.Now(), t.ID)
	return err
}

func (r *PostgresRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := db.Pool.Exec(ctx, `
		UPDATE tests SET deleted_at = NOW() WHERE id = $1
	`, id)
	return err
}

func (r *PostgresRepo) GetByID(ctx context.Context, id uuid.UUID) (*Test, error) {
	row := db.Pool.QueryRow(ctx, `
		SELECT id, title, time_limit, shuffle, attempts, show_score, show_answer, access_from, access_to, created_at, updated_at
		FROM tests
		WHERE id = $1 AND deleted_at IS NULL
	`, id)

	var t Test
	err := row.Scan(&t.ID, &t.Title, &t.TimeLimit, &t.Shuffle, &t.Attempts, &t.ShowScore, &t.ShowAnswer, &t.AccessFrom, &t.AccessTo, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *PostgresRepo) GetByLessonID(ctx context.Context, lessonID uuid.UUID) ([]Test, error) {
	rows, err := db.Pool.Query(ctx, `
		SELECT id, lesson_id, title, time_limit, shuffle, attempts, show_score, show_answer, access_from, access_to, author_id, created_at, updated_at, deleted_at
		FROM tests
		WHERE lesson_id = $1 AND deleted_at IS NULL
	`, lessonID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tests []Test
	for rows.Next() {
		var t Test
		err := rows.Scan(&t.ID, &t.LessonID, &t.Title, &t.TimeLimit, &t.Shuffle, &t.Attempts, &t.ShowScore, &t.ShowAnswer, &t.AccessFrom, &t.AccessTo, &t.AuthorID, &t.CreatedAt, &t.UpdatedAt, &t.DeletedAt)
		if err != nil {
			return nil, err
		}
		tests = append(tests, t)
	}
	return tests, nil
}
