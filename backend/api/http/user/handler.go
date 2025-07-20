package user

import (
	"net/http"

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
		if err := c.Bind(&reqDTO); err != nil {
			return c.JSON(http.StatusBadRequest, ErrorResponseDTO{Error: "invalid request"})
		}

		// Здесь можно добавить валидацию, если подключен echo/validator

		// Конвертация DTO HTTP в внутренний DTO
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
			return c.JSON(http.StatusInternalServerError, ErrorResponseDTO{Error: "failed to save user"})
		}

		return c.JSON(http.StatusOK, AuthTelegramResponseDTO{
			Status: "ok",
			UserID: u.ID,
		})
	}
}
