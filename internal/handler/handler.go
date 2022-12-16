package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/domain"
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

func (h *Handler) ConnectToWs(c echo.Context) error {
	uid := uuid.New()

	err := h.stream.ServeWS(c.Response().Writer, c.Request(), uid)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return nil
}

func (h *Handler) JoinRoom(c echo.Context) error {
	req := new(oapi.JoinRoomRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	userID := uuid.New()
	player := domain.NewPlayer(userID, req.PlayerName)

	err := h.r.JoinRoom(req.RoomId, player)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	room, err := h.r.FindRoom(req.RoomId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, room)
}

func (h *Handler) CreateRoom(c echo.Context) error {
	req := new(oapi.CreateRoomRequest)

	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	userID := uuid.New()
	player := domain.NewPlayer(userID, req.PlayerName)

	room, err := h.r.CreateRoom(player)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, room)
}

func (h *Handler) GetRoom(c echo.Context, roomID oapi.RoomId) error {
	room, err := h.r.FindRoom(roomID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, room)
}
