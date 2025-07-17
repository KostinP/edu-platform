package progress

import (
	"context"

	"github.com/google/uuid"

	"github.com/kostinp/edu-platform-backend/internal/course/lesson"
	"github.com/kostinp/edu-platform-backend/internal/course/module"
	"github.com/kostinp/edu-platform-backend/internal/homework/homework"
	"github.com/kostinp/edu-platform-backend/internal/homework/submission"
	"github.com/kostinp/edu-platform-backend/internal/testing/session"
	"github.com/kostinp/edu-platform-backend/internal/testing/test"
)

type ProgressRepo struct {
	moduleRepo     module.Repository
	lessonRepo     lesson.Repository
	testRepo       test.Repository
	sessionRepo    session.Repository
	homeworkRepo   homework.Repository
	submissionRepo submission.Repository
	passScore      float64
}

func NewProgressRepo(
	mr module.Repository,
	lr lesson.Repository,
	tr test.Repository,
	sr session.Repository,
	hr homework.Repository,
	sbr submission.Repository,
	passScore float64,
) *ProgressRepo {
	return &ProgressRepo{
		moduleRepo:     mr,
		lessonRepo:     lr,
		testRepo:       tr,
		sessionRepo:    sr,
		homeworkRepo:   hr,
		submissionRepo: sbr,
		passScore:      passScore,
	}
}

// GetCourseProgress возвращает прогресс пользователя по курсу от 0.0 до 1.0
func (r *ProgressRepo) GetCourseProgress(ctx context.Context, userID, courseID uuid.UUID) (float64, error) {
	modules, err := r.moduleRepo.GetByCourseID(ctx, courseID)
	if err != nil {
		return 0, err
	}

	var totalLessons int
	var doneLessons int

	for _, m := range modules {
		lessons, err := r.lessonRepo.GetByModuleID(ctx, m.ID)
		if err != nil {
			return 0, err
		}

		totalLessons += len(lessons)

		for _, l := range lessons {
			done, err := r.isLessonDone(ctx, userID, l.ID)
			if err != nil {
				return 0, err
			}
			if done {
				doneLessons++
			}
		}
	}

	if totalLessons == 0 {
		return 0, nil
	}

	return float64(doneLessons) / float64(totalLessons), nil
}

// isLessonDone проверяет, выполнен ли урок пользователем
func (r *ProgressRepo) isLessonDone(ctx context.Context, userID, lessonID uuid.UUID) (bool, error) {
	// Проверка тестов
	tests, err := r.testRepo.GetByLessonID(ctx, lessonID)
	if err != nil {
		return false, err
	}
	for _, t := range tests {
		session, err := r.sessionRepo.GetLastFinishedByUserAndTest(ctx, userID, t.ID)
		if err == nil && session != nil && session.Score != nil && *session.Score >= r.passScore {
			return true, nil
		}
	}

	// Проверка домашних заданий
	homeworks, err := r.homeworkRepo.ListByLesson(ctx, lessonID)
	if err != nil {
		return false, err
	}
	for _, hw := range homeworks {
		submission, err := r.submissionRepo.GetByUserAndHomework(ctx, userID, hw.ID)
		if err == nil && submission != nil && (submission.Status == "submitted" || submission.Status == "reviewed") {
			return true, nil
		}
	}

	return false, nil
}
