package http

import (
	authHandlers "github.com/kostinp/edu-platform-backend/api/http/handlers/auth"
	"github.com/kostinp/edu-platform-backend/internal/auth/user"

	courseHandlers "github.com/kostinp/edu-platform-backend/api/http/handlers/courses"
	"github.com/kostinp/edu-platform-backend/internal/course/course"
	"github.com/kostinp/edu-platform-backend/internal/course/lesson"
	"github.com/kostinp/edu-platform-backend/internal/course/module"

	testingHandlers "github.com/kostinp/edu-platform-backend/api/http/handlers/testing"
	"github.com/kostinp/edu-platform-backend/internal/testing/answer"
	"github.com/kostinp/edu-platform-backend/internal/testing/question"
	"github.com/kostinp/edu-platform-backend/internal/testing/session"
	"github.com/kostinp/edu-platform-backend/internal/testing/test"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo) {
	userRepo := user.NewPostgresRepo()
	e.POST("/auth/telegram", authHandlers.AuthTelegramHandler(userRepo))

	courseRepo := course.NewPostgresRepo()
	e.POST("/courses", courseHandlers.CreateCourseHandler(courseRepo))
	e.GET("/courses", courseHandlers.GetCoursesHandler(courseRepo))
	e.GET("/courses/:slug", courseHandlers.GetCourseBySlugHandler(courseRepo))
	e.PUT("/courses/:id", courseHandlers.UpdateCourseHandler(courseRepo))
	e.DELETE("/courses/:id", courseHandlers.DeleteCourseHandler(courseRepo))

	moduleRepo := module.NewPostgresRepo()
	e.POST("/modules", courseHandlers.CreateModuleHandler(moduleRepo))
	e.GET("/courses/:courseID/modules", courseHandlers.GetModulesByCourseIDHandler(moduleRepo))
	e.PUT("/modules/:id", courseHandlers.UpdateModuleHandler(moduleRepo))
	e.DELETE("/modules/:id", courseHandlers.DeleteModuleHandler(moduleRepo))

	lessonRepo := lesson.NewPostgresRepo()
	e.POST("/lessons", courseHandlers.CreateLessonHandler(lessonRepo))
	e.GET("/modules/:moduleID/lessons", courseHandlers.GetLessonsByModuleIDHandler(lessonRepo))
	e.PUT("/lessons/:id", courseHandlers.UpdateLessonHandler(lessonRepo))
	e.DELETE("/lessons/:id", courseHandlers.DeleteLessonHandler(lessonRepo))

	testRepo := test.NewPostgresRepo()
	e.POST("/tests", testingHandlers.CreateTestHandler(testRepo))
	e.PATCH("/tests/:id", testingHandlers.UpdateTestHandler(testRepo))
	e.DELETE("/tests/:id", testingHandlers.DeleteTestHandler(testRepo))

	questionRepo := question.NewPostgresRepo()
	e.POST("/tests/:id/questions", testingHandlers.CreateQuestionHandler(questionRepo))
	e.PATCH("/questions/:id", testingHandlers.UpdateQuestionHandler(questionRepo))
	e.DELETE("/questions/:id", testingHandlers.DeleteQuestionHandler(questionRepo))

	sessionRepo := session.NewPostgresRepo()
	e.POST("/tests/:id/start", testingHandlers.StartTestHandler(sessionRepo))
	e.POST("/tests/:id/submit", testingHandlers.SubmitTestHandler(sessionRepo))

	answerRepo := answer.NewPostgresRepo()
	e.PATCH("/answers/:id/review", testingHandlers.ReviewAnswerHandler(answerRepo))
	e.POST("/questions/:id/answer", testingHandlers.SubmitAnswerHandler(answerRepo, questionRepo))
	e.GET("/sessions/:id/result", testingHandlers.GetSessionResultHandler(sessionRepo, answerRepo, questionRepo))

	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Server is running!")
	})
}
