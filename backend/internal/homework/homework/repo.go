package homework

import (
	"context"
	"time"

	"github.com/kostinp/edu-platform-backend/pkg/db"

	"github.com/google/uuid"
)

type PostgresRepo struct{}

type Repository interface {
	Create(ctx context.Context, hw *Homework) error
	GetByID(ctx context.Context, id uuid.UUID) (*Homework, error)
	List(ctx context.Context) ([]Homework, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

func NewRepository() Repository {
	return &PostgresRepo{}
}

func (r *PostgresRepo) Create(ctx context.Context, hw *Homework) error {
	_, err := db.Pool.Exec(ctx, `
		INSERT INTO homeworks
		(id, title, description, type, course_id, module_id, lesson_id, group_id, user_id, author_id, due_at, is_required, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)
	`, hw.ID, hw.Title, hw.Description, hw.Type, hw.CourseID, hw.ModuleID, hw.LessonID, hw.GroupID, hw.UserID, hw.AuthorID, hw.DueAt, hw.IsRequired, time.Now(), time.Now())
	return err
}

func (r *PostgresRepo) GetByID(ctx context.Context, id uuid.UUID) (*Homework, error) {
	row := db.Pool.QueryRow(ctx, `
		SELECT id, title, description, type, course_id, module_id, lesson_id, group_id, user_id, author_id, due_at, is_required, created_at, updated_at
		FROM homeworks WHERE id = $1 AND deleted_at IS NULL
	`, id)

	var hw Homework
	err := row.Scan(&hw.ID, &hw.Title, &hw.Description, &hw.Type, &hw.CourseID, &hw.ModuleID, &hw.LessonID, &hw.GroupID, &hw.UserID, &hw.AuthorID, &hw.DueAt, &hw.IsRequired, &hw.CreatedAt, &hw.UpdatedAt)
	return &hw, err
}

func (r *PostgresRepo) List(ctx context.Context) ([]Homework, error) {
	rows, err := db.Pool.Query(ctx, `
		SELECT id, title, description, type, course_id, module_id, lesson_id, group_id, user_id, author_id, due_at, is_required, created_at, updated_at
		FROM homeworks WHERE deleted_at IS NULL
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []Homework
	for rows.Next() {
		var hw Homework
		err := rows.Scan(&hw.ID, &hw.Title, &hw.Description, &hw.Type, &hw.CourseID, &hw.ModuleID, &hw.LessonID, &hw.GroupID, &hw.UserID, &hw.AuthorID, &hw.DueAt, &hw.IsRequired, &hw.CreatedAt, &hw.UpdatedAt)
		if err != nil {
			return nil, err
		}
		result = append(result, hw)
	}
	return result, nil
}

func (r *PostgresRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := db.Pool.Exec(ctx, `
		UPDATE homeworks SET deleted_at = NOW() WHERE id = $1
	`, id)
	return err
}
