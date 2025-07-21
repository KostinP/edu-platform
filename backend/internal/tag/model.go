package tag

import (
	"time"

	"github.com/google/uuid"
)

// Tag представляет сущность тега в системе.
// Используется для хранения и обработки данных в бизнес-логике и слое доступа к данным.
type Tag struct {
	ID        uuid.UUID  `json:"id"`         // ID — уникальный идентификатор тега
	Name      string     `json:"name"`       // Name — название тега
	AuthorID  uuid.UUID  `json:"author_id"`  // AuthorID — идентификатор пользователя, создавшего тег
	CreatedAt time.Time  `json:"created_at"` // CreatedAt — дата и время создания тега
	UpdatedAt time.Time  `json:"updated_at"` // UpdatedAt — дата и время последнего обновления тега
	DeletedAt *time.Time `json:"-"`          // DeletedAt — nullable поле для soft delete (мягкого удаления). Не сериализуется в JSON
}

// UpdateTagDTO используется для передачи данных при обновлении тега
type UpdateTagDTO struct {
	Name string `json:"name"` // Name — новое название тега
}
