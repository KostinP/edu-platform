package user

import "context"

type AuthTelegramRequest struct {
	TelegramID            string
	FirstName             string
	LastName              string
	Username              string
	PhotoURL              string
	Email                 *string
	SubscribeToNewsletter bool
	Role                  string
}

type Service interface {
	AuthTelegram(ctx context.Context, req AuthTelegramRequest) (*User, error)
}
