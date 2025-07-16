package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/kostinp/edu-platform-backend/internal/question"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type CreateQuestionRequest struct {
	Type     string          `json:"type"`
	Title    string          `json:"title"`
	ImageURL *string         `json:"image_url"`
	Data     json.RawMessage `json:"data"`
	Feedback string          `json:"feedback"`
	Score    float64         `json:"score"`
	Ordinal  int             `json:"ordinal"`
}

type UpdateQuestionRequest struct {
	Type     string          `json:"type"`
	Title    string          `json:"title"`
	ImageURL *string         `json:"image_url"`
	Data     json.RawMessage `json:"data"`
	Feedback string          `json:"feedback"`
	Score    float64         `json:"score"`
	Ordinal  int             `json:"ordinal"`
}

func CreateQuestionHandler(repo question.Repository) echo.HandlerFunc {
	return func(c echo.Context) error {
		testID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid test id"})
		}

		var req CreateQuestionRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid body"})
		}

		q := &question.Question{
			ID:        uuid.New(),
			TestID:    testID,
			Type:      req.Type,
			Title:     req.Title,
			ImageURL:  req.ImageURL,
			Data:      req.Data,
			Feedback:  req.Feedback,
			Score:     req.Score,
			Ordinal:   req.Ordinal,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := repo.Create(c.Request().Context(), q); err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to create question"})
		}

		return c.JSON(http.StatusCreated, echo.Map{"id": q.ID})
	}
}

func GetQuestionsByTestHandler(repo question.Repository) echo.HandlerFunc {
	return func(c echo.Context) error {
		testID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			return c.JSON(400, echo.Map{"error": "invalid test id"})
		}

		questions, err := repo.GetByTestID(c.Request().Context(), testID)
		if err != nil {
			return c.JSON(500, echo.Map{"error": "fetch failed"})
		}

		return c.JSON(200, questions)
	}
}

func UpdateQuestionHandler(repo question.Repository) echo.HandlerFunc {
	return func(c echo.Context) error {
		questionID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			return c.JSON(400, echo.Map{"error": "invalid id"})
		}

		var req UpdateQuestionRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(400, echo.Map{"error": "invalid body"})
		}

		q := &question.Question{
			ID:        questionID,
			Type:      req.Type,
			Title:     req.Title,
			ImageURL:  req.ImageURL,
			Data:      req.Data,
			Feedback:  req.Feedback,
			Score:     req.Score,
			Ordinal:   req.Ordinal,
			UpdatedAt: time.Now(),
		}

		if err := repo.Update(c.Request().Context(), q); err != nil {
			return c.JSON(500, echo.Map{"error": "update failed"})
		}

		return c.JSON(200, echo.Map{"status": "updated"})
	}
}

func DeleteQuestionHandler(repo question.Repository) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			return c.JSON(400, echo.Map{"error": "invalid id"})
		}
		if err := repo.Delete(c.Request().Context(), id); err != nil {
			return c.JSON(500, echo.Map{"error": "delete failed"})
		}
		return c.JSON(200, echo.Map{"status": "deleted"})
	}
}
