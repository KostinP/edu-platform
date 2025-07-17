package group

import (
	"time"

	"github.com/google/uuid"
)

type Group struct {
	ID          uuid.UUID
	Name        string
	Description string
	OwnerID     uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
