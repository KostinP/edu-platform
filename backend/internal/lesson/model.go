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
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
