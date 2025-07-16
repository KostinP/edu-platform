package handlers

import (
	"net/http"
	"time"

	"github.com/kostinp/edu-platform-backend/internal/homework/homework"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type CreateHomeworkRequest struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Type        string     `json:"type"`
	CourseID    *uuid.UUID `json:"course_id"`
	ModuleID    *uuid.UUID `json:"module_id"`
	LessonID    *uuid.UUID `json:"lesson_id"`
	GroupID     *uuid.UUID `json:"group_id"`
	UserID      *uuid.UUID `json:"user_id"`
	AuthorID    uuid.UUID  `json:"author_id"`
	DueAt       *time.Time `json:"due_at"`
	IsRequired  bool       `json:"is_required"`
}

func CreateHomeworkHandler(repo homework.Repository) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req CreateHomeworkRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid input"})
		}

		hw := &homework.Homework{
			ID:          uuid.New(),
			Title:       req.Title,
			Description: req.Description,
			Type:        req.Type,
			CourseID:    req.CourseID,
			ModuleID:    req.ModuleID,
			LessonID:    req.LessonID,
			GroupID:     req.GroupID,
			UserID:      req.UserID,
			AuthorID:    req.AuthorID,
			DueAt:       req.DueAt,
			IsRequired:  req.IsRequired,
		}

		if err := repo.Create(c.Request().Context(), hw); err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to create"})
		}
		return c.JSON(http.StatusCreated, hw)
	}
}

func ListHomeworksHandler(repo homework.Repository) echo.HandlerFunc {
	return func(c echo.Context) error {
		hws, err := repo.List(c.Request().Context())
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to list"})
		}
		return c.JSON(http.StatusOK, hws)
	}
}
