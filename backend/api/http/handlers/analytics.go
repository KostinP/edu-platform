package handlers

import (
	"github.com/google/uuid"
	"github.com/kostinp/edu-platform-backend/internal/answer"
	"github.com/labstack/echo/v4"
)

func GetQuestionAnalyticsHandler(answerRepo answer.Repository) echo.HandlerFunc {
	return func(c echo.Context) error {
		questionID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			return c.JSON(400, echo.Map{"error": "invalid id"})
		}

		stats, err := answerRepo.GetAnswerStats(c.Request().Context(), questionID)
		if err != nil {
			return c.JSON(500, echo.Map{"error": "stats error"})
		}

		return c.JSON(200, stats)
	}
}
