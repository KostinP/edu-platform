package user

import (
	"bytes"
	"io"
	"net/http"
	"os"

	"github.com/kostinp/edu-platform-backend/internal/user/user"
	"github.com/labstack/echo/v4"
)

// AuthTelegramHandler godoc
// @Summary Авторизация через Telegram
// @Description Авторизация или обновление пользователя через Telegram ID
// @Tags auth
// @Accept json
// @Produce json
// @Param input body AuthTelegramRequestDTO true "Параметры авторизации через Telegram"
// @Success 200 {object} AuthTelegramResponseDTO
// @Failure 400 {object} ErrorResponseDTO
// @Failure 500 {object} ErrorResponseDTO
// @Router /auth/telegram [post]
func AuthTelegramHandler(svc user.Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		var reqDTO AuthTelegramRequestDTO
		bodyBytes, err := io.ReadAll(c.Request().Body)
		if err != nil {
			c.Logger().Errorf("Failed to read request body: %v", err)
		} else {
			c.Logger().Infof("Request body: %s", string(bodyBytes))
			// Восстановим тело для повторного чтения
			c.Request().Body = io.NopCloser(bytes.NewReader(bodyBytes))
		}
		if err := c.Bind(&reqDTO); err != nil {
			c.Logger().Errorf("Bind error: %v", err)
			return c.JSON(http.StatusBadRequest, ErrorResponseDTO{Error: "invalid request"})
		}

		req := user.AuthTelegramRequest{
			TelegramID:            reqDTO.TelegramID,
			FirstName:             reqDTO.FirstName,
			LastName:              reqDTO.LastName,
			Username:              reqDTO.Username,
			PhotoURL:              reqDTO.PhotoURL,
			Email:                 reqDTO.Email,
			SubscribeToNewsletter: reqDTO.SubscribeToNewsletter,
			Role:                  reqDTO.Role,
		}

		u, err := svc.AuthTelegram(c.Request().Context(), req)
		if err != nil {
			c.Logger().Errorf("AuthTelegram service error: %v", err)

			// Возвращаем ошибку с текстом только в dev среде
			if os.Getenv("ENV") == "local" {
				return c.JSON(http.StatusInternalServerError, ErrorResponseDTO{Error: err.Error()})
			}
			return c.JSON(http.StatusInternalServerError, ErrorResponseDTO{Error: "failed to save user"})
		}

		return c.JSON(http.StatusOK, AuthTelegramResponseDTO{
			Status: "ok",
			UserID: u.ID,
		})
	}
}
