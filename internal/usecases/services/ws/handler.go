package ws

import (
	"errors"

	"github.com/google/uuid"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/oapi"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/usecases/repository"
)

func (h *Hub) handleEvent(req *oapi.WsRequest) error {
	switch req.Type {
	case oapi.WsRequestTypeGameStartEvent:
		return h.handleGameStartEvent(req.Body)

	case oapi.WsRequestTypeCardEvent:
		return h.handleCardEvent(req.Body)

	case oapi.WsRequestTypeLifeEvent:
		return h.handleLifeEvent(req.Body)

	case oapi.WsRequestTypeRailMergeEvent:
		return h.handleRailMergeEvent(req.Body)

	default:
		return errors.New("invalid request type")
	}
}

func (h *Hub) handleGameStartEvent(body oapi.WsRequest_Body) error {
	_, err := body.AsWsRequestBodyGameStartEvent()
	if err != nil {
		return err
	}

	// TODO: 初期カードを決めるロジックを書く
	cards := []oapi.Card{
		{Id: uuid.New(), Type: oapi.CardTypeCreateRail},
		{Id: uuid.New(), Type: oapi.CardTypeCreateRail},
		{Id: uuid.New(), Type: oapi.CardTypeCreateBlock},
		{Id: uuid.New(), Type: oapi.CardTypeCreateBlock},
		{Id: uuid.New(), Type: oapi.CardTypeCreateBlock},
	}

	room, err := h.roomRepo.FindRoom(repository.CommonRoomID) // TODO 適切なIDを指定する
	if err != nil {
		return err
	}

	players := make([]oapi.Player, len(room.Players))
	for i, p := range room.Players {
		players[i] = oapi.PlayerFromDomain(p)
	}

	res, err := oapi.NewWsResponseGameStarted(cards, players)
	if err != nil {
		return err
	}

	h.bloadcast(res)

	return nil
}

func (h *Hub) handleCardEvent(body oapi.WsRequest_Body) error {
	b, err := body.AsWsRequestBodyCardEvent()
	if err != nil {
		return err
	}

	var res *oapi.WsResponse

	switch b.Type {
	case oapi.CardTypeCreateRail:
		res, err = oapi.NewWsResponseRailCreated()
		if err != nil {
			return err
		}

	case oapi.CardTypeCreateBlock:
		res, err = oapi.NewWsResponseBlockCreated()
		if err != nil {
			return err
		}

	default:
		return errors.New("invalid card type")
	}

	h.bloadcast(res)

	return nil
}

func (h *Hub) handleLifeEvent(body oapi.WsRequest_Body) error {
	b, err := body.AsWsRequestBodyLifeEvent()
	if err != nil {
		return err
	}

	switch b.Type {
	case oapi.LifeEventTypeDecrement:
		res, err := oapi.NewWsResponseLifeChanged()
		if err != nil {
			return err
		}

		h.bloadcast(res)

	default:
		return errors.New("invalid life type")
	}

	return nil
}

func (h *Hub) handleRailMergeEvent(_ oapi.WsRequest_Body) error {
	res, err := oapi.NewWsResponseRailMerged()
	if err != nil {
		return err
	}

	h.bloadcast(res)

	return nil
}
