package ws

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/oapi"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/pkg/jst"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/pkg/log"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/pkg/optional"
	"github.com/shiguredo/websocket"
	"go.uber.org/zap"
)

type Streamer interface {
	Run()
	ServeWS(w http.ResponseWriter, r *http.Request, opts ServeWsOpts) error
}

type ServeWsOpts struct {
	PlayerID   uuid.UUID
	PlayerName string
	RoomID     optional.Of[uuid.UUID]
}

type streamer struct {
	hub      *Hub
	upgrader websocket.Upgrader
}

func NewStreamer(hub *Hub) Streamer {
	stream := &streamer{
		hub: hub,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {

				return true
			},
		},
	}
	stream.Run()

	return stream
}

func (s *streamer) Run() {
	go s.hub.Run()
}

func (s *streamer) ServeWS(w http.ResponseWriter, r *http.Request, opts ServeWsOpts) error {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return fmt.Errorf("failed to upgrade the HTTP server connection to the WebSocket protocol: %w", err)
	}

	client, err := s.addNewClient(opts.PlayerID, conn)
	if err != nil {
		return fmt.Errorf("failed to add new client: %w", err)
	}

	go func() {
		if err := client.writePump(); err != nil {
			log.L().Error("failed to write pump", zap.Error(err))
		}
	}()
	go func() {
		if err := client.readPump(); err != nil {
			log.L().Error("failed to read pump", zap.Error(err))
		}
	}()

	res, err := oapi.NewWsResponseConnected(jst.Now(), opts.PlayerID)
	if err != nil {
		return fmt.Errorf("failed to send connected: %w", err)
	}

	// TODO: roomIDが指定されていれば参加、されていなければ作成する
	room, err := s.hub.roomRepo.FindRoomFromPlayerID(opts.PlayerID)
	if err != nil {
		return fmt.Errorf("failed to find room from player id: %w", err)
	}

	if err := client.Broadcast(room.ID, res); err != nil {
		return fmt.Errorf("failed to broadcast: %w", err)
	}

	return nil
}

func (s *streamer) addNewClient(playerID uuid.UUID, conn *websocket.Conn) (*Client, error) {
	client := NewClient(s.hub, playerID, conn)
	s.hub.Register(client)
	s.hub.clients.LoadOrStore(playerID, client)

	return client, nil
}
