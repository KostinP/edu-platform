package user

import "github.com/google/uuid"

// AuthTelegramRequestDTO представляет структуру запроса авторизации через Telegram.
// swagger:model AuthTelegramRequestDTO
type AuthTelegramRequestDTO struct {
	// Telegram ID пользователя
	// required: true
	// example: "123456789"
	TelegramID string `json:"telegram_id" validate:"required"`

	// Имя пользователя
	// example: "Pavel"
	FirstName string `json:"first_name"`

	// Фамилия пользователя
	// example: "Kostin"
	LastName string `json:"last_name"`

	// Логин пользователя в Telegram
	// example: "kostinp"
	Username string `json:"username"`

	// URL фото пользователя
	// example: "https://t.me/photo.jpg"
	PhotoURL string `json:"photo_url"`

	// Email пользователя
	// example: "user@example.com"
	Email *string `json:"email"`

	// Подписка на рассылку
	// example: true
	SubscribeToNewsletter bool `json:"subscribe_to_newsletter"`

	// Роль пользователя (например, "student")
	// required: true
	// example: "student"
	Role string `json:"role" validate:"required"`
}

// AuthTelegramResponseDTO представляет успешный ответ авторизации.
// swagger:model AuthTelegramResponseDTO
type AuthTelegramResponseDTO struct {
	// Статус ответа
	// example: "ok"
	Status string `json:"status"`

	// ID пользователя
	// example: "e7bfc72d-34cf-4d92-8de1-f349ae5f0370"
	UserID uuid.UUID `json:"user_id"`
}

// ErrorResponseDTO представляет структуру ошибки API.
// swagger:model ErrorResponseDTO
type ErrorResponseDTO struct {
	// Сообщение об ошибке
	// example: "invalid request"
	Error string `json:"error"`
}
