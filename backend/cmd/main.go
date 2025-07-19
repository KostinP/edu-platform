package main

import (
	nethttp "net/http"

	httpapi "github.com/kostinp/edu-platform-backend/api/http"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"

	"github.com/kostinp/edu-platform-backend/pkg/db"
)

func main() {
	_ = godotenv.Load()

	if err := db.Connect(); err != nil {
		logrus.Fatalf("Failed to connect to DB: %v", err)
	}

	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{nethttp.MethodGet, nethttp.MethodPost, nethttp.MethodPut, nethttp.MethodDelete, nethttp.MethodOptions},
	}))

	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Server is running!")
	})

	logrus.Info("Server started on http://localhost:8080")

	httpapi.RegisterRoutes(e)

	e.Logger.Fatal(e.Start(":8080"))
}
