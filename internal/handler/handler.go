package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/oapi"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/usecases/services/ws"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	stream ws.Streamer
}

func New(stream ws.Streamer) oapi.ServerInterface {
	return &Handler{stream}
}

func (h *Handler) Ping(c echo.Context) error {
	return c.String(200, "pong")
}

func (h *Handler) ConnectToWs(c echo.Context) error {
	uid := uuid.New()

	err := h.stream.ServeWS(c.Response().Writer, c.Request(), uid)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return nil
}
