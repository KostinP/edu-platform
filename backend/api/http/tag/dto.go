package tag

import (
	"github.com/google/uuid"
	"github.com/kostinp/edu-platform-backend/internal/tag"
)

// CreateTagRequest представляет тело запроса для создания тега
type CreateTagRequest struct {
	// Name — название тега
	Name string `json:"name" example:"golang"`

	// AuthorID — UUID автора тега, обязателен
	AuthorID uuid.UUID `json:"author_id" example:"550e8400-e29b-41d4-a716-446655440000"`
}

// UpdateTagRequest представляет тело запроса для обновления тега
type UpdateTagRequest struct {
	// Name — новое название тега
	Name string `json:"name" example:"go-programming"`
}

// TagResponse представляет данные тега для ответа API
type TagResponse struct {
	// ID — UUID тега
	ID uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`

	// Name — название тега
	Name string `json:"name" example:"golang"`

	// AuthorID — UUID автора тега
	AuthorID uuid.UUID `json:"author_id" example:"550e8400-e29b-41d4-a716-446655440000"`

	// CreatedAt — дата и время создания тега в формате RFC3339
	CreatedAt string `json:"created_at" example:"2025-07-20T15:04:05Z07:00"`

	// UpdatedAt — дата и время последнего обновления тега в формате RFC3339
	UpdatedAt string `json:"updated_at" example:"2025-07-20T15:04:05Z07:00"`
}

// ToTagResponse преобразует внутреннюю модель Tag в структуру TagResponse для API
func ToTagResponse(t *tag.Tag) TagResponse {
	return TagResponse{
		ID:        t.ID,
		Name:      t.Name,
		AuthorID:  t.AuthorID,
		CreatedAt: t.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: t.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// ToTagResponseList преобразует срез внутренних моделей Tag в срез TagResponse
func ToTagResponseList(tags []*tag.Tag) []TagResponse {
	resp := make([]TagResponse, 0, len(tags))
	for _, t := range tags {
		resp = append(resp, ToTagResponse(t))
	}
	return resp
}
