package main

import (
	"os"

	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/handler"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/oapi"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/usecases/repository/repoimpl"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/usecases/services/ws"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/pkg/log"
	"go.uber.org/zap"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const baseURL = "/api/v1"

func main() {
	port := getEnv("APP_PORT", ":8080")

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	roomRepo := repoimpl.NewRoomRepository()
	hub := ws.NewHub(roomRepo)
	streamer := ws.NewStreamer(hub)
	h := handler.New(roomRepo, streamer)
	oapi.RegisterHandlersWithBaseURL(e, h, baseURL)

	log.L().Fatal("exit", zap.Error(e.Start(port)))
}

func getEnv(key string, def string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return def
	}

	return value
}
