package lesson

import (
	"time"

	"github.com/google/uuid"
)

type Lesson struct {
	ID        uuid.UUID
	ModuleID  uuid.UUID
	Title     string
	Content   string
	Ordinal   int
	AuthorID  uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
