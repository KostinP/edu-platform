package session

import (
	"context"

	"github.com/kostinp/edu-platform-backend/pkg/db"

	"github.com/google/uuid"
)

type Repository interface {
	Start(ctx context.Context, s *TestSession) error
	Finish(ctx context.Context, sessionID uuid.UUID, score float64) error
	GetByID(ctx context.Context, id uuid.UUID) (*TestSession, error)
	GetLastFinishedByUserAndTest(ctx context.Context, userID, testID uuid.UUID) (*TestSession, error)
}

type PostgresRepo struct{}

func NewPostgresRepo() Repository {
	return &PostgresRepo{}
}

func (r *PostgresRepo) Start(ctx context.Context, s *TestSession) error {
	_, err := db.Pool.Exec(ctx, `
		INSERT INTO test_sessions (id, user_id, test_id, started_at, attempts, status, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), $4, $5, NOW(), NOW())
	`, s.ID, s.UserID, s.TestID, s.Attempts, s.Status)
	return err
}

func (r *PostgresRepo) Finish(ctx context.Context, id uuid.UUID, score float64) error {
	_, err := db.Pool.Exec(ctx, `
		UPDATE test_sessions
		SET score = $1, finished_at = NOW(), status = 'finished', updated_at = NOW()
		WHERE id = $2
	`, score, id)
	return err
}

func (r *PostgresRepo) GetByID(ctx context.Context, id uuid.UUID) (*TestSession, error) {
	row := db.Pool.QueryRow(ctx, `
		SELECT id, user_id, test_id, started_at, finished_at, score, attempts, status, created_at, updated_at
		FROM test_sessions
		WHERE id = $1
	`, id)

	var s TestSession
	err := row.Scan(&s.ID, &s.UserID, &s.TestID, &s.StartedAt, &s.FinishedAt, &s.Score, &s.Attempts, &s.Status, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *PostgresRepo) GetLastFinishedByUserAndTest(ctx context.Context, userID, testID uuid.UUID) (*TestSession, error) {
	row := db.Pool.QueryRow(ctx, `
		SELECT id, user_id, test_id, started_at, finished_at, score, attempts, status, created_at, updated_at
		FROM test_sessions
		WHERE user_id = $1 AND test_id = $2 AND status = 'finished'
		ORDER BY finished_at DESC
		LIMIT 1
	`, userID, testID)

	var s TestSession
	err := row.Scan(&s.ID, &s.UserID, &s.TestID, &s.StartedAt, &s.FinishedAt, &s.Score, &s.Attempts, &s.Status, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &s, nil
}
