package ws

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/oapi"
	"github.com/labstack/echo/v4"
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

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (client *Client) readPump() {
	defer func() {
		client.hub.Unregister(client)
		client.conn.Close()
	}()
	client.conn.SetReadLimit(maxMessageSize)
	if err := client.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		client.logger.Error("error:", err.Error())
	}

	client.conn.SetPongHandler(func(string) error { return client.conn.SetReadDeadline(time.Now().Add(pongWait)) })
	for {
		req := new(oapi.WsRequest)
		if err := client.conn.ReadJSON(req); err != nil {
			if !websocket.IsCloseError(err) && !websocket.IsUnexpectedCloseError(err) {
				client.logger.Error("websocket error occurred:", err.Error())
			}

			break
		}

		//TODO Handler
		if err := client.callEventHandler(req); err != nil {
			client.logger.Error("error: %v", err)

			break
		}
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (client *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		client.conn.Close()
	}()
	for {
		select {
		case message, ok := <-client.send:
			if err := client.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				client.logger.Error("error:", err.Error())

				return
			}
			if !ok {
				// The hub closed the channel.
				if err := client.conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
					client.logger.Error("error:", err.Error())
				}

				return
			}

			w, err := client.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				client.logger.Error("error: %v", err)

				return
			}

			buf, err := json.Marshal(message)
			if err != nil {
				client.logger.Error("error: %v", err)

				return
			}

			if _, err := w.Write(buf); err != nil {
				client.logger.Error("error:", err.Error())

				return
			}

			for i := 0; i < len(client.send); i++ {
				buf, err = json.Marshal(<-client.send)
				if err != nil {
					client.logger.Error("error: %v", err)

					return
				}

				if _, err := w.Write(buf); err != nil {
					client.logger.Error("error:", err.Error())

					return
				}
			}

			if err := w.Close(); err != nil {
				client.logger.Error("error: %v", err)

				return
			}
		case <-ticker.C:
			if err := client.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				client.logger.Error("error:", err.Error())

				return
			}
			if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {

				return
			}
		}
	}
}
