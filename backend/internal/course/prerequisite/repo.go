package prerequisite

import (
	"context"

	"github.com/google/uuid"
	"github.com/kostinp/edu-platform-backend/pkg/db"
)

type Repository interface {
	Add(ctx context.Context, courseID, prerequisiteID uuid.UUID) error
	Remove(ctx context.Context, courseID, prerequisiteID uuid.UUID) error
	List(ctx context.Context, courseID uuid.UUID) ([]uuid.UUID, error)
}

type PostgresRepo struct{}

func NewPostgresRepo() Repository {
	return &PostgresRepo{}
}

func (r *PostgresRepo) Add(ctx context.Context, courseID, prerequisiteID uuid.UUID) error {
	_, err := db.DB().Exec(ctx, `
		INSERT INTO course_prerequisites (course_id, prerequisite_id)
		VALUES ($1, $2) ON CONFLICT DO NOTHING
	`, courseID, prerequisiteID)
	return err
}

func (r *PostgresRepo) Remove(ctx context.Context, courseID, prerequisiteID uuid.UUID) error {
	_, err := db.DB().Exec(ctx, `
		DELETE FROM course_prerequisites
		WHERE course_id = $1 AND prerequisite_id = $2
	`, courseID, prerequisiteID)
	return err
}

func (r *PostgresRepo) List(ctx context.Context, courseID uuid.UUID) ([]uuid.UUID, error) {
	rows, err := db.DB().Query(ctx, `
		SELECT prerequisite_id FROM course_prerequisites
		WHERE course_id = $1
	`, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prerequisites []uuid.UUID
	for rows.Next() {
		var pid uuid.UUID
		if err := rows.Scan(&pid); err != nil {
			return nil, err
		}
		prerequisites = append(prerequisites, pid)
	}
	return prerequisites, nil
}
