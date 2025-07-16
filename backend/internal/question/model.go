package question

import (
	"time"

	"github.com/google/uuid"
)

type Question struct {
	ID        uuid.UUID
	TestID    uuid.UUID
	Type      string
	Title     string
	ImageURL  *string
	Data      []byte // json.RawMessage
	Feedback  string
	Score     float64
	Ordinal   int
	AuthorID  uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
