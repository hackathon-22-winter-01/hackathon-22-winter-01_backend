package ws

import (
	"sync"

	"github.com/google/uuid"
)

type Hub struct {
	clients      map[uuid.UUID]*Client
	registerCh   chan *Client
	unregisterCh chan *Client
	mux          sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients:      make(map[uuid.UUID]*Client),
		registerCh:   make(chan *Client),
		unregisterCh: make(chan *Client),
		mux:          sync.RWMutex{},
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
			h.register(client)
		case client := <-h.unregisterCh:
			if _, ok := h.clients[client.userID]; !ok {
				h.unregister(client)
			}
		}
	}
}

func (h *Hub) register(client *Client) {
	h.mux.Lock()
	defer h.mux.Unlock()

	client.logger.Infof("register client: %s", client.userID.String())

	h.clients[client.userID] = client
}

func (h *Hub) unregister(client *Client) {
	h.mux.Lock()
	defer h.mux.Unlock()

	client.logger.Infof("unregister client: %s", client.userID.String())

	close(client.send)
	delete(h.clients, client.userID)
}
