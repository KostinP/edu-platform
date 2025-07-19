package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                    uuid.UUID
	TelegramID            string
	FirstName             string
	LastName              string
	Username              string
	PhotoURL              string
	CreatedAt             time.Time
	Email                 *string
	SubscribeToNewsletter bool
	Role                  string
}
