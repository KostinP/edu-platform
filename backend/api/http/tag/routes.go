package tag

import (
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(g *echo.Group, h *Handler) {
	g.POST("/tags", h.Create)
	g.GET("/tags", h.GetAll)
	g.PUT("/tags/:id", h.Update)
	g.DELETE("/tags/:id", h.Delete)
}
