package course

import (
	"context"

	"github.com/kostinp/edu-platform-backend/pkg/db"
)

type Repository interface {
	GetAll(ctx context.Context) ([]Course, error)
	Create(ctx context.Context, c *Course) error
	GetBySlug(ctx context.Context, slug string) (*Course, error)
}

type PostgresRepo struct{}

func NewPostgresRepo() Repository {
	return &PostgresRepo{}
}

func (r *PostgresRepo) GetAll(ctx context.Context) ([]Course, error) {
	rows, err := db.Pool.Query(ctx, `
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
	_, err := db.Pool.Exec(ctx, `
		INSERT INTO courses (id, slug, title, description, author_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, c.ID, c.Slug, c.Title, c.Description, c.AuthorID, c.CreatedAt)
	return err
}

func (r *PostgresRepo) GetBySlug(ctx context.Context, slug string) (*Course, error) {
	row := db.Pool.QueryRow(ctx, `
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
