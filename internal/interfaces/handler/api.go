package handler

import (
	"context"

	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/interfaces/handler/oapi"
)

type API struct{}

func NewAPI() oapi.StrictServerInterface {
	return &API{}
}

func (a *API) Ping(ctx context.Context, request oapi.PingRequestObject) (oapi.PingResponseObject, error) {
	return oapi.Ping200JSONResponse{
		Message: "pong",
	}, nil
}

func (a *API) GetWs(ctx context.Context, request oapi.GetWsRequestObject) (oapi.GetWsResponseObject, error) {
	panic("implement me")
}

// deprecated
func (a *API) GetWsSchemas(ctx context.Context, request oapi.GetWsSchemasRequestObject) (oapi.GetWsSchemasResponseObject, error) {
	return oapi.GetWsSchemas200JSONResponse{}, nil
}
