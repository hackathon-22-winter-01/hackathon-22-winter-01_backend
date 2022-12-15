package ws

import (
	"errors"

	"github.com/google/uuid"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/oapi"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/usecases/repository"
)

func (c *Client) bloadcast(res *oapi.WsResponse) {
	// TODO: 全クライアントに送信してしまうためルーム内のクライアントだけに絞る
	c.hub.clients.Range(func(_ uuid.UUID, client *Client) bool {
		client.send <- res
		return true
	})
}

func (c *Client) handleEvent(req *oapi.WsRequest) error {
	switch req.Type {
	case oapi.WsRequestTypeGameStartEvent:
		return c.handleGameStartEvent(req.Body)

	case oapi.WsRequestTypeCardEvent:
		return c.handleCardEvent(req.Body)

	case oapi.WsRequestTypeLifeEvent:
		return c.handleLifeEvent(req.Body)

	case oapi.WsRequestTypeRailMergeEvent:
		return c.handleRailMergeEvent(req.Body)

	default:
		return errors.New("invalid request type")
	}
}

func (c *Client) handleGameStartEvent(body oapi.WsRequest_Body) error {
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

	room, err := c.hub.roomRepo.FindRoom(repository.CommonRoomID) // TODO 適切なIDを指定する
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

	c.bloadcast(res)

	return nil
}

func (c *Client) handleCardEvent(body oapi.WsRequest_Body) error {
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

	c.bloadcast(res)

	return nil
}

func (c *Client) handleLifeEvent(body oapi.WsRequest_Body) error {
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

		c.bloadcast(res)

	default:
		return errors.New("invalid life type")
	}

	return nil
}

func (c *Client) handleRailMergeEvent(_ oapi.WsRequest_Body) error {
	res, err := oapi.NewWsResponseRailMerged()
	if err != nil {
		return err
	}

	c.bloadcast(res)

	return nil
}
