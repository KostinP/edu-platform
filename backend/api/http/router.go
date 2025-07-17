package http

import (
	authHandlers "github.com/kostinp/edu-platform-backend/api/http/handlers/auth"
	"github.com/kostinp/edu-platform-backend/internal/auth/user"
	"github.com/kostinp/edu-platform-backend/internal/progress"

	courseHandlers "github.com/kostinp/edu-platform-backend/api/http/handlers/courses"
	"github.com/kostinp/edu-platform-backend/internal/course/course"
	"github.com/kostinp/edu-platform-backend/internal/course/lesson"
	"github.com/kostinp/edu-platform-backend/internal/course/module"

	testingHandlers "github.com/kostinp/edu-platform-backend/api/http/handlers/testing"
	"github.com/kostinp/edu-platform-backend/internal/testing/answer"
	"github.com/kostinp/edu-platform-backend/internal/testing/question"
	"github.com/kostinp/edu-platform-backend/internal/testing/session"
	"github.com/kostinp/edu-platform-backend/internal/testing/test"

	homeworkHandlers "github.com/kostinp/edu-platform-backend/api/http/handlers/homework"
	"github.com/kostinp/edu-platform-backend/internal/homework/homework"
	"github.com/kostinp/edu-platform-backend/internal/homework/submission"

	groupsHandlers "github.com/kostinp/edu-platform-backend/api/http/handlers/groups"
	"github.com/kostinp/edu-platform-backend/internal/group/group"

	progressHandler "github.com/kostinp/edu-platform-backend/api/http/handlers/progress"

	"github.com/labstack/echo/v4"
)

// RegisterRoutes регистрирует все маршруты приложения
func RegisterRoutes(e *echo.Echo) {
	registerAuthRoutes(e)
	registerCourseRoutes(e)
	registerTestingRoutes(e)
	registerHomeworkRoutes(e)
	registerHealthCheckRoute(e)
}

func registerAuthRoutes(e *echo.Echo) {
	userRepo := user.NewPostgresRepo()
	e.POST("/auth/telegram", authHandlers.AuthTelegramHandler(userRepo))
}

func registerCourseRoutes(e *echo.Echo) {
	courseRepo := course.NewPostgresRepo()
	moduleRepo := module.NewPostgresRepo()
	lessonRepo := lesson.NewPostgresRepo()

	// Course routes
	e.POST("/courses", courseHandlers.CreateCourseHandler(courseRepo))
	e.GET("/courses", courseHandlers.GetCoursesHandler(courseRepo))
	e.GET("/courses/:slug", courseHandlers.GetCourseBySlugHandler(courseRepo))
	e.PUT("/courses/:id", courseHandlers.UpdateCourseHandler(courseRepo))
	e.DELETE("/courses/:id", courseHandlers.DeleteCourseHandler(courseRepo))

	// Module routes
	e.POST("/modules", courseHandlers.CreateModuleHandler(moduleRepo))
	e.GET("/courses/:courseID/modules", courseHandlers.GetModulesByCourseIDHandler(moduleRepo))
	e.PUT("/modules/:id", courseHandlers.UpdateModuleHandler(moduleRepo))
	e.DELETE("/modules/:id", courseHandlers.DeleteModuleHandler(moduleRepo))

	// Lesson routes
	e.POST("/lessons", courseHandlers.CreateLessonHandler(lessonRepo))
	e.GET("/modules/:moduleID/lessons", courseHandlers.GetLessonsByModuleIDHandler(lessonRepo))
	e.PUT("/lessons/:id", courseHandlers.UpdateLessonHandler(lessonRepo))
	e.DELETE("/lessons/:id", courseHandlers.DeleteLessonHandler(lessonRepo))
}

