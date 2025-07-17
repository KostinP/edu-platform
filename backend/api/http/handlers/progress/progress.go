package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/kostinp/edu-platform-backend/internal/progress"
)

type ProgressHandler struct {
	progressRepo *progress.ProgressRepo
}

func NewProgressHandler(pr *progress.ProgressRepo) *ProgressHandler {
	return &ProgressHandler{progressRepo: pr}
}

func (h *ProgressHandler) GetCourseProgress(c echo.Context) error {
	userIDStr := c.QueryParam("user_id")
	courseIDStr := c.Param("course_id")

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user_id"})
	}

	courseID, err := uuid.Parse(courseIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid course_id"})
	}

	progress, err := h.progressRepo.GetCourseProgress(c.Request().Context(), userID, courseID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to get progress"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"user_id":   userID,
		"course_id": courseID,
		"progress":  progress, // 0.0 - 1.0
	})
}
