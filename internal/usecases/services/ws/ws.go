package ws

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/oapi"
	"github.com/labstack/echo/v4"
)

type Streamer interface {
	Run()
	ServeWS(w http.ResponseWriter, r *http.Request, uid uuid.UUID) error
}

type streamer struct {
	hub      *Hub
	upgrader websocket.Upgrader
	logger   echo.Logger
}

func NewStreamer(hub *Hub, logger echo.Logger) Streamer {
	stream := &streamer{
		hub: hub,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		logger: logger,
	}
	stream.Run()
	return stream
}

func (s *streamer) Run() {
	go s.hub.Run()
}

func (s *streamer) ServeWS(w http.ResponseWriter, r *http.Request, userId uuid.UUID) error {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return fmt.Errorf("failed to upgrade the HTTP server connection to the WebSocket protocol: %w", err)
	}

	client, err := s.addNewClient(userId, conn)
	if err != nil {
		return fmt.Errorf("failed to add new client: %w", err)
	}

	go client.writePump()
	go client.readPump()

	client.send <- &oapi.WsResponse{
		Type: "Hello",
		Body: oapi.WsResponse_Body{},
	}

	return nil
}

func (s *streamer) addNewClient(userId uuid.UUID, conn *websocket.Conn) (*Client, error) {
	client := NewClient(s.hub, userId, conn, s.logger)
	s.hub.Register(client)

	s.hub.mux.Lock()
	defer s.hub.mux.Unlock()

	_, ok := s.hub.clients[userId]

	if !ok {
		s.hub.clients[userId] = client
	}
	return client, nil
}
