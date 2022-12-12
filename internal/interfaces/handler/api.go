package handler

import (
	"net/http"

	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/interfaces/handler/oapi"
	"github.com/labstack/echo/v4"
)

type API struct{}

func NewAPI() oapi.ServerInterface {
	return &API{}
}

func (a *API) Ping(c echo.Context) error {
	return c.String(200, "pong")
}

func (a *API) GetWs(cardReset echo.Context) error {
	// TODO: 実装する
	return echo.NewHTTPError(http.StatusNotImplemented)
}

// deprecated
func (a *API) GetWsSchemas(c echo.Context) error {
	return echo.NewHTTPError(http.StatusNotImplemented)
}
