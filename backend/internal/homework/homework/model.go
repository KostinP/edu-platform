package homework

import (
	"time"

	"github.com/google/uuid"
)

type Homework struct {
	ID          uuid.UUID
	Title       string
	Description string
	Type        string // "file", "text", "link", "code"
	CourseID    *uuid.UUID
	ModuleID    *uuid.UUID
	LessonID    *uuid.UUID
	GroupID     *uuid.UUID
	UserID      *uuid.UUID
	AuthorID    uuid.UUID
	DueAt       *time.Time
	IsRequired  bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}
