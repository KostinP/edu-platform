package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kostinp/edu-platform-backend/pkg/db"

	"github.com/jackc/pgx/v5"
)

// PostgresRepo — реализация Repository для PostgreSQL.
type PostgresRepo struct{}

// NewPostgresRepo возвращает новый экземпляр PostgreSQL репозитория.
func NewPostgresRepo() Repository {
	return &PostgresRepo{}
}

func (r *PostgresRepo) GetByTelegramID(ctx context.Context, telegramID string) (*User, error) {
	row := db.DB().QueryRow(ctx, `
		SELECT id, telegram_id, first_name, last_name, username, photo_url, created_at, email, subscribe_to_newsletter, role
		FROM users
		WHERE telegram_id = $1
	`, telegramID)

	var u User
	err := row.Scan(&u.ID, &u.TelegramID, &u.FirstName, &u.LastName, &u.Username, &u.PhotoURL, &u.CreatedAt, &u.Email, &u.SubscribeToNewsletter, &u.Role)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *PostgresRepo) CreateOrUpdate(ctx context.Context, u *User) error {
	// Если ID не задан — генерируем
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}

	// Подстраховка: если CreatedAt не установлен, ставим текущее время
	if u.CreatedAt.IsZero() {
		u.CreatedAt = time.Now()
	}

	query := `
		INSERT INTO users (
			id, telegram_id, first_name, last_name, username, photo_url,
			created_at, email, subscribe_to_newsletter, role
		) VALUES (
			$1, $2, $3, $4, $5, $6,
			$7, $8, $9, $10
		)
		ON CONFLICT (telegram_id) DO UPDATE SET
			first_name = EXCLUDED.first_name,
			last_name = EXCLUDED.last_name,
			username = EXCLUDED.username,
			photo_url = EXCLUDED.photo_url,
			email = EXCLUDED.email,
			subscribe_to_newsletter = EXCLUDED.subscribe_to_newsletter,
			role = EXCLUDED.role
	`

	_, err := db.DB().Exec(ctx, query,
		u.ID,
		u.TelegramID,
		u.FirstName,
		u.LastName,
		u.Username,
		u.PhotoURL,
		u.CreatedAt,
		u.Email,
		u.SubscribeToNewsletter,
		u.Role,
	)

	if err != nil {
		fmt.Printf("❌ CreateOrUpdate error for telegram_id=%s: %v\n", u.TelegramID, err)
		return err
	}

	fmt.Printf("✅ User %s (%s) created/updated\n", u.TelegramID, u.Username)
	return nil
}
