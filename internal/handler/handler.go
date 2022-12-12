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
	return &Handler{}
}

func (h *Handler) Ping(c echo.Context) error {
	return c.String(200, "pong")
}

func (h *Handler) ConnectToWs(c echo.Context) error {
	// TODO: (usecases/service/ws).ServeWsを読んで、websocketの接続を確立する

	uid := uuid.New()
	err := h.stream.ServeWS(c.Response().Writer, c.Request(), uid)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return echo.NewHTTPError(http.StatusNoContent)
}

// deprecated
func (h *Handler) UseWsSchemas(c echo.Context) error {
	return echo.NewHTTPError(http.StatusNotImplemented)
}
