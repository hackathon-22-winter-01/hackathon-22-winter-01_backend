package ws

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/domain"
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
	roomID, err := s.createOrJoinRoom(opts)
	if err != nil {
		return fmt.Errorf("failed to create or join room: %w", err)
	}

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

	if err := client.Broadcast(roomID, res); err != nil {
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

// createOrJoinRoom roomIDが指定されていれば参加、されていなければ作成する
func (s *streamer) createOrJoinRoom(opts ServeWsOpts) (uuid.UUID, error) {
	player := domain.NewPlayer(opts.PlayerID, opts.PlayerName)

	roomID, ok := opts.RoomID.Value()
	if ok {
		if err := s.hub.roomRepo.JoinRoom(roomID, player); err != nil {
			return uuid.Nil, fmt.Errorf("failed to join room: %w", err)
		}
	} else {
		room, err := s.hub.roomRepo.CreateRoom(player)
		if err != nil {
			return uuid.Nil, fmt.Errorf("failed to create room: %w", err)
		}

		roomID = room.ID
	}

	return roomID, nil
}
