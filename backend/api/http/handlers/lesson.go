package handlers

import (
	"net/http"
	"time"

	"github.com/kostinp/edu-platform-backend/internal/lesson"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type CreateLessonRequest struct {
	ModuleID string `json:"module_id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	Ordinal  int    `json:"ordinal"`
}

type UpdateLessonRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Ordinal int    `json:"ordinal"`
}

func CreateLessonHandler(repo lesson.Repository) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req CreateLessonRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(400, echo.Map{"error": "invalid request"})
		}

		moduleID, err := uuid.Parse(req.ModuleID)
		if err != nil {
			return c.JSON(400, echo.Map{"error": "invalid module_id"})
		}

		newLesson := &lesson.Lesson{
			ID:        uuid.New(),
			ModuleID:  moduleID,
			Title:     req.Title,
			Content:   req.Content,
			Ordinal:   req.Ordinal,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := repo.Create(c.Request().Context(), newLesson); err != nil {
			return c.JSON(500, echo.Map{"error": "failed to create lesson"})
		}

		return c.JSON(201, echo.Map{"status": "created", "id": newLesson.ID})
	}
}

func GetLessonsByModuleIDHandler(repo lesson.Repository) echo.HandlerFunc {
	return func(c echo.Context) error {
		moduleIDStr := c.Param("moduleID")
		moduleID, err := uuid.Parse(moduleIDStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid module ID"})
		}

		lessons, err := repo.GetByModuleID(c.Request().Context(), moduleID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to load lessons"})
		}

		return c.JSON(http.StatusOK, lessons)
	}
}

func UpdateLessonHandler(repo lesson.Repository) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, _ := uuid.Parse(c.Param("id"))
		var req UpdateLessonRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(400, echo.Map{"error": "bad request"})
		}

		err := repo.Update(c.Request().Context(), &lesson.Lesson{
			ID:      id,
			Title:   req.Title,
			Content: req.Content,
			Ordinal: req.Ordinal,
		})
		if err != nil {
			return c.JSON(500, echo.Map{"error": "update failed"})
		}
		return c.JSON(200, echo.Map{"status": "updated"})
	}
}

func DeleteLessonHandler(repo lesson.Repository) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, _ := uuid.Parse(c.Param("id"))
		if err := repo.SoftDelete(c.Request().Context(), id); err != nil {
			return c.JSON(500, echo.Map{"error": "delete failed"})
		}
		return c.JSON(200, echo.Map{"status": "deleted"})
	}
}
