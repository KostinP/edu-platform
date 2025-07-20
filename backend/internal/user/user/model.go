package user

import (
	"time"

	"github.com/google/uuid"
)

// User представляет пользователя платформы.
type User struct {
	ID                    uuid.UUID // Уникальный идентификатор
	TelegramID            string    // Telegram ID (уникальный идентификатор внешней системы)
	FirstName             string    // Имя пользователя
	LastName              string    // Фамилия пользователя
	Username              string    // Никнейм
	PhotoURL              string    // Ссылка на аватар
	CreatedAt             time.Time // Дата создания
	Email                 *string   // Email (опционально)
	SubscribeToNewsletter bool      // Подписка на рассылку
	Role                  string    // Роль (например, "student", "teacher", "admin")
}
