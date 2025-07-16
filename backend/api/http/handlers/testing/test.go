package handlers

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/kostinp/edu-platform-backend/internal/testing/session"
	"github.com/kostinp/edu-platform-backend/internal/testing/test"
)

type CreateTestRequest struct {
	Title      string     `json:"title"`
	TimeLimit  int        `json:"time_limit"`
	Shuffle    bool       `json:"shuffle"`
	Attempts   int        `json:"attempts"`
	ShowScore  bool       `json:"show_score"`
	ShowAnswer bool       `json:"show_answer"`
	AccessFrom *time.Time `json:"access_from"`
	AccessTo   *time.Time `json:"access_to"`
}

type UpdateTestRequest struct {
	Title      string     `json:"title"`
	TimeLimit  int        `json:"time_limit"`
	Shuffle    bool       `json:"shuffle"`
	Attempts   int        `json:"attempts"`
	ShowScore  bool       `json:"show_score"`
	ShowAnswer bool       `json:"show_answer"`
	AccessFrom *time.Time `json:"access_from"`
	AccessTo   *time.Time `json:"access_to"`
}

type StartTestRequest struct {
	UserID uuid.UUID `json:"user_id"`
}

type SubmitTestRequest struct {
	SessionID uuid.UUID `json:"session_id"`
	Score     float64   `json:"score"`
}

func CreateTestHandler(repo test.Repository) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req CreateTestRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid body"})
		}

		t := &test.Test{
			ID:         uuid.New(),
			Title:      req.Title,
			TimeLimit:  req.TimeLimit,
			Shuffle:    req.Shuffle,
			Attempts:   req.Attempts,
			ShowScore:  req.ShowScore,
			ShowAnswer: req.ShowAnswer,
			AccessFrom: req.AccessFrom,
			AccessTo:   req.AccessTo,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		if err := repo.Create(c.Request().Context(), t); err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to create test"})
		}

		return c.JSON(http.StatusCreated, echo.Map{"id": t.ID})
	}
}

func UpdateTestHandler(repo test.Repository) echo.HandlerFunc {
	return func(c echo.Context) error {
		testID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			return c.JSON(400, echo.Map{"error": "invalid test id"})
		}

		var req UpdateTestRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(400, echo.Map{"error": "invalid body"})
		}

		t := &test.Test{
			ID:         testID,
			Title:      req.Title,
			TimeLimit:  req.TimeLimit,
			Shuffle:    req.Shuffle,
			Attempts:   req.Attempts,
			ShowScore:  req.ShowScore,
			ShowAnswer: req.ShowAnswer,
			AccessFrom: req.AccessFrom,
			AccessTo:   req.AccessTo,
			UpdatedAt:  time.Now(),
		}

		if err := repo.Update(c.Request().Context(), t); err != nil {
			return c.JSON(500, echo.Map{"error": "update failed"})
		}

		return c.JSON(200, echo.Map{"status": "updated"})
	}
}

func DeleteTestHandler(repo test.Repository) echo.HandlerFunc {
	return func(c echo.Context) error {
		testID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			return c.JSON(400, echo.Map{"error": "invalid test id"})
		}

		if err := repo.Delete(c.Request().Context(), testID); err != nil {
			return c.JSON(500, echo.Map{"error": "delete failed"})
		}

		return c.JSON(200, echo.Map{"status": "deleted"})
	}
}

func StartTestHandler(repo session.Repository) echo.HandlerFunc {
	return func(c echo.Context) error {
		testID, _ := uuid.Parse(c.Param("id"))
		var req StartTestRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(400, echo.Map{"error": "bad request"})
		}

		s := &session.TestSession{
			ID:       uuid.New(),
			UserID:   req.UserID,
			TestID:   testID,
			Attempts: 1,
			Status:   "in_progress",
		}

		err := repo.Start(c.Request().Context(), s)
		if err != nil {
			return c.JSON(500, echo.Map{"error": "failed to start test"})
		}

		return c.JSON(201, echo.Map{"session_id": s.ID})
	}
}

func SubmitTestHandler(repo session.Repository) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req SubmitTestRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(400, echo.Map{"error": "bad request"})
		}

		err := repo.Finish(c.Request().Context(), req.SessionID, req.Score)
		if err != nil {
			return c.JSON(500, echo.Map{"error": "failed to submit test"})
		}

		return c.JSON(200, echo.Map{"status": "submitted"})
	}
}
