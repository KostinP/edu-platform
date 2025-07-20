package user

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) AuthTelegram(ctx context.Context, req AuthTelegramRequest) (*User, error) {
	existing, err := s.repo.GetByTelegramID(ctx, req.TelegramID)
	if err != nil && err != ErrUserNotFound {
		return nil, err
	}

	var u *User
	if existing == nil {
		u = &User{
			ID:                    uuid.New(),
			TelegramID:            req.TelegramID,
			FirstName:             req.FirstName,
			LastName:              req.LastName,
			Username:              req.Username,
			PhotoURL:              req.PhotoURL,
			Email:                 req.Email,
			SubscribeToNewsletter: req.SubscribeToNewsletter,
			Role:                  req.Role,
			CreatedAt:             time.Now(),
		}
	} else {
		u = existing
		u.FirstName = req.FirstName
		u.LastName = req.LastName
		u.Username = req.Username
		u.PhotoURL = req.PhotoURL
		u.Email = req.Email
		u.SubscribeToNewsletter = req.SubscribeToNewsletter
		u.Role = req.Role
	}

	if err := s.repo.CreateOrUpdate(ctx, u); err != nil {
		return nil, err
	}

	return u, nil
}
