package course

import (
	"context"

	"github.com/google/uuid"
	"github.com/kostinp/edu-platform-backend/pkg/db"
)

type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*Course, error)
	GetAll(ctx context.Context) ([]Course, error)
	GetStats(ctx context.Context, courseID uuid.UUID) (*Stats, error)
	Create(ctx context.Context, c *Course) error
	GetBySlug(ctx context.Context, slug string) (*Course, error)
	Update(ctx context.Context, c *Course) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
}

type PostgresRepo struct{}

func NewPostgresRepo() Repository {
	return &PostgresRepo{}
}

func (r *PostgresRepo) GetByID(ctx context.Context, id uuid.UUID) (*Course, error) {
	row := db.DB().QueryRow(ctx, `
		SELECT id, slug, title, description, author_id, created_at, updated_at, deleted_at
		FROM courses
		WHERE id = $1
	`, id)

	var c Course
	err := row.Scan(&c.ID, &c.Slug, &c.Title, &c.Description, &c.AuthorID, &c.CreatedAt, &c.UpdatedAt, &c.DeletedAt)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *PostgresRepo) GetAll(ctx context.Context) ([]Course, error) {
	rows, err := db.DB().Query(ctx, `
		SELECT id, slug, title, description, author_id, created_at
		FROM courses
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	courses := make([]Course, 0)
	for rows.Next() {
		var c Course
		if err := rows.Scan(&c.ID, &c.Slug, &c.Title, &c.Description, &c.AuthorID, &c.CreatedAt); err != nil {
			return nil, err
		}
		courses = append(courses, c)
	}
	return courses, nil
}

func (r *PostgresRepo) Create(ctx context.Context, c *Course) error {
	_, err := db.DB().Exec(ctx, `
		INSERT INTO courses (id, slug, title, description, author_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, c.ID, c.Slug, c.Title, c.Description, c.AuthorID, c.CreatedAt)
	return err
}

func (r *PostgresRepo) GetBySlug(ctx context.Context, slug string) (*Course, error) {
	row := db.DB().QueryRow(ctx, `
		SELECT id, slug, title, description, author_id, created_at
		FROM courses
		WHERE slug = $1
	`, slug)

	var c Course
	err := row.Scan(&c.ID, &c.Slug, &c.Title, &c.Description, &c.AuthorID, &c.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *PostgresRepo) Update(ctx context.Context, c *Course) error {
	_, err := db.DB().Exec(ctx, `
		UPDATE courses SET
			slug = $1,
			title = $2,
			description = $3,
			updated_at = NOW()
		WHERE id = $4 AND deleted_at IS NULL
	`, c.Slug, c.Title, c.Description, c.ID)
	return err
}

func (r *PostgresRepo) SoftDelete(ctx context.Context, id uuid.UUID) error {
	_, err := db.DB().Exec(ctx, `
		UPDATE courses SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`, id)
	return err
}

func (r *PostgresRepo) GetStats(ctx context.Context, courseID uuid.UUID) (*Stats, error) {
	query := `
		SELECT
			COALESCE(AVG(r.rating), 0) as rating,
			(SELECT COUNT(*) FROM lessons l
				JOIN modules m ON l.module_id = m.id
				WHERE m.course_id = $1) as lessons_count,
			(SELECT COALESCE(SUM(l.duration), 0) FROM lessons l
				JOIN modules m ON l.module_id = m.id
				WHERE m.course_id = $1) as total_duration,
			(SELECT COUNT(DISTINCT e.user_id) FROM enrollments e
				WHERE e.course_id = $1) as students
		FROM course_reviews r
		WHERE r.course_id = $1
	`

	var s Stats
	err := db.DB().QueryRow(ctx, query, courseID).Scan(
		&s.Rating, &s.LessonsCount, &s.Duration, &s.Students,
	)
	if err != nil {
		return nil, err
	}
	return &s, nil
}
