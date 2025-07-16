package test

import (
	"time"

	"github.com/google/uuid"
)

type Test struct {
	ID         uuid.UUID
	LessonID   uuid.UUID
	Title      string
	TimeLimit  int
	Shuffle    bool
	Attempts   int
	ShowScore  bool
	ShowAnswer bool
	AccessFrom *time.Time
	AccessTo   *time.Time
	AuthorID   uuid.UUID
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time
}
