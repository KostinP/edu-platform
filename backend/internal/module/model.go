package module

import (
	"time"

	"github.com/google/uuid"
)

type Module struct {
	ID          uuid.UUID
	CourseID    uuid.UUID
	Title       string
	Description string
	Ordinal     int
	AuthorID    uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}
