package main

import (
	"net/http"
	"os"

	"github.com/kostinp/edu-platform-backend/api/http/user"
	_ "github.com/kostinp/edu-platform-backend/docs"
	"github.com/kostinp/edu-platform-backend/pkg/db"

	httpapi "github.com/kostinp/edu-platform-backend/api/http"
	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

func main() {
	env := os.Getenv("ENV")
	if env == "" {
		env = "local" // по умолчанию локальное окружение
	}

	// Загружаем переменные из соответствующего .env файла
	if env == "local" {
		if err := godotenv.Load(".env.local"); err != nil {
			logrus.Warn("Failed to load .env.local file, make sure it exists")
		} else {
			logrus.Info("Loaded .env.local file")
		}
	} else {
		if err := godotenv.Load(".env"); err != nil {
			logrus.Warn("Failed to load .env file, make sure it exists")
		} else {
			logrus.Info("Loaded .env file")
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logrus.Infof("Starting server in %s mode", env)

	if err := db.Connect(); err != nil {
		logrus.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Close()

	e := echo.New()

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

	user.RegisterRoutes(e)
	httpapi.RegisterRoutes(e)

	logrus.Infof("Server listening on :%s", port)
	e.Logger.Fatal(e.Start(":" + port))
}
