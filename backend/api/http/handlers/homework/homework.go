package handlers

import (
	"net/http"
	"strconv"
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

func ListHomeworksForUserHandler(repo homework.Repository) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			return c.JSON(400, echo.Map{"error": "invalid user ID"})
		}

		filters := make(map[string]interface{})

		if t := c.QueryParam("type"); t != "" {
			filters["type"] = t
		}
		if r := c.QueryParam("required"); r != "" {
			val := r == "true"
			filters["required"] = &val
		}
		if cid := c.QueryParam("course_id"); cid != "" {
			if courseID, err := uuid.Parse(cid); err == nil {
				filters["course_id"] = courseID
			}
		}
		if s := c.QueryParam("status"); s != "" {
			filters["status"] = s
		}
		if df := c.QueryParam("due_from"); df != "" {
			if dueFrom, err := time.Parse(time.RFC3339, df); err == nil {
				filters["due_from"] = dueFrom
			}
		}
		if dt := c.QueryParam("due_to"); dt != "" {
			if dueTo, err := time.Parse(time.RFC3339, dt); err == nil {
				filters["due_to"] = dueTo
			}
		}
		if sb := c.QueryParam("sort_by"); sb != "" {
			filters["sort_by"] = sb
		}
		if so := c.QueryParam("sort_order"); so != "" {
			filters["sort_order"] = so
		}
		if l := c.QueryParam("limit"); l != "" {
			if limit, err := strconv.Atoi(l); err == nil {
				filters["limit"] = limit
			}
		}
		if o := c.QueryParam("offset"); o != "" {
			if offset, err := strconv.Atoi(o); err == nil {
				filters["offset"] = offset
			}
		}

		hws, err := repo.ListForUserFiltered(c.Request().Context(), userID, filters)
		if err != nil {
			return c.JSON(500, echo.Map{"error": "failed to fetch homeworks"})
		}
		return c.JSON(200, hws)
	}
}

func GetHomeworkStatsHandler(repo homework.Repository) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			return c.JSON(400, echo.Map{"error": "invalid user ID"})
		}
		stats, err := repo.GetStatsForUser(c.Request().Context(), userID)
		if err != nil {
			return c.JSON(500, echo.Map{"error": "failed to get stats"})
		}
		return c.JSON(200, stats)
	}
}

func ListHomeworksByAuthorHandler(repo homework.Repository) echo.HandlerFunc {
	return func(c echo.Context) error {
		authorID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			return c.JSON(400, echo.Map{"error": "invalid author ID"})
		}
		hws, err := repo.ListByAuthor(c.Request().Context(), authorID)
		if err != nil {
			return c.JSON(500, echo.Map{"error": "failed to list"})
		}
		return c.JSON(200, hws)
	}
}

func FilterHomeworksHandler(repo homework.Repository) echo.HandlerFunc {
	return func(c echo.Context) error {
		filters := map[string]interface{}{}

		// UUID фильтры
		for _, key := range []string{"course_id", "module_id", "lesson_id", "group_id", "user_id"} {
			if v := c.QueryParam(key); v != "" {
				if id, err := uuid.Parse(v); err == nil {
					filters[key] = id
				}
			}
		}

		// status=done|not_done
		status := c.QueryParam("status")
		if status == "done" || status == "not_done" {
			filters["status"] = status
		}

		// is_required
		if isReq := c.QueryParam("is_required"); isReq == "true" || isReq == "false" {
			filters["is_required"] = isReq == "true"
		}

		// due_before
		if due := c.QueryParam("due_before"); due != "" {
			if t, err := time.Parse(time.RFC3339, due); err == nil {
				filters["due_at <="] = t
			}
		}

		// Pagination
		limit := 20
		offset := 0
		if l := c.QueryParam("limit"); l != "" {
			if v, err := strconv.Atoi(l); err == nil {
				limit = v
			}
		}
		if o := c.QueryParam("offset"); o != "" {
			if v, err := strconv.Atoi(o); err == nil {
				offset = v
			}
		}

		// Требуется user_id для done/not_done
		if _, hasStatus := filters["status"]; hasStatus {
			if _, ok := filters["user_id"]; !ok {
				return c.JSON(http.StatusBadRequest, echo.Map{"error": "user_id required for status filter"})
			}
		}

		homeworks, err := repo.Filter(c.Request().Context(), filters, limit, offset)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to filter"})
		}
		return c.JSON(http.StatusOK, homeworks)
	}
}
