package tag

import (
	"time"

	"github.com/google/uuid"
)

type Tag struct {
	ID        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	AuthorID  uuid.UUID  `json:"author_id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-"`
}

type UpdateTagDTO struct {
	Name string `json:"name"`
}
