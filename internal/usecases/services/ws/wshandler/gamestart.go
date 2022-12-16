package wshandler

import (
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/oapi"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/pkg/jst"
)

func (h *wsHandler) handleGameStartEvent(body oapi.WsRequest_Body) error {
	_, err := body.AsWsRequestBodyGameStartEvent()
	if err != nil {
		return err
	}

	players := make([]oapi.Player, len(h.room.Players))
	for i, p := range h.room.Players {
		players[i] = oapi.PlayerFromDomain(p)
	}

	resFunc := func() (*oapi.WsResponse, error) {
		res, err := oapi.NewWsResponseGameStarted(jst.Now(), players)
		if err != nil {
			return nil, err
		}

		return res, nil
	}

	h.room.StartedAt = jst.Now()

	if err := h.sender.BroadcastDynamic(h.room.ID, resFunc); err != nil {
		return err
	}

	return nil
}
