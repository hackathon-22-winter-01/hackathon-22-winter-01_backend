package ws

import (
	"github.com/google/uuid"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/usecases/repository"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/pkg/log"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/pkg/sync"
	"go.uber.org/zap"
)

type Hub struct {
	clients      sync.Map[uuid.UUID, *Client]
	registerCh   chan *Client
	unregisterCh chan *Client
	roomRepo     repository.RoomRepository
}

func NewHub(roomRepo repository.RoomRepository) *Hub {
	return &Hub{
		clients:      sync.Map[uuid.UUID, *Client]{},
		registerCh:   make(chan *Client),
		unregisterCh: make(chan *Client),
		roomRepo:     roomRepo,
	}
}

func (h *Hub) Register(client *Client) {
	h.registerCh <- client
}

func (h *Hub) Unregister(client *Client) {
	h.unregisterCh <- client
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.registerCh:
			log.L().Info("new client has registered", zap.Stringer("playerID", client.playerID))
			h.clients.LoadOrStore(client.playerID, client)
		case client := <-h.unregisterCh:
			log.L().Info("a client has unregistered", zap.Stringer("playerID", client.playerID))
			close(client.send)
			h.clients.Delete(client.playerID)
		}
	}
}
