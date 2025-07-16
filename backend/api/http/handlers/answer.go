package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/kostinp/edu-platform-backend/internal/answer"
	"github.com/kostinp/edu-platform-backend/internal/question"
)

type AnswerRequest struct {
	SessionID       uuid.UUID       `json:"session_id"`
	UserID          uuid.UUID       `json:"user_id"`
	QuestionID      uuid.UUID       `json:"question_id"`
	Answer          json.RawMessage `json:"answer"`
	DurationSeconds int             `json:"duration_seconds"`
}

type ManualReviewRequest struct {
	MarkAsCorrect bool `json:"mark_as_correct"`
}

func SubmitAnswerHandler(answerRepo answer.Repository, questionRepo question.Repository) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req AnswerRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "bad request"})
		}

		// Получить вопрос из базы
		q, err := questionRepo.GetByID(c.Request().Context(), req.QuestionID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "question not found"})
		}

		// Автоматическая проверка
		correct, _, needsManual := answer.EvaluateAnswer(q.Type, q.Data, req.Answer)

		isChecked := !needsManual

		// Сохранить ответ
		a := &answer.UserAnswer{
			ID:              uuid.New(),
			TestSessionID:   req.SessionID,
			UserID:          req.UserID,
			QuestionID:      req.QuestionID,
			Answer:          req.Answer,
			DurationSeconds: req.DurationSeconds,
			IsChecked:       isChecked,
			CreatedAt:       time.Now(),
		}

		if err := answerRepo.Submit(c.Request().Context(), a); err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to submit answer"})
		}

		return c.JSON(http.StatusCreated, echo.Map{
			"status":      "answer saved",
			"autoChecked": isChecked,
			"correct":     correct,
		})
	}
}

func ReviewAnswerHandler(repo answer.Repository) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, _ := uuid.Parse(c.Param("id"))
		var req ManualReviewRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(400, echo.Map{"error": "bad request"})
		}
		if err := repo.ManualCheck(c.Request().Context(), id, req.MarkAsCorrect); err != nil {
			return c.JSON(500, echo.Map{"error": "review failed"})
		}
		return c.JSON(200, echo.Map{"status": "checked"})
	}
}
