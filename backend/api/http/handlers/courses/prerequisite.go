package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/kostinp/edu-platform-backend/internal/course/prerequisite"

	"github.com/labstack/echo/v4"
)

type PrerequisiteRequest struct {
	PrerequisiteID string `json:"prerequisite_id"`
}

func AddCoursePrerequisiteHandler(repo prerequisite.Repository) echo.HandlerFunc {
	return func(c echo.Context) error {
		courseID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid course ID"})
		}

		var req PrerequisiteRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request"})
		}

		prereqID, err := uuid.Parse(req.PrerequisiteID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid prerequisite ID"})
		}

		if err := repo.Add(c.Request().Context(), courseID, prereqID); err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to add prerequisite"})
		}

		return c.JSON(http.StatusCreated, echo.Map{"status": "added"})
	}
}

func RemoveCoursePrerequisiteHandler(repo prerequisite.Repository) echo.HandlerFunc {
	return func(c echo.Context) error {
		courseID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid course ID"})
		}

		var req PrerequisiteRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request"})
		}

		prereqID, err := uuid.Parse(req.PrerequisiteID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid prerequisite ID"})
		}

		if err := repo.Remove(c.Request().Context(), courseID, prereqID); err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to remove prerequisite"})
		}

		return c.JSON(http.StatusOK, echo.Map{"status": "removed"})
	}
}

func ListCoursePrerequisitesHandler(repo prerequisite.Repository) echo.HandlerFunc {
	return func(c echo.Context) error {
		courseID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid course ID"})
		}

		list, err := repo.List(c.Request().Context(), courseID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to list prerequisites"})
		}

		return c.JSON(http.StatusOK, list)
	}
}
