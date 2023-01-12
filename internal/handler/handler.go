package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/oapi"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/usecases/repository"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/usecases/services/ws"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/pkg/optional"
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
	opts := ws.ServeWsOpts{
		PlayerID:   uuid.New(),
		PlayerName: params.Name,
		RoomID:     optional.NewFromPtr(params.RoomId),
	}

	err := h.stream.ServeWS(c.Response().Writer, c.Request(), opts)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return nil
}
