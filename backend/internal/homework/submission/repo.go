package submission

import (
	"context"
	"time"

	"github.com/kostinp/edu-platform-backend/pkg/db"

	"github.com/google/uuid"
)

type PostgresRepo struct{}

type Repository interface {
	Submit(ctx context.Context, s *Submission) error
	ListByUser(ctx context.Context, userID uuid.UUID) ([]Submission, error)
	Review(ctx context.Context, id uuid.UUID, reviewText string, score float64) error
	GetByID(ctx context.Context, id uuid.UUID) (*Submission, error)
	AddPeerReview(ctx context.Context, id uuid.UUID, peerID uuid.UUID, comment string) error
	GetByUserAndHomework(ctx context.Context, userID, homeworkID uuid.UUID) (*Submission, error)
}

func NewPostgresRepo() Repository {
	return &PostgresRepo{}
}

func (r *PostgresRepo) Submit(ctx context.Context, s *Submission) error {
	_, err := db.DB().Exec(ctx, `
		INSERT INTO homework_submissions
		(id, homework_id, user_id, status, answer, file_url, created_at, updated_at)
		VALUES ($1,$2,$3,'submitted',$4,$5,$6,$7)
	`, s.ID, s.HomeworkID, s.UserID, s.Answer, s.FileURL, time.Now(), time.Now())
	return err
}

func (r *PostgresRepo) ListByUser(ctx context.Context, userID uuid.UUID) ([]Submission, error) {
	rows, err := db.DB().Query(ctx, `
		SELECT id, homework_id, user_id, status, answer, file_url, review, score, created_at, updated_at
		FROM homework_submissions WHERE user_id = $1
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []Submission
	for rows.Next() {
		var s Submission
		err := rows.Scan(&s.ID, &s.HomeworkID, &s.UserID, &s.Status, &s.Answer, &s.FileURL, &s.Review, &s.Score, &s.CreatedAt, &s.UpdatedAt)
		if err != nil {
			return nil, err
		}
		result = append(result, s)
	}
	return result, nil
}

func (r *PostgresRepo) Review(ctx context.Context, id uuid.UUID, review string, score float64) error {
	_, err := db.DB().Exec(ctx, `
		UPDATE homework_submissions
		SET review = $1, score = $2, status = 'reviewed', updated_at = NOW()
		WHERE id = $3
	`, review, score, id)
	return err
}

func (r *PostgresRepo) GetByID(ctx context.Context, id uuid.UUID) (*Submission, error) {
	row := db.DB().QueryRow(ctx, `
		SELECT id, homework_id, user_id, status, answer, file_url, review, score, created_at, updated_at
		FROM homework_submissions WHERE id = $1
	`, id)

	var s Submission
	err := row.Scan(&s.ID, &s.HomeworkID, &s.UserID, &s.Status, &s.Answer, &s.FileURL, &s.Review, &s.Score, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *PostgresRepo) AddPeerReview(ctx context.Context, id uuid.UUID, peerID uuid.UUID, comment string) error {
	_, err := db.DB().Exec(ctx, `
		UPDATE homework_submissions
		SET review = COALESCE(review, '') || E'\n[peer-review ' || $2 || ']: ' || $3,
		    updated_at = NOW()
		WHERE id = $1
	`, id, peerID.String(), comment)
	return err
}

func (r *PostgresRepo) GetByUserAndHomework(ctx context.Context, userID, homeworkID uuid.UUID) (*Submission, error) {
	row := db.DB().QueryRow(ctx, `
		SELECT id, homework_id, user_id, status, answer, file_url, review, score, created_at, updated_at
		FROM homework_submissions
		WHERE user_id = $1 AND homework_id = $2
		ORDER BY created_at DESC
		LIMIT 1
	`, userID, homeworkID)

	var s Submission
	err := row.Scan(&s.ID, &s.HomeworkID, &s.UserID, &s.Status, &s.Answer, &s.FileURL, &s.Review, &s.Score, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &s, nil
}
