package user

import (
	"context"
	"errors"

	"github.com/kostinp/edu-platform-backend/pkg/db"

	"github.com/jackc/pgx/v5"
)

var ErrUserNotFound = errors.New("user not found")

type Repository interface {
	GetByTelegramID(ctx context.Context, telegramID string) (*User, error)
	CreateOrUpdate(ctx context.Context, u *User) error
}

type PostgresRepo struct{}

func NewPostgresRepo() Repository {
	return &PostgresRepo{}
}

func (r *PostgresRepo) GetByTelegramID(ctx context.Context, telegramID string) (*User, error) {
	row := db.Pool.QueryRow(ctx, `
		SELECT id, telegram_id, first_name, last_name, username, photo_url, created_at
		FROM users
		WHERE telegram_id = $1
	`, telegramID)

	var u User
	err := row.Scan(&u.ID, &u.TelegramID, &u.FirstName, &u.LastName, &u.Username, &u.PhotoURL, &u.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *PostgresRepo) CreateOrUpdate(ctx context.Context, u *User) error {
	_, err := db.Pool.Exec(ctx, `
		INSERT INTO users (id, telegram_id, first_name, last_name, username, photo_url, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (telegram_id) DO UPDATE SET
			first_name = EXCLUDED.first_name,
			last_name = EXCLUDED.last_name,
			username = EXCLUDED.username,
			photo_url = EXCLUDED.photo_url
	`, u.ID, u.TelegramID, u.FirstName, u.LastName, u.Username, u.PhotoURL, u.CreatedAt)
	return err
}
