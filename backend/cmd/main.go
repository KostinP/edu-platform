package main

import (
	"net/http"
	"os"

	"github.com/joho/godotenv"
	httpapi "github.com/kostinp/edu-platform-backend/api/http"
	"github.com/kostinp/edu-platform-backend/api/http/user"
	_ "github.com/kostinp/edu-platform-backend/docs"
	"github.com/kostinp/edu-platform-backend/pkg/db"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func main() {
	env := os.Getenv("ENV")
	if env == "" {
		env = "local"
	}

	// Загружаем .env
	if env == "local" {
		if err := godotenv.Load(".env.local"); err != nil {
			logrus.Warn("Failed to load .env.local")
		} else {
			logrus.Info("Loaded .env.local")
		}
	} else {
		if err := godotenv.Load(".env"); err != nil {
			logrus.Warn("Failed to load .env")
		} else {
			logrus.Info("Loaded .env")
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logrus.Infof("Starting server in %s mode", env)

	// Подключение к БД
	if err := db.Connect(); err != nil {
		logrus.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Close()

	// ✅ Автоматическое применение миграций
	// applyMigrations()

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())  // логирует HTTP запросы (метод, путь, статус, время)
	e.Use(middleware.Recover()) // чтобы сервер не падал на panic

	// CORS
	frontendOrigin := "http://localhost:3000"
	if env == "production" {
		frontendOrigin = "https://codesigned.ru"
	}
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{frontendOrigin},
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodOptions,
		},
	}))

	// Swagger
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Healthcheck
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Server is running!")
	})

	// Регистрация роутов
	user.RegisterRoutes(e)
	httpapi.RegisterRoutes(e)

	logrus.Infof("Server listening on :%s", port)
	e.Logger.Fatal(e.Start(":" + port))
}

func applyMigrations() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		logrus.Fatal("DATABASE_URL not set")
	}

	m, err := migrate.New(
		"file://migrations",
		databaseURL,
	)
	if err != nil {
		logrus.Fatalf("Failed to init migrate: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logrus.Fatalf("Failed to apply migrations: %v", err)
	}

	logrus.Info("Migrations applied successfully")
}
