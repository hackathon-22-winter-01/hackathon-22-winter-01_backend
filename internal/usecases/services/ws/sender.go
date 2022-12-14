package ws

import (
	"github.com/google/uuid"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/oapi"
)

func (h *Hub) sendConnected(playerID uuid.UUID) (*oapi.WsResponse, error) {
	b := oapi.WsResponseBodyConnected{
		PlayerId: playerID,
	}

	res := oapi.WsResponseFromType(oapi.WsResponseTypeConnected)
	if err := res.Body.FromWsResponseBodyConnected(b); err != nil {
		return nil, err
	}

	return res, nil
}

func (h *Hub) sendLifeChanged() (*oapi.WsResponse, error) {
	b := oapi.WsResponseBodyLifeChanged{}

	res := oapi.WsResponseFromType(oapi.WsResponseTypeLifeChanged)
	if err := res.Body.FromWsResponseBodyLifeChanged(b); err != nil {
		return nil, err
	}

	return res, nil
}

func (h *Hub) sendCardReset() (*oapi.WsResponse, error) {
	b := oapi.WsResponseBodyCardReset{}

	res := oapi.WsResponseFromType(oapi.WsResponseTypeCardReset)
	if err := res.Body.FromWsResponseBodyCardReset(b); err != nil {
		return nil, err
	}

	return res, nil
}

func (h *Hub) sendRailCreated() (*oapi.WsResponse, error) {
	b := oapi.WsResponseBodyRailCreated{}

	res := oapi.WsResponseFromType(oapi.WsResponseTypeRailCreated)
	if err := res.Body.FromWsResponseBodyRailCreated(b); err != nil {
		return nil, err
	}

	return res, nil
}

func (h *Hub) sendRailMerged() (*oapi.WsResponse, error) {
	b := oapi.WsResponseBodyRailMerged{}

	res := oapi.WsResponseFromType(oapi.WsResponseTypeRailMerged)
	if err := res.Body.FromWsResponseBodyRailMerged(b); err != nil {
		return nil, err
	}

	return res, nil
}

func (h *Hub) sendBlockCreated() (*oapi.WsResponse, error) {
	b := oapi.WsResponseBodyBlockCreated{}

	res := oapi.WsResponseFromType(oapi.WsResponseTypeBlockCreated)
	if err := res.Body.FromWsResponseBodyBlockCreated(b); err != nil {
		return nil, err
	}

	return res, nil
}
