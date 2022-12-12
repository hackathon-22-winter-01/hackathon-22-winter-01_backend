package main

import (
	"os"

	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/handler"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/handler/oapi"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const baseURL = "/api/v1"

var port = getEnv("APP_PORT", ":8080")

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	h := handler.New()
	oapi.RegisterHandlersWithBaseURL(e, h, baseURL)

	e.Logger.Fatal(e.Start(port))
}

func getEnv(key string, def string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return def
	}

	return value
}
