package ws

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/oapi"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/pkg/jst"
	"github.com/labstack/echo/v4"
	"github.com/shiguredo/websocket"
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

func (s *streamer) ServeWS(w http.ResponseWriter, r *http.Request, userID uuid.UUID) error {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return fmt.Errorf("failed to upgrade the HTTP server connection to the WebSocket protocol: %w", err)
	}

	client, err := s.addNewClient(userID, conn)
	if err != nil {
		return fmt.Errorf("failed to add new client: %w", err)
	}

	go func() {
		if err := client.writePump(); err != nil {
			s.logger.Error(err)
		}
	}()
	go func() {
		if err := client.readPump(); err != nil {
			s.logger.Error(err)
		}
	}()

	res, err := oapi.NewWsResponseConnected(jst.Now(), userID)
	if err != nil {
		return fmt.Errorf("failed to send connected: %w", err)
	}

	client.send <- res

	return nil
}

func (s *streamer) addNewClient(userID uuid.UUID, conn *websocket.Conn) (*Client, error) {
	client := NewClient(s.hub, userID, conn, s.logger)
	s.hub.Register(client)
	s.hub.clients.LoadOrStore(userID, client)

	return client, nil
}
