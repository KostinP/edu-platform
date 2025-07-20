package user

import (
	"context"
)

// Repository определяет интерфейс доступа к данным пользователя.
type Repository interface {
	// GetByTelegramID возвращает пользователя по Telegram ID.
	GetByTelegramID(ctx context.Context, telegramID string) (*User, error)

	// CreateOrUpdate сохраняет или обновляет пользователя по Telegram ID.
	CreateOrUpdate(ctx context.Context, u *User) error
}
