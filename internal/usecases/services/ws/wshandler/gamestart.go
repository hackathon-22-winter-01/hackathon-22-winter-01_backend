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

	domainCards, err := h.cardRepo.DrawCards(h.room.ID, h.playerID, 5)
	if err != nil {
		return err
	}

	cards := make([]oapi.Card, len(domainCards))
	for i, c := range domainCards {
		cards[i] = oapi.CardFromDomain(c)
	}

	players := make([]oapi.Player, len(h.room.Players))
	for i, p := range h.room.Players {
		players[i] = oapi.PlayerFromDomain(p)
	}

	res, err := oapi.NewWsResponseGameStarted(jst.Now(), cards, players)
	if err != nil {
		return err
	}

	// TODO: 全員に同じカードが渡ってしまう
	if err := h.sender.Broadcast(h.room.ID, res); err != nil {
		return err
	}

	return nil
}
