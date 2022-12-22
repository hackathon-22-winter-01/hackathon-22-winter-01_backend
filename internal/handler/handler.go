package handler

import (
	"net/http"

	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/oapi"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/usecases/repository"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/usecases/services/ws"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	r      repository.RoomRepository
	stream ws.Streamer
}

func New(r repository.RoomRepository, stream ws.Streamer) oapi.ServerInterface {
	return &Handler{r, stream}
}

func (h *Handler) Ping(c echo.Context) error {
	return c.String(200, "pong")
}

func (h *Handler) ConnectToWs(c echo.Context, params oapi.ConnectToWsParams) error {
	err := h.stream.ServeWS(c.Response().Writer, c.Request(), params.PlayerId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return nil
}
