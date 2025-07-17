package tag

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/kostinp/edu-platform-backend/internal/tag"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	repo tag.Repository
}

func NewHandler(repo tag.Repository) *Handler {
	return &Handler{repo: repo}
}

// CreateTag создаёт новый тег
func (h *Handler) CreateTag(c echo.Context) error {
	var t tag.Tag
	if err := c.Bind(&t); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request body"})
	}
	if t.AuthorID == uuid.Nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "author_id is required"})
	}

	if err := h.repo.Create(c.Request().Context(), &t); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to create tag"})
	}

	return c.JSON(http.StatusCreated, t)
}

// GetAllTags возвращает все теги
func (h *Handler) GetAllTags(c echo.Context) error {
	tags, err := h.repo.GetAll(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to get tags"})
	}
	return c.JSON(http.StatusOK, tags)
}

// UpdateTag обновляет тег по ID
func (h *Handler) UpdateTag(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid tag id"})
	}

	var dto tag.UpdateTagDTO
	if err := c.Bind(&dto); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request body"})
	}

	if err := h.repo.Update(c.Request().Context(), id, dto); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to update tag"})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "tag updated"})
}

// DeleteTag выполняет мягкое удаление тега
func (h *Handler) DeleteTag(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid tag id"})
	}

	if err := h.repo.Delete(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to delete tag"})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "tag deleted"})
}

// --- Методы для связи с курсами ---

func (h *Handler) AddTagToCourse(c echo.Context) error {
	courseIDStr := c.Param("course_id")
	tagIDStr := c.Param("tag_id")

	courseID, err := uuid.Parse(courseIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid course id"})
	}
	tagID, err := uuid.Parse(tagIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid tag id"})
	}

	err = h.repo.AddTagToCourse(c.Request().Context(), courseID, tagID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to add tag to course"})
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "tag added to course"})
}

func (h *Handler) RemoveTagFromCourse(c echo.Context) error {
	courseIDStr := c.Param("course_id")
	tagIDStr := c.Param("tag_id")

	courseID, err := uuid.Parse(courseIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid course id"})
	}
	tagID, err := uuid.Parse(tagIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid tag id"})
	}

	err = h.repo.RemoveTagFromCourse(c.Request().Context(), courseID, tagID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to remove tag from course"})
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "tag removed from course"})
}

func (h *Handler) ListTagsByCourse(c echo.Context) error {
	courseIDStr := c.Param("course_id")
	courseID, err := uuid.Parse(courseIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid course id"})
	}

	tags, err := h.repo.ListTagsByCourse(c.Request().Context(), courseID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to list tags for course"})
	}
	return c.JSON(http.StatusOK, tags)
}

// --- Методы для уроков ---

func (h *Handler) AddTagToLesson(c echo.Context) error {
	lessonIDStr := c.Param("lesson_id")
	tagIDStr := c.Param("tag_id")

	lessonID, err := uuid.Parse(lessonIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid lesson id"})
	}
	tagID, err := uuid.Parse(tagIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid tag id"})
	}

	err = h.repo.AddTagToLesson(c.Request().Context(), lessonID, tagID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to add tag to lesson"})
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "tag added to lesson"})
}

func (h *Handler) RemoveTagFromLesson(c echo.Context) error {
	lessonIDStr := c.Param("lesson_id")
	tagIDStr := c.Param("tag_id")

	lessonID, err := uuid.Parse(lessonIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid lesson id"})
	}
	tagID, err := uuid.Parse(tagIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid tag id"})
	}

	err = h.repo.RemoveTagFromLesson(c.Request().Context(), lessonID, tagID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to remove tag from lesson"})
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "tag removed from lesson"})
}

func (h *Handler) ListTagsByLesson(c echo.Context) error {
	lessonIDStr := c.Param("lesson_id")
	lessonID, err := uuid.Parse(lessonIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid lesson id"})
	}

	tags, err := h.repo.ListTagsByLesson(c.Request().Context(), lessonID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to list tags for lesson"})
	}
	return c.JSON(http.StatusOK, tags)
}

// --- Методы для тестов ---

func (h *Handler) AddTagToTest(c echo.Context) error {
	testIDStr := c.Param("test_id")
	tagIDStr := c.Param("tag_id")

	testID, err := uuid.Parse(testIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid test id"})
	}
	tagID, err := uuid.Parse(tagIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid tag id"})
	}

	err = h.repo.AddTagToTest(c.Request().Context(), testID, tagID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to add tag to test"})
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "tag added to test"})
}

func (h *Handler) RemoveTagFromTest(c echo.Context) error {
	testIDStr := c.Param("test_id")
	tagIDStr := c.Param("tag_id")

	testID, err := uuid.Parse(testIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid test id"})
	}
	tagID, err := uuid.Parse(tagIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid tag id"})
	}

	err = h.repo.RemoveTagFromTest(c.Request().Context(), testID, tagID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to remove tag from test"})
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "tag removed from test"})
}

func (h *Handler) ListTagsByTest(c echo.Context) error {
	testIDStr := c.Param("test_id")
	testID, err := uuid.Parse(testIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid test id"})
	}

	tags, err := h.repo.ListTagsByTest(c.Request().Context(), testID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to list tags for test"})
	}
	return c.JSON(http.StatusOK, tags)
}
