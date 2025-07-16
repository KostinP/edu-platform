package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/kostinp/edu-platform-backend/internal/user"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type AuthTelegramRequest struct {
	TelegramID string `json:"telegram_id" validate:"required"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Username   string `json:"username"`
	PhotoURL   string `json:"photo_url"`
}

func AuthTelegramHandler(repo user.Repository) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req AuthTelegramRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request"})
		}

		ctx := context.Background()

		existing, err := repo.GetByTelegramID(ctx, req.TelegramID)
		if err != nil && err != user.ErrUserNotFound {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "db error"})
		}

		var newUser *user.User
		if existing == nil {
			newUser = &user.User{
				ID:         uuid.New(),
				TelegramID: req.TelegramID,
				FirstName:  req.FirstName,
				LastName:   req.LastName,
				Username:   req.Username,
				PhotoURL:   req.PhotoURL,
				CreatedAt:  time.Now(),
			}
		} else {
			newUser = existing
			newUser.FirstName = req.FirstName
			newUser.LastName = req.LastName
			newUser.Username = req.Username
			newUser.PhotoURL = req.PhotoURL
		}

		if err := repo.CreateOrUpdate(ctx, newUser); err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to save user"})
		}

		return c.JSON(http.StatusOK, echo.Map{"status": "ok", "user_id": newUser.ID})
	}
}
