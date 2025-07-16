package answer

import (
	"context"

	"github.com/google/uuid"
	"github.com/kostinp/edu-platform-backend/pkg/db"
)

type Repository interface {
	Submit(ctx context.Context, a *UserAnswer) error
	GetBySession(ctx context.Context, sessionID uuid.UUID) ([]UserAnswer, error)
	GetAnswerStats(ctx context.Context, questionID uuid.UUID) (map[string]int, error)
	ManualCheck(ctx context.Context, id uuid.UUID, markCorrect bool) error
}

type PostgresRepo struct{}

func NewPostgresRepo() Repository {
	return &PostgresRepo{}
}

func (r *PostgresRepo) Submit(ctx context.Context, a *UserAnswer) error {
	_, err := db.Pool.Exec(ctx, `
		INSERT INTO user_answers (id, test_session_id, user_id, question_id, answer, duration_seconds, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW())
	`, a.ID, a.TestSessionID, a.UserID, a.QuestionID, a.Answer, a.DurationSeconds)
	return err
}

func (r *PostgresRepo) GetBySession(ctx context.Context, sessionID uuid.UUID) ([]UserAnswer, error) {
	rows, err := db.Pool.Query(ctx, `
		SELECT id, test_session_id, user_id, question_id, answer, duration_seconds, is_checked, created_at
		FROM user_answers
		WHERE test_session_id = $1
	`, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []UserAnswer
	for rows.Next() {
		var a UserAnswer
		err := rows.Scan(&a.ID, &a.TestSessionID, &a.UserID, &a.QuestionID, &a.Answer, &a.DurationSeconds, &a.IsChecked, &a.CreatedAt)
		if err != nil {
			return nil, err
		}
		result = append(result, a)
	}
	return result, nil
}

func (r *PostgresRepo) GetAnswerStats(ctx context.Context, questionID uuid.UUID) (map[string]int, error) {
	rows, err := db.Pool.Query(ctx, `
		SELECT answer::TEXT, COUNT(*) FROM user_answers
		WHERE question_id = $1
		GROUP BY answer::TEXT
	`, questionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := map[string]int{}
	for rows.Next() {
		var answerText string
		var count int
		if err := rows.Scan(&answerText, &count); err != nil {
			return nil, err
		}
		stats[answerText] = count
	}
	return stats, nil
}

func (r *PostgresRepo) ManualCheck(ctx context.Context, id uuid.UUID, markCorrect bool) error {
	_, err := db.Pool.Exec(ctx, `
		UPDATE user_answers
		SET is_checked = true
		WHERE id = $1
	`, id)
	return err
}
