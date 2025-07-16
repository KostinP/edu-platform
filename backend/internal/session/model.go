package session

import (
	"time"

	"github.com/google/uuid"
)

type TestSession struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	TestID     uuid.UUID
	StartedAt  time.Time
	FinishedAt *time.Time
	Score      *float64
	Attempts   int
	Status     string // in_progress | finished | canceled
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
