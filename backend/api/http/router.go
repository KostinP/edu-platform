package http

import (
	"github.com/kostinp/edu-platform-backend/api/http/handlers"
	"github.com/kostinp/edu-platform-backend/internal/user"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo) {
	userRepo := user.NewPostgresRepo()

	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Server is running!")
	})

	e.POST("/auth/telegram", handlers.AuthTelegramHandler(userRepo))
}
