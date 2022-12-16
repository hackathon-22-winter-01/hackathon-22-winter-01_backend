package wshandler

import (
	"github.com/google/uuid"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/oapi"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/pkg/jst"
)

func (h *wsHandler) handleGameStartEvent(body oapi.WsRequest_Body) error {
	_, err := body.AsWsRequestBodyGameStartEvent()
	if err != nil {
		return err
	}

	// TODO: 初期カードを決めるロジックを書く
	// テスト時は固定する
	cards := []oapi.Card{
		{Id: uuid.New(), Type: oapi.CardTypeCreateRail},
		{Id: uuid.New(), Type: oapi.CardTypeCreateBlock},
		{Id: uuid.New(), Type: oapi.CardTypeCreateRail},
		{Id: uuid.New(), Type: oapi.CardTypeCreateBlock},
		{Id: uuid.New(), Type: oapi.CardTypeCreateRail},
	}

	players := make([]oapi.Player, len(h.room.Players))
	for i, p := range h.room.Players {
		players[i] = oapi.PlayerFromDomain(p)
	}

	res, err := oapi.NewWsResponseGameStarted(jst.Now(), cards, players)
	if err != nil {
		return err
	}

	if err := h.sender.Broadcast(h.room.ID, res); err != nil {
		return err
	}

	return nil
}
