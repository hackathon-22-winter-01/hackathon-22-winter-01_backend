package main

import (
	"hackathon-22-winter-01/internal/interfaces/handler"
	"hackathon-22-winter-01/internal/interfaces/handler/oapi"
	"os"

	"github.com/labstack/echo/v4"
)

const baseURL = "/api/v1"

var port = getEnv("APP_PORT", ":8080")

func main() {
	e := echo.New()
	api := handler.NewAPI()

	h := oapi.NewStrictHandler(api, nil)
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
