package submission

import (
	"time"

	"github.com/google/uuid"
)

type Submission struct {
	ID         uuid.UUID
	HomeworkID uuid.UUID
	UserID     uuid.UUID
	Status     string
	Answer     string
	FileURL    *string
	Review     *string
	Score      *float64
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
