package ws

import (
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/oapi"
)

func (h *Hub) sendLifeChanged() (*oapi.WsResponse, error) {
	b := oapi.WsResponseBodyLifeChanged{}

	res := oapi.NewWsResponse(oapi.WsResponseTypeLifeChanged)
	if err := res.Body.FromWsResponseBodyLifeChanged(b); err != nil {
		return nil, err
	}

	return res, nil
}

func (h *Hub) sendCardUsed() (*oapi.WsResponse, error) {
	b := oapi.WsResponseBodyCardUsed{}

	res := oapi.NewWsResponse(oapi.WsResponseTypeCardUsed)
	if err := res.Body.FromWsResponseBodyCardUsed(b); err != nil {
		return nil, err
	}

	return res, nil
}

func (h *Hub) sendCardReset() (*oapi.WsResponse, error) {
	b := oapi.WsResponseBodyCardReset{}

	res := oapi.NewWsResponse(oapi.WsResponseTypeCardReset)
	if err := res.Body.FromWsResponseBodyCardReset(b); err != nil {
		return nil, err
	}

	return res, nil
}

func (h *Hub) sendRailCreated() (*oapi.WsResponse, error) {
	b := oapi.WsResponseBodyRailCreated{}

	res := oapi.NewWsResponse(oapi.WsResponseTypeRailCreated)
	if err := res.Body.FromWsResponseBodyRailCreated(b); err != nil {
		return nil, err
	}

	return res, nil
}

func (h *Hub) sendRailMerged() (*oapi.WsResponse, error) {
	b := oapi.WsResponseBodyRailMerged{}

	res := oapi.NewWsResponse(oapi.WsResponseTypeRailMerged)
	if err := res.Body.FromWsResponseBodyRailMerged(b); err != nil {
		return nil, err
	}

	return res, nil
}

func (h *Hub) sendBlockCreated() (*oapi.WsResponse, error) {
	b := oapi.WsResponseBodyBlockCreated{}

	res := oapi.NewWsResponse(oapi.WsResponseTypeBlockCreated)
	if err := res.Body.FromWsResponseBodyBlockCreated(b); err != nil {
		return nil, err
	}

	return res, nil
}
