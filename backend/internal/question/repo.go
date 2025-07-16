package question

import (
	"context"

	"github.com/kostinp/edu-platform-backend/pkg/db"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, q *Question) error
	GetByID(ctx context.Context, id uuid.UUID) (*Question, error)
	GetByTestID(ctx context.Context, testID uuid.UUID) ([]Question, error)
	Update(ctx context.Context, q *Question) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type PostgresRepo struct{}

func NewPostgresRepo() Repository {
	return &PostgresRepo{}
}

func (r *PostgresRepo) Create(ctx context.Context, q *Question) error {
	_, err := db.Pool.Exec(ctx, `
		INSERT INTO questions (id, test_id, type, title, image_url, data, feedback, score, ordinal, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
	`, q.ID, q.TestID, q.Type, q.Title, q.ImageURL, q.Data, q.Feedback, q.Score, q.Ordinal, q.CreatedAt, q.UpdatedAt)
	return err
}

func (r *PostgresRepo) GetByID(ctx context.Context, id uuid.UUID) (*Question, error) {
	row := db.Pool.QueryRow(ctx, `
		SELECT id, test_id, type, title, image_url, data, feedback, score, ordinal, created_at, updated_at, deleted_at
		FROM questions WHERE id = $1 AND deleted_at IS NULL
	`, id)

	var q Question
	err := row.Scan(&q.ID, &q.TestID, &q.Type, &q.Title, &q.ImageURL, &q.Data, &q.Feedback, &q.Score, &q.Ordinal, &q.CreatedAt, &q.UpdatedAt, &q.DeletedAt)
	if err != nil {
		return nil, err
	}
	return &q, nil
}

func (r *PostgresRepo) GetByTestID(ctx context.Context, testID uuid.UUID) ([]Question, error) {
	rows, err := db.Pool.Query(ctx, `
		SELECT id, test_id, author_id, type, title, image_url, data, feedback, score, ordinal, created_at, updated_at
		FROM questions WHERE test_id = $1 AND deleted_at IS NULL
		ORDER BY ordinal
	`, testID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []Question
	for rows.Next() {
		var q Question
		err := rows.Scan(&q.ID, &q.TestID, &q.AuthorID, &q.Type, &q.Title, &q.ImageURL, &q.Data, &q.Feedback, &q.Score, &q.Ordinal, &q.CreatedAt, &q.UpdatedAt)
		if err != nil {
			return nil, err
		}
		result = append(result, q)
	}
	return result, nil
}

func (r *PostgresRepo) Update(ctx context.Context, q *Question) error {
	_, err := db.Pool.Exec(ctx, `
		UPDATE questions
		SET type = $1, title = $2, image_url = $3, data = $4, feedback = $5,
			score = $6, ordinal = $7, updated_at = $8
		WHERE id = $9 AND deleted_at IS NULL
	`, q.Type, q.Title, q.ImageURL, q.Data, q.Feedback, q.Score, q.Ordinal, q.UpdatedAt, q.ID)
	return err
}

func (r *PostgresRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := db.Pool.Exec(ctx, `
		UPDATE questions SET deleted_at = NOW() WHERE id = $1
	`, id)
	return err
}