func registerTestingRoutes(e *echo.Echo) {
	testRepo := test.NewPostgresRepo()
	questionRepo := question.NewPostgresRepo()
	sessionRepo := session.NewPostgresRepo()
	answerRepo := answer.NewPostgresRepo()

	// Test routes
	e.POST("/tests", testingHandlers.CreateTestHandler(testRepo))
	e.PATCH("/tests/:id", testingHandlers.UpdateTestHandler(testRepo))
	e.DELETE("/tests/:id", testingHandlers.DeleteTestHandler(testRepo))

	// Question routes
	e.POST("/tests/:id/questions", testingHandlers.CreateQuestionHandler(questionRepo))
	e.PATCH("/questions/:id", testingHandlers.UpdateQuestionHandler(questionRepo))
	e.DELETE("/questions/:id", testingHandlers.DeleteQuestionHandler(questionRepo))

	// Session routes
	e.POST("/tests/:id/start", testingHandlers.StartTestHandler(sessionRepo))
	e.POST("/tests/:id/submit", testingHandlers.SubmitTestHandler(sessionRepo))

	// Answer and analytics routes
	e.PATCH("/answers/:id/review", testingHandlers.ReviewAnswerHandler(answerRepo))
	e.POST("/questions/:id/answer", testingHandlers.SubmitAnswerHandler(answerRepo, questionRepo))
	e.GET("/sessions/:id/result", testingHandlers.GetSessionResultHandler(sessionRepo, answerRepo, questionRepo))
}

func registerHomeworkRoutes(e *echo.Echo) {
	hwRepoInstance := homework.NewPostgresRepo()
	sbRepoInstance := submission.NewPostgresRepo()

	// Homework routes
	e.POST("/homeworks", homeworkHandlers.CreateHomeworkHandler(hwRepoInstance))
	// e.GET("/homeworks", homeworkHandlers.ListHomeworksHandler(hwRepoInstance))
	e.GET("/homeworks", homeworkHandlers.FilterHomeworksHandler(hwRepoInstance))
	e.GET("/users/:id/homeworks", homeworkHandlers.ListHomeworksForUserHandler(hwRepoInstance))
	e.GET("/users/:id/homeworks/stats", homeworkHandlers.GetHomeworkStatsHandler(hwRepoInstance))
	e.GET("/authors/:id/homeworks", homeworkHandlers.ListHomeworksByAuthorHandler(hwRepoInstance))

	// Submission routes
	e.POST("/submissions", homeworkHandlers.SubmitHomeworkHandler(sbRepoInstance))
	e.GET("/users/:user_id/submissions", homeworkHandlers.ListUserSubmissionsHandler(sbRepoInstance))
	e.PUT("/submissions/:id/review", homeworkHandlers.ReviewHomeworkSubmissionHandler(sbRepoInstance))
	e.GET("/submissions/:id", homeworkHandlers.GetSubmissionHandler(sbRepoInstance))
	e.POST("/submissions/:id/peer-review", homeworkHandlers.PeerReviewHandler(sbRepoInstance))

}

func registerGroupRoutes(e *echo.Echo) {
	grRepoInstance := group.NewRepository()

	// Group routes
	e.POST("/groups", groupsHandlers.CreateGroupHandler(grRepoInstance))
	e.POST("/groups/:id/members", groupsHandlers.AddGroupMemberHandler(grRepoInstance))
	e.GET("/users/:id/groups", groupsHandlers.ListGroupsByUserHandler(grRepoInstance))
}

func registerProgressRoutes(e *echo.Echo) {
	// Импортируем репозитории из внутренних пакетов
	moduleRepo := module.NewPostgresRepo()
	lessonRepo := lesson.NewPostgresRepo()
	testRepo := test.NewPostgresRepo()
	sessionRepo := session.NewPostgresRepo()
	homeworkRepo := homework.NewPostgresRepo()
	submissionRepo := submission.NewPostgresRepo()

	// Создаем ProgressRepo с passing score, например 100%
	passScore := 100.0
	progressRepo := progress.NewProgressRepo(
		moduleRepo,
		lessonRepo,
		testRepo,
		sessionRepo,
		homeworkRepo,
		submissionRepo,
		passScore,
	)

	progressH := progressHandler.NewProgressHandler(progressRepo)

	// Регистрируем роутинг с методом
	e.GET("/courses/:course_id/progress", progressH.GetCourseProgress)
}

func registerHealthCheckRoute(e *echo.Echo) {
	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Server is running!")
	})
}
