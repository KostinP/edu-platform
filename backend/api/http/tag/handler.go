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

// Create создаёт новый тег.
// @Summary Создать тег
// @Description Создаёт новый тег с именем и автором
// @Tags tags
// @Accept json
// @Produce json
// @Param input body CreateTagRequest true "Тело запроса создания тега"
// @Success 201 {object} TagResponse
// @Failure 400 {object} echo.Map "Неверный формат запроса или отсутствует author_id"
// @Failure 500 {object} echo.Map "Ошибка сервера при создании тега"
// @Router /tags [post]
func (h *Handler) Create(c echo.Context) error {
	var req CreateTagRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request body"})
	}
	if req.AuthorID == uuid.Nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "author_id is required"})
	}

	t := &tag.Tag{
		Name:     req.Name,
		AuthorID: req.AuthorID,
	}
	if err := h.repo.Create(c.Request().Context(), t); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to create tag"})
	}
	return c.JSON(http.StatusCreated, ToTagResponse(t))
}

// GetAll возвращает список всех тегов.
// @Summary Получить все теги
// @Description Возвращает список всех тегов, которые не удалены
// @Tags tags
// @Produce json
// @Success 200 {array} TagResponse
// @Failure 500 {object} echo.Map "Ошибка сервера при получении тегов"
// @Router /tags [get]
func (h *Handler) GetAll(c echo.Context) error {
	tags, err := h.repo.GetAll(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to get tags"})
	}
	return c.JSON(http.StatusOK, ToTagResponseList(tags))
}

// Update обновляет тег по ID.
// @Summary Обновить тег
// @Description Обновляет имя тега по его ID
// @Tags tags
// @Accept json
// @Produce json
// @Param id path string true "ID тега"
// @Param input body UpdateTagRequest true "Данные для обновления тега"
// @Success 200 {object} echo.Map "Тег успешно обновлён"
// @Failure 400 {object} echo.Map "Неверный ID или тело запроса"
// @Failure 500 {object} echo.Map "Ошибка сервера при обновлении тега"
// @Router /tags/{id} [put]
func (h *Handler) Update(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid tag id"})
	}

	var req UpdateTagRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request body"})
	}

	dto := tag.UpdateTagDTO{Name: req.Name}
	if err := h.repo.Update(c.Request().Context(), id, dto); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to update tag"})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "tag updated"})
}

// Delete выполняет мягкое удаление тега по ID.
// @Summary Удалить тег
// @Description Выполняет мягкое удаление тега по ID (устанавливает deleted_at)
// @Tags tags
// @Produce json
// @Param id path string true "ID тега"
// @Success 200 {object} echo.Map "Тег успешно удалён"
// @Failure 400 {object} echo.Map "Неверный ID"
// @Failure 500 {object} echo.Map "Ошибка сервера при удалении тега"
// @Router /tags/{id} [delete]
func (h *Handler) Delete(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid tag id"})
	}

	if err := h.repo.Delete(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to delete tag"})
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "tag deleted"})
}
