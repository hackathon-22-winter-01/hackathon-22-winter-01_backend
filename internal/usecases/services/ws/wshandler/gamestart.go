package wshandler

import (
	"github.com/google/uuid"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/oapi"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/usecases/repository"
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

	room, err := h.roomRepo.FindRoom(repository.CommonRoomID) // TODO 適切なIDを指定する
	if err != nil {
		return err
	}

	players := make([]oapi.Player, len(room.Players))
	for i, p := range room.Players {
		players[i] = oapi.PlayerFromDomain(p)
	}

	res, err := oapi.NewWsResponseGameStarted(jst.Now(), cards, players)
	if err != nil {
		return err
	}

	h.sender.Bloadcast(room.ID, res)

	return nil
}
