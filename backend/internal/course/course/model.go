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
	Category    string
	Price       int
	ImageURL    *string
	AuthorID    uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

type Stats struct {
	Rating       float64
	LessonsCount int
	Duration     int
	Students     int
}
