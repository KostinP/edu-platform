package http

import (
	"github.com/kostinp/edu-platform-backend/api/http/handlers"
	"github.com/kostinp/edu-platform-backend/internal/answer"
	"github.com/kostinp/edu-platform-backend/internal/course"
	"github.com/kostinp/edu-platform-backend/internal/lesson"
	"github.com/kostinp/edu-platform-backend/internal/module"
	"github.com/kostinp/edu-platform-backend/internal/question"
	"github.com/kostinp/edu-platform-backend/internal/session"
	"github.com/kostinp/edu-platform-backend/internal/test"
	"github.com/kostinp/edu-platform-backend/internal/user"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo) {
	userRepo := user.NewPostgresRepo()
	e.POST("/auth/telegram", handlers.AuthTelegramHandler(userRepo))

	courseRepo := course.NewPostgresRepo()
	e.POST("/courses", handlers.CreateCourseHandler(courseRepo))
	e.GET("/courses", handlers.GetCoursesHandler(courseRepo))
	e.GET("/courses/:slug", handlers.GetCourseBySlugHandler(courseRepo))
	e.PUT("/courses/:id", handlers.UpdateCourseHandler(courseRepo))
	e.DELETE("/courses/:id", handlers.DeleteCourseHandler(courseRepo))

	moduleRepo := module.NewPostgresRepo()
	e.POST("/modules", handlers.CreateModuleHandler(moduleRepo))
	e.GET("/courses/:courseID/modules", handlers.GetModulesByCourseIDHandler(moduleRepo))
	e.PUT("/modules/:id", handlers.UpdateModuleHandler(moduleRepo))
	e.DELETE("/modules/:id", handlers.DeleteModuleHandler(moduleRepo))

	lessonRepo := lesson.NewPostgresRepo()
	e.POST("/lessons", handlers.CreateLessonHandler(lessonRepo))
	e.GET("/modules/:moduleID/lessons", handlers.GetLessonsByModuleIDHandler(lessonRepo))
	e.PUT("/lessons/:id", handlers.UpdateLessonHandler(lessonRepo))
	e.DELETE("/lessons/:id", handlers.DeleteLessonHandler(lessonRepo))

	testRepo := test.NewPostgresRepo()
	e.POST("/tests", handlers.CreateTestHandler(testRepo))
	e.PATCH("/tests/:id", handlers.UpdateTestHandler(testRepo))
	e.DELETE("/tests/:id", handlers.DeleteTestHandler(testRepo))

	questionRepo := question.NewPostgresRepo()
	e.POST("/tests/:id/questions", handlers.CreateQuestionHandler(questionRepo))
	e.PATCH("/questions/:id", handlers.UpdateQuestionHandler(questionRepo))
	e.DELETE("/questions/:id", handlers.DeleteQuestionHandler(questionRepo))

	sessionRepo := session.NewPostgresRepo()
	e.POST("/tests/:id/start", handlers.StartTestHandler(sessionRepo))
	e.POST("/tests/:id/submit", handlers.SubmitTestHandler(sessionRepo))
	answerRepo := answer.NewPostgresRepo()
	e.PATCH("/answers/:id/review", handlers.ReviewAnswerHandler(answerRepo))
	e.POST("/questions/:id/answer", handlers.SubmitAnswerHandler(answerRepo, questionRepo))
	e.GET("/sessions/:id/result", handlers.GetSessionResultHandler(sessionRepo, answerRepo, questionRepo))

	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Server is running!")
	})
}
