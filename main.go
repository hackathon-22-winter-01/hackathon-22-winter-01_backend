package main

import (
	"hackathon-22-winter-01/internal/interfaces/handler"
	"hackathon-22-winter-01/internal/interfaces/handler/oapi"
	"os"

	"github.com/labstack/echo/v4"
)

const baseURL = "/api/v1"

var port = getEnvOrPanic("APP_PORT")

func main() {
	e := echo.New()
	api := handler.NewAPI()

	h := oapi.NewStrictHandler(api, nil)
	oapi.RegisterHandlersWithBaseURL(e, h, baseURL)

	e.Logger.Fatal(e.Start(port))
}

func getEnvOrPanic(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		panic("environment variable " + key + " is not set")
	}

	return value
}
