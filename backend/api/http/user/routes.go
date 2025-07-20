package user

import "github.com/labstack/echo/v4"

func RegisterRoutes(e *echo.Echo) {
	group := e.Group("/auth")
	group.POST("/telegram", AuthTelegramHandler(NewUserService()))
}
