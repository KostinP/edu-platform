package answer

import (
	"time"

	"github.com/google/uuid"
)

type UserAnswer struct {
	ID              uuid.UUID
	TestSessionID   uuid.UUID
	UserID          uuid.UUID
	QuestionID      uuid.UUID
	Answer          []byte // json.RawMessage
	DurationSeconds int
	CreatedAt       time.Time
	IsChecked       bool
}
