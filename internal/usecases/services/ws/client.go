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

func (client *Client) readPump() error {
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

		if err := client.callEventHandler(req); err != nil {
			client.logger.Error("error: %v", err)

			break
		}
	}

	return nil
}

func (client *Client) writePump() error {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		client.conn.Close()
	}()
	for {
		select {
		case message, ok := <-client.send:
			if err := client.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {

				return err
			}
			if !ok {
				if err := client.conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
					return err
				}

				return nil
			}

			w, err := client.conn.NextWriter(websocket.TextMessage)
			if err != nil {

				return err
			}

			buf, err := json.Marshal(message)
			if err != nil {

				return err
			}

			if _, err := w.Write(buf); err != nil {

				return err
			}

			for i := 0; i < len(client.send); i++ {
				buf, err = json.Marshal(<-client.send)
				if err != nil {

					return err
				}

				if _, err := w.Write(buf); err != nil {

					return err
				}
			}

			if err := w.Close(); err != nil {

				return err
			}
		case <-ticker.C:
			if err := client.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {

				return err
			}
			if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {

				return err
			}
		}
	}
}
