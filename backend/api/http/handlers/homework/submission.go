package handlers

import (
	"net/http"

	"github.com/kostinp/edu-platform-backend/internal/homework/submission"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type SubmitRequest struct {
	HomeworkID uuid.UUID `json:"homework_id"`
	UserID     uuid.UUID `json:"user_id"`
	Answer     string    `json:"answer"`
	FileURL    *string   `json:"file_url"`
}

type ReviewRequest struct {
	Review string  `json:"review"`
	Score  float64 `json:"score"`
}

type PeerReviewRequest struct {
	PeerID  uuid.UUID `json:"peer_id"`
	Comment string    `json:"comment"`
}

func PeerReviewHandler(repo submission.Repository) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			return c.JSON(400, echo.Map{"error": "invalid ID"})
		}

		var req PeerReviewRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(400, echo.Map{"error": "invalid input"})
		}

		if err := repo.AddPeerReview(c.Request().Context(), id, req.PeerID, req.Comment); err != nil {
			return c.JSON(500, echo.Map{"error": "peer review failed"})
		}

		return c.JSON(200, echo.Map{"status": "peer-reviewed"})
	}
}

func SubmitHomeworkHandler(repo submission.Repository) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req SubmitRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request"})
		}

		s := &submission.Submission{
			ID:         uuid.New(),
			HomeworkID: req.HomeworkID,
			UserID:     req.UserID,
			Answer:     req.Answer,
			FileURL:    req.FileURL,
		}

		if err := repo.Submit(c.Request().Context(), s); err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to submit"})
		}
		return c.JSON(http.StatusCreated, s)
	}
}

func ListUserSubmissionsHandler(repo submission.Repository) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID, err := uuid.Parse(c.Param("user_id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
		}

		submissions, err := repo.ListByUser(c.Request().Context(), userID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to list"})
		}
		return c.JSON(http.StatusOK, submissions)
	}
}

func ReviewHomeworkSubmissionHandler(repo submission.Repository) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			return c.JSON(400, echo.Map{"error": "invalid ID"})
		}

		var req ReviewRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(400, echo.Map{"error": "invalid input"})
		}

		if err := repo.Review(c.Request().Context(), id, req.Review, req.Score); err != nil {
			return c.JSON(500, echo.Map{"error": "failed to update submission"})
		}

		return c.JSON(200, echo.Map{"status": "reviewed"})
	}
}

func GetSubmissionHandler(repo submission.Repository) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			return c.JSON(400, echo.Map{"error": "invalid ID"})
		}

		sub, err := repo.GetByID(c.Request().Context(), id)
		if err != nil {
			return c.JSON(404, echo.Map{"error": "submission not found"})
		}

		return c.JSON(200, sub)
	}
}
