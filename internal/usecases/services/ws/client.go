package ws

import (
	"time"

	"github.com/google/uuid"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/oapi"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/usecases/services/ws/wshandler"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/pkg/log"
	"github.com/shiguredo/websocket"
	"go.uber.org/zap"
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
}

func NewClient(hub *Hub, userID uuid.UUID, conn *websocket.Conn) *Client {
	return &Client{
		hub:    hub,
		userID: userID,
		conn:   conn,
		send:   make(chan *oapi.WsResponse, 256),
	}
}

func (c *Client) readPump() error {
	defer c.hub.Unregister(c)
	c.conn.SetReadLimit(maxMessageSize)

	if err := c.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		return err
	}

	c.conn.SetPongHandler(func(string) error {
		return c.conn.SetReadDeadline(time.Now().Add(pongWait))
	})

	wh := wshandler.NewWsHandler(c.userID, c.hub.roomRepo, c)

	for {
		req := new(oapi.WsRequest)
		if err := c.conn.ReadJSON(req); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				return err
			}

			break
		}

		if len(req.Type) == 0 {
			continue
		}

		if err := wh.HandleEvent(req); err != nil {
			log.L().Error("failed to handle event", zap.Error(err), zap.String("eventType", string(req.Type)))
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
				log.L().Error("failed to set write deadline", zap.Error(err))
			}

			if !ok {
				if err := c.conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
					log.L().Error("failed to write close message", zap.Error(err))
				}

				return nil
			}

			if err := c.conn.WriteJSON(message); err != nil {
				log.L().Error("failed to write json", zap.Error(err))
			}

		case <-ticker.C:
			if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				log.L().Error("failed to set write deadline", zap.Error(err))
			}

			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.L().Error("failed to write ping message", zap.Error(err))
			}

			res, err := oapi.NewWsResponseCardReset(time.Now())
			if err != nil {
				log.L().Error("failed to create card reset response", zap.Error(err))
			}

			if err := c.conn.WriteJSON(res); err != nil {
				log.L().Error("failed to write json", zap.Error(err))
			}
		}
	}
}

func (c *Client) Bloadcast(roomID uuid.UUID, res *oapi.WsResponse) {
	// TODO: 全クライアントに送信してしまうためルーム内のクライアントだけに絞る
	c.hub.clients.Range(func(_ uuid.UUID, client *Client) bool {
		client.send <- res
		return true
	})
}
