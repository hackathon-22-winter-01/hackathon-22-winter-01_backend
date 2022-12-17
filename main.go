package main

import (
	"embed"
	"fmt"
	"net/http"

	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/handler"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/oapi"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/usecases/repository/repoimpl"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/usecases/services/ws"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/pkg/config"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/pkg/log"
	"go.uber.org/zap"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

//go:embed bin/frontend/dist
var dist embed.FS

const baseURL = "/api/v1"

func main() {
	config.ParseFlags()

	e := echo.New()

	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:       "bin/frontend/dist",
		Filesystem: http.FS(dist),
	}))
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	roomRepo := repoimpl.NewRoomRepository()
	hub := ws.NewHub(roomRepo)
	streamer := ws.NewStreamer(hub)
	h := handler.New(roomRepo, streamer)
	oapi.RegisterHandlersWithBaseURL(e, h, baseURL)

	if err := e.Start(fmt.Sprintf(":%d", *config.Port)); err != nil {
		log.L().Fatal("failed to start server", zap.Error(err))
	}
}
