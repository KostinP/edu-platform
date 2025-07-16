package course

import (
	"time"

	"github.com/google/uuid"
)

type Course struct {
	ID          uuid.UUID
	Slug        string
	Title       string
	Description string
	AuthorID    uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}
