package handlers

import (
	"net/http"
	"time"

	"github.com/kostinp/edu-platform-backend/internal/module"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type CreateModuleRequest struct {
	CourseID    string `json:"course_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Ordinal     int    `json:"ordinal"`
}

type UpdateModuleRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Ordinal     int    `json:"ordinal"`
}

func CreateModuleHandler(repo module.Repository) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req CreateModuleRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(400, echo.Map{"error": "invalid request"})
		}

		courseID, err := uuid.Parse(req.CourseID)
		if err != nil {
			return c.JSON(400, echo.Map{"error": "invalid course_id"})
		}

		newModule := &module.Module{
			ID:          uuid.New(),
			CourseID:    courseID,
			Title:       req.Title,
			Description: req.Description,
			Ordinal:     req.Ordinal,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		if err := repo.Create(c.Request().Context(), newModule); err != nil {
			return c.JSON(500, echo.Map{"error": "failed to create module"})
		}

		return c.JSON(201, echo.Map{"status": "created", "id": newModule.ID})
	}
}

func GetModulesByCourseIDHandler(repo module.Repository) echo.HandlerFunc {
	return func(c echo.Context) error {
		courseIDStr := c.Param("courseID")
		courseID, err := uuid.Parse(courseIDStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid course ID"})
		}

		modules, err := repo.GetByCourseID(c.Request().Context(), courseID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to load modules"})
		}

		return c.JSON(http.StatusOK, modules)
	}
}

func UpdateModuleHandler(repo module.Repository) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, _ := uuid.Parse(c.Param("id"))
		var req UpdateModuleRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(400, echo.Map{"error": "bad request"})
		}

		err := repo.Update(c.Request().Context(), &module.Module{
			ID:          id,
			Title:       req.Title,
			Description: req.Description,
			Ordinal:     req.Ordinal,
		})
		if err != nil {
			return c.JSON(500, echo.Map{"error": "update failed"})
		}
		return c.JSON(200, echo.Map{"status": "updated"})
	}
}

func DeleteModuleHandler(repo module.Repository) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, _ := uuid.Parse(c.Param("id"))
		if err := repo.SoftDelete(c.Request().Context(), id); err != nil {
			return c.JSON(500, echo.Map{"error": "delete failed"})
		}
		return c.JSON(200, echo.Map{"status": "deleted"})
	}
}
