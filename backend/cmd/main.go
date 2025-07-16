package main

import (
	"github.com/kostinp/edu-platform-backend/api/http"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"

	"github.com/kostinp/edu-platform-backend/pkg/db"
)

func main() {
	_ = godotenv.Load()

	if err := db.Connect(); err != nil {
		logrus.Fatalf("Failed to connect to DB: %v", err)
	}

	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Server is running!")
	})

	logrus.Info("Server started on http://localhost:8080")

	http.RegisterRoutes(e)

	e.Logger.Fatal(e.Start(":8080"))
}
