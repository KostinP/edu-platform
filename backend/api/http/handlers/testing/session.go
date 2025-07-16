package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/kostinp/edu-platform-backend/internal/testing/answer"
	"github.com/kostinp/edu-platform-backend/internal/testing/question"
	"github.com/kostinp/edu-platform-backend/internal/testing/session"
)

func GetSessionResultHandler(
	sessionRepo session.Repository,
	answerRepo answer.Repository,
	questionRepo question.Repository,
) echo.HandlerFunc {
	return func(c echo.Context) error {
		sessionID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid session id"})
		}

		s, err := sessionRepo.GetByID(c.Request().Context(), sessionID)
		if err != nil {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "session not found"})
		}

		answers, err := answerRepo.GetBySession(c.Request().Context(), sessionID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "answers not found"})
		}

		var total float64
		var max float64
		var result []map[string]interface{}

		for _, a := range answers {
			q, err := questionRepo.GetByID(c.Request().Context(), a.QuestionID)
			if err != nil {
				continue
			}
			correct, _, _ := answer.EvaluateAnswer(q.Type, q.Data, a.Answer)

			score := 0.0
			if correct {
				score = q.Score
			}
			if a.IsChecked {
				total += score
			}
			max += q.Score

			result = append(result, map[string]interface{}{
				"question_id":     q.ID,
				"title":           q.Title,
				"user_answer":     a.Answer,
				"correct":         correct,
				"score":           score,
				"max_score":       q.Score,
				"manual_required": !a.IsChecked,
			})
		}

		return c.JSON(http.StatusOK, echo.Map{
			"score":        total,
			"max_score":    max,
			"percent":      int(total / max * 100),
			"questions":    result,
			"finished_at":  s.FinishedAt,
			"status":       s.Status,
			"show_answers": true,
		})
	}
}
