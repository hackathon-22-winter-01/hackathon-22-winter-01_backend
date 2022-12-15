package ws

import (
	"github.com/google/uuid"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/oapi"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/usecases/repository"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/pkg/sync"
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
			client.logger.Infof("register client: %s", client.userID.String())
			h.clients.LoadOrStore(client.userID, client)
		case client := <-h.unregisterCh:
			client.logger.Infof("unregister client: %s", client.userID.String())
			close(client.send)
			h.clients.Delete(client.userID)
		}
	}
}

func (h *Hub) bloadcast(res *oapi.WsResponse) {
	h.clients.Range(func(_ uuid.UUID, client *Client) bool {
		client.send <- res
		return true
	})
}
