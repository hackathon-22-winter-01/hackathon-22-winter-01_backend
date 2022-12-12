package handler

import (
	"net/http"

	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/handler/oapi"
	"github.com/labstack/echo/v4"
)

type Handler struct{}

func New() oapi.ServerInterface {
	return &Handler{}
}

func (h *Handler) Ping(c echo.Context) error {
	return c.String(200, "pong")
}

func (h *Handler) ConnectToWs(cardReset echo.Context) error {
	// TODO: (usecases/service/ws).ServeWsを読んで、websocketの接続を確立する
	return echo.NewHTTPError(http.StatusNotImplemented)
}

// deprecated
func (h *Handler) UseWsSchemas(c echo.Context) error {
	return echo.NewHTTPError(http.StatusNotImplemented)
}
