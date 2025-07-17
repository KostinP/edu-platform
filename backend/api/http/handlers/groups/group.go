package handlers

import (
	"net/http"
	"time"

	"github.com/kostinp/edu-platform-backend/internal/group/group"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type CreateGroupRequest struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	OwnerID     uuid.UUID `json:"owner_id"`
}

func CreateGroupHandler(repo group.Repository) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req CreateGroupRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid input"})
		}

		g := &group.Group{
			ID:          uuid.New(),
			Name:        req.Name,
			Description: req.Description,
			OwnerID:     req.OwnerID,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		if err := repo.Create(c.Request().Context(), g); err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to create"})
		}
		return c.JSON(http.StatusCreated, g)
	}
}

type AddMemberRequest struct {
	UserID uuid.UUID `json:"user_id"`
	Role   string    `json:"role"` // default: "student"
}

func AddGroupMemberHandler(repo group.Repository) echo.HandlerFunc {
	return func(c echo.Context) error {
		groupID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			return c.JSON(400, echo.Map{"error": "invalid group ID"})
		}

		var req AddMemberRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(400, echo.Map{"error": "invalid body"})
		}

		if req.Role == "" {
			req.Role = "student"
		}

		err = repo.AddMember(c.Request().Context(), groupID, req.UserID, req.Role)
		if err != nil {
			return c.JSON(500, echo.Map{"error": "failed to add member"})
		}
		return c.JSON(200, echo.Map{"status": "added"})
	}
}

func ListGroupsByUserHandler(repo group.Repository) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			return c.JSON(400, echo.Map{"error": "invalid user ID"})
		}

		groups, err := repo.ListByUser(c.Request().Context(), userID)
		if err != nil {
			return c.JSON(500, echo.Map{"error": "failed to fetch groups"})
		}
		return c.JSON(200, groups)
	}
}
