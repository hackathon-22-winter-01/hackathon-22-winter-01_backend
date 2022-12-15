package ws

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/domain"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/oapi"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/usecases/repository"
	"github.com/labstack/echo/v4"
	"github.com/shiguredo/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 300000
)

type Client struct {
	hub    *Hub
	userID uuid.UUID
	conn   *websocket.Conn
	send   chan *oapi.WsResponse
	logger echo.Logger
}

func NewClient(hub *Hub, userID uuid.UUID, conn *websocket.Conn, logger echo.Logger) *Client {
	return &Client{
		hub:    hub,
		userID: userID,
		conn:   conn,
		send:   make(chan *oapi.WsResponse, 256),
		logger: logger,
	}
}

func (c *Client) readPump() error {
	defer func() {
		c.hub.Unregister(c)
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)

	if err := c.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		return err
	}

	c.conn.SetPongHandler(func(string) error {
		return c.conn.SetReadDeadline(time.Now().Add(pongWait))
	})

	for {
		req := new(oapi.WsRequest)
		if err := c.conn.ReadJSON(req); err != nil {
			if !websocket.IsCloseError(err) && !websocket.IsUnexpectedCloseError(err) {
				return err
			}

			break
		}

		if err := c.handleEvent(req); err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) writePump() error {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				return err
			}

			if !ok {
				if err := c.conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
					return err
				}

				return nil
			}

			if err := c.conn.WriteJSON(message); err != nil {
				return err
			}

		case <-ticker.C:
			if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				return err
			}

			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return err
			}

			res, err := oapi.NewWsResponseCardReset(time.Now())
			if err != nil {
				return err
			}

			if err := c.conn.WriteJSON(res); err != nil {
				return err
			}
		}
	}
}

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
	// テスト時は固定する
	cards := []oapi.Card{
		{Id: uuid.New(), Type: oapi.CardTypeCreateRail},
		{Id: uuid.New(), Type: oapi.CardTypeCreateBlock},
		{Id: uuid.New(), Type: oapi.CardTypeCreateRail},
		{Id: uuid.New(), Type: oapi.CardTypeCreateBlock},
		{Id: uuid.New(), Type: oapi.CardTypeCreateRail},
	}

	room, err := c.hub.roomRepo.FindRoom(repository.CommonRoomID) // TODO 適切なIDを指定する
	if err != nil {
		return err
	}

	players := make([]oapi.Player, len(room.Players))
	for i, p := range room.Players {
		players[i] = oapi.PlayerFromDomain(p)
	}

	res, err := oapi.NewWsResponseGameStarted(nowInJST(), cards, players)
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
		room, err := c.hub.roomRepo.FindRoom(repository.CommonRoomID)
		if err != nil {
			return err
		}

		target, ok := room.FindPlayer(b.TargetId)
		if !ok {
			return errors.New("player not found")
		}

		newRail := domain.NewRail()
		beforeRails := []*domain.Rail{newRail}
		afterRails := []*domain.Rail{}

		if l := len(target.Events); l > 0 {
			lastEvent := target.Events[l-1]
			beforeRails = lastEvent.AfterRails
			afterRails = append(beforeRails, newRail)
		}

		target.Events = append(target.Events, domain.NewCardEvent(
			uuid.New(),
			domain.RailCreated,
			c.userID,
			target.ID,
			beforeRails,
			afterRails,
		))

		res, err = oapi.NewWsResponseRailCreated(nowInJST(), uuid.New(), target.Main.ID, c.userID, b.TargetId)
		if err != nil {
			return err
		}

	case oapi.CardTypeCreateBlock:
		res, err = oapi.NewWsResponseBlockCreated(nowInJST(), c.userID, b.TargetId)
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
		res, err := oapi.NewWsResponseLifeChanged(nowInJST())
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
	res, err := oapi.NewWsResponseRailMerged(nowInJST())
	if err != nil {
		return err
	}

	c.bloadcast(res)

	return nil
}

func nowInJST() time.Time {
	return time.Now().In(time.FixedZone("Asia/Tokyo", 9*60*60))
}
