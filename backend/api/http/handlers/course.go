package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/kostinp/edu-platform-backend/internal/course"

	"github.com/labstack/echo/v4"
)

type CreateCourseRequest struct {
	Slug        string `json:"slug"`
	Title       string `json:"title"`
	Description string `json:"description"`
	AuthorID    string `json:"author_id"`
}

type UpdateCourseRequest struct {
	Slug        string `json:"slug"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func CreateCourseHandler(repo course.Repository) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req CreateCourseRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(400, echo.Map{"error": "invalid request"})
		}

		authorID, err := uuid.Parse(req.AuthorID)
		if err != nil {
			return c.JSON(400, echo.Map{"error": "invalid author_id"})
		}

		course := &course.Course{
			ID:          uuid.New(),
			Slug:        req.Slug,
			Title:       req.Title,
			Description: req.Description,
			AuthorID:    authorID,
			CreatedAt:   time.Now(),
		}

		if err := repo.Create(c.Request().Context(), course); err != nil {
			return c.JSON(500, echo.Map{"error": "failed to create course"})
		}

		return c.JSON(201, echo.Map{"status": "created", "id": course.ID})
	}
}

func GetCoursesHandler(repo course.Repository) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := context.Background()
		courses, err := repo.GetAll(ctx)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to fetch courses"})
		}
		return c.JSON(http.StatusOK, courses)
	}
}

func GetCourseBySlugHandler(repo course.Repository) echo.HandlerFunc {
	return func(c echo.Context) error {
		slug := c.Param("slug")
		course, err := repo.GetBySlug(c.Request().Context(), slug)
		if err != nil {
			return c.JSON(404, echo.Map{"error": "not found"})
		}
		return c.JSON(200, course)
	}
}

func UpdateCourseHandler(repo course.Repository) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, _ := uuid.Parse(c.Param("id"))
		var req UpdateCourseRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(400, echo.Map{"error": "bad request"})
		}

		err := repo.Update(c.Request().Context(), &course.Course{
			ID:          id,
			Slug:        req.Slug,
			Title:       req.Title,
			Description: req.Description,
		})
		if err != nil {
			return c.JSON(500, echo.Map{"error": "update failed"})
		}

		return c.JSON(200, echo.Map{"status": "updated"})
	}
}

func DeleteCourseHandler(repo course.Repository) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, _ := uuid.Parse(c.Param("id"))
		if err := repo.SoftDelete(c.Request().Context(), id); err != nil {
			return c.JSON(500, echo.Map{"error": "delete failed"})
		}
		return c.JSON(200, echo.Map{"status": "deleted"})
	}
}
